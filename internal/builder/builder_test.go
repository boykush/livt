package builder

import (
	"os"
	"path/filepath"
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

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
