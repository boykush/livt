package automation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repo writes files into a scratch directory and returns it, so each test
// scans a tree it fully describes.
func repo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func scan(t *testing.T, root string, opts Options) (Report, []string) {
	t.Helper()
	report, warnings, err := Scan(root, opts)
	if err != nil {
		t.Fatal(err)
	}
	return report, warnings
}

func uris(r Report) []string {
	out := make([]string, 0, len(r.Citations))
	for _, c := range r.Citations {
		out = append(out, c.URI)
	}
	return out
}

// livt:automates livt://mapping/collect-automations/rule/R-01/example/EX-01
func TestScanCollectsAMarkedLine(t *testing.T) {
	root := repo(t, map[string]string{"a_test.go": "" +
		"// " + Marker + " livt://mapping/checkout/rule/R-01/example/EX-02\n" +
		"func TestOutOfStock(t *testing.T) {}\n"})

	report, warnings := scan(t, root, Options{})

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(report.Citations) != 1 {
		t.Fatalf("got %d citations, want 1: %v", len(report.Citations), uris(report))
	}
	got := report.Citations[0]
	if got.URI != "livt://mapping/checkout/rule/R-01/example/EX-02" || got.File != "a_test.go" || got.Line != 1 {
		t.Errorf("got %+v", got)
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-01/example/EX-02
func TestScanReadsEveryLanguagesCommentSyntax(t *testing.T) {
	root := repo(t, map[string]string{
		"a_test.go":  "// " + Marker + " livt://mapping/checkout/rule/R-01\n",
		"b_test.py":  "# " + Marker + " livt://mapping/checkout/rule/R-02\n",
		"c_test.sql": "-- " + Marker + " livt://mapping/checkout/rule/R-03\n",
		"d_test.rb":  "  #  " + Marker + " livt://mapping/checkout/rule/R-04\n",
	})

	report, warnings := scan(t, root, Options{})

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(report.Citations) != 4 {
		t.Fatalf("got %d citations, want 4: %v", len(report.Citations), uris(report))
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-01/example/EX-03
func TestScanLeavesAnUnmarkedURIAlone(t *testing.T) {
	root := repo(t, map[string]string{"handler.go": "" +
		"// livt://mapping/checkout/rule/R-01 explains why this branch exists\n"})

	report, _ := scan(t, root, Options{})

	if len(report.Citations) != 0 {
		t.Fatalf("a reference was collected as a claim: %v", uris(report))
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-01/example/EX-04
func TestScanTakesOneURIPerLine(t *testing.T) {
	root := repo(t, map[string]string{"a_test.go": "" +
		"// " + Marker + " livt://mapping/checkout/rule/R-01/example/EX-01\n" +
		"// " + Marker + " livt://mapping/checkout/rule/R-01/example/EX-02\n" +
		"func TestBoth(t *testing.T) {}\n"})

	report, warnings := scan(t, root, Options{})

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(report.Citations) != 2 {
		t.Fatalf("got %d citations, want 2: %v", len(report.Citations), uris(report))
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-01
func TestScanWarnsRatherThanDroppingAMalformedClaim(t *testing.T) {
	root := repo(t, map[string]string{"a_test.go": "" +
		"// " + Marker + " livt://mapping/checkout/rule/R-01/example/EX-01 and EX-02\n"})

	report, warnings := scan(t, root, Options{})

	if len(report.Citations) != 0 {
		t.Errorf("a half-read claim was collected: %v", uris(report))
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "a_test.go:1") {
		t.Fatalf("got warnings %v, want one naming the line", warnings)
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-02/example/EX-01
func TestScanCollectsRulesAndExamplesIndependently(t *testing.T) {
	root := repo(t, map[string]string{"a_test.go": "" +
		"// " + Marker + " livt://mapping/checkout/rule/R-01\n" +
		"// " + Marker + " livt://mapping/checkout/rule/R-01/example/EX-01\n"})

	report, _ := scan(t, root, Options{})

	want := []string{"livt://mapping/checkout/rule/R-01", "livt://mapping/checkout/rule/R-01/example/EX-01"}
	got := uris(report)
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-02/example/EX-03
func TestScanIgnoresTheStructureAroundACitation(t *testing.T) {
	nested := "" +
		"func TestRule(t *testing.T) {\n" +
		"\t// " + Marker + " livt://mapping/checkout/rule/R-01/example/EX-01\n" +
		"\tt.Run(\"first\", func(t *testing.T) {})\n" +
		"}\n"

	report, _ := scan(t, repo(t, map[string]string{"a_test.go": nested}), Options{})

	if len(report.Citations) != 1 || report.Citations[0].Line != 2 {
		t.Fatalf("got %+v", report.Citations)
	}
}

func TestScanSkipsPlaceholdersAndHiddenTrees(t *testing.T) {
	root := repo(t, map[string]string{
		"README.md":              "// " + Marker + " livt://mapping/{story-key}/rule/{rule-id}\n",
		".claude/worktrees/a.go": "// " + Marker + " livt://mapping/checkout/rule/R-01\n",
		"node_modules/b.go":      "// " + Marker + " livt://mapping/checkout/rule/R-02\n",
	})

	report, warnings := scan(t, root, Options{})

	if len(report.Citations) != 0 {
		t.Fatalf("got %v, want nothing collected", uris(report))
	}
	if len(warnings) != 0 {
		t.Fatalf("a documented template warned: %v", warnings)
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-06/example/EX-03
func TestScanPinsLineURLsToTheScannedRevision(t *testing.T) {
	root := repo(t, map[string]string{"a_test.go": "" +
		"// " + Marker + " livt://mapping/checkout/rule/R-01\n"})

	report, _ := scan(t, root, Options{
		Rev:         "abc123",
		Base:        "https://github.com/acme/backend",
		URLTemplate: githubURLTemplate,
	})

	want := "https://github.com/acme/backend/blob/abc123/a_test.go#L1"
	if got := report.Citations[0].URL; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-06/example/EX-04
func TestScanLeavesTheURLEmptyWhenNoForgeIsKnown(t *testing.T) {
	root := repo(t, map[string]string{"a_test.go": "" +
		"// " + Marker + " livt://mapping/checkout/rule/R-01\n"})

	report, _ := scan(t, root, Options{Rev: "abc123", Base: "https://git.acme.com/backend"})

	if got := report.Citations[0]; got.URL != "" || got.File == "" {
		t.Errorf("got %+v, want the file and line without a URL", got)
	}
}
