package builder

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// TestBuildDropsPagesForRemovedResources covers the rebuild path. A page whose
// source file was renamed or deleted between builds must not survive: the
// output directory is not cleared by anything else, so `livt serve` would keep
// serving a resource that no longer exists.
func TestBuildDropsPagesForRemovedResources(t *testing.T) {
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.MappingsDir, "kept.yaml"), "rules: []\n")
	writeFile(t, filepath.Join(b.StoriesDir, "kept.md"), "---\nname: Kept story\n---\n\nbody\n")
	writeFile(t, filepath.Join(b.USMDir, "kept.yaml"), "name: kept\nactivities: []\n")

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	// Stand in for pages left behind by resources that are now gone.
	orphans := []string{
		filepath.Join(b.OutDir, "mapping", "removed.html"),
		filepath.Join(b.OutDir, "story", "removed.html"),
		filepath.Join(b.OutDir, "story-map", "removed.html"),
	}
	for _, p := range orphans {
		writeFile(t, p, "<html>stale</html>")
	}

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	for _, p := range orphans {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Fatalf("%s survived the rebuild (stat error: %v)", p, err)
		}
	}
	for _, p := range []string{
		filepath.Join(b.OutDir, "mapping", "kept.html"),
		filepath.Join(b.OutDir, "story", "kept.html"),
		filepath.Join(b.OutDir, "story-map", "kept.html"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("rebuild dropped a live page: %v", err)
		}
	}
}

// TestBuildKeepsUnrelatedFilesInOutDir guards the blast radius of that clearing:
// --out is user supplied, so a build must not remove files it did not write,
// such as the CNAME or .nojekyll a Pages deploy keeps beside the output.
func TestBuildKeepsUnrelatedFilesInOutDir(t *testing.T) {
	b := emptyDirsBuilder(t)
	keep := filepath.Join(b.OutDir, "CNAME")
	writeFile(t, keep, "example.com\n")

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(keep); err != nil {
		t.Fatalf("build removed a file it did not write: %v", err)
	}
}

// livt://mapping/trace-test-to-rule/rule/R-05/example/EX-02: the sidebar's Tasks
// badge counts what the Tasks page lists, so a retired question or rule is out
// of the number as well as out of the list.
func TestComputeCountsLeavesRetiredItemsOutOfTasks(t *testing.T) {
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.MappingsDir, "demo.yaml"),
		"rules:\n"+
			"  - id: R-01\n"+
			"    name: 現役の未自動化ルール\n"+
			"  - id: R-02\n"+
			"    name: 退役したルール\n"+
			"    retired: true\n"+
			"questions:\n"+
			"  - id: Q-01\n"+
			"    text: 現役の疑問\n"+
			"  - id: Q-02\n"+
			"    text: 退役した疑問\n"+
			"    retired: true\n")

	counts, err := b.computeCounts()
	if err != nil {
		t.Fatal(err)
	}
	if counts.tasks != 2 {
		t.Fatalf("tasks = %d, want 2 (the live question and the live un-automated rule)", counts.tasks)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// livt://mapping/reflect-every-artifact-edit-in-preview/rule/R-02/example/EX-01: a new input
// the build reads joins what is reflected, without anyone remembering to add
// it. This test is what keeps `livt serve` honest.
// The watch list is derived from InputDirs, so an input directory added to
// Builder and left out of it would be built from but never watched — the site
// would go stale under an edit with no sign anything was missed. Fields are
// read by reflection rather than listed here so a new one fails this test
// instead of quietly slipping past it.
func TestInputDirsCoversEveryInputDirectory(t *testing.T) {
	b := emptyDirsBuilder(t)

	watched := make(map[string]bool)
	for _, dir := range b.InputDirs() {
		watched[dir] = true
	}

	typ := reflect.TypeOf(b)
	val := reflect.ValueOf(b)
	for i := range typ.NumField() {
		name := typ.Field(i).Name
		// OutDir is written, not read: watching it would rebuild on the output
		// of the rebuild that just ran.
		if !strings.HasSuffix(name, "Dir") || name == "OutDir" {
			continue
		}
		if !watched[val.Field(i).String()] {
			t.Errorf("%s is an input directory that InputDirs does not return", name)
		}
	}
}

// livt://mapping/reflect-every-artifact-edit-in-preview/rule/R-02/example/EX-02: the output
// is left out. It lives inside the repository in the default layout, so
// watching it would make every rebuild trigger the next one.
func TestInputDirsExcludesOutDir(t *testing.T) {
	b := emptyDirsBuilder(t)
	if slices.Contains(b.InputDirs(), b.OutDir) {
		t.Error("InputDirs returned OutDir; a rebuild would retrigger itself")
	}
}
