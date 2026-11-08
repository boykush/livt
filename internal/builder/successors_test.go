package builder

import (
	"path/filepath"
	"strings"
	"testing"
)

// livt:automates livt://mapping/trace-test-to-rule/rule/R-11
func TestBuildChecksSuccessors(t *testing.T) {
	// livt:automates livt://mapping/trace-test-to-rule/rule/R-11/example/EX-01
	t.Run("a retirement pointing at a rule that does not exist fails the build", func(t *testing.T) {
		b := emptyDirsBuilder(t)
		writeFile(t, filepath.Join(b.MappingsDir, "demo.yaml"),
			"rules:\n"+
				"  - id: R-01\n"+
				"    name: old\n"+
				"    status: retired\n"+
				"    superseded_by:\n"+
				"      - livt://mapping/demo/rule/R-99\n")

		err := b.Build()
		if err == nil {
			t.Fatal("build passed with a successor that does not resolve")
		}
		for _, want := range []string{"demo.yaml", "rule R-01", "livt://mapping/demo/rule/R-99"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error does not name %q:\n%v", want, err)
			}
		}
	})

	t.Run("every broken successor is reported in one run", func(t *testing.T) {
		b := emptyDirsBuilder(t)
		writeFile(t, filepath.Join(b.MappingsDir, "demo.yaml"),
			"rules:\n"+
				"  - id: R-01\n"+
				"    name: kept\n"+
				"    examples:\n"+
				"      - id: EX-01\n"+
				"        name: old example\n"+
				"        retired: true\n"+
				"        superseded_by:\n"+
				"          - livt://mapping/demo/rule/R-01/example/EX-09\n"+
				"questions:\n"+
				"  - id: Q-01\n"+
				"    text: old question\n"+
				"    retired: true\n"+
				"    superseded_by:\n"+
				"      - livt://mapping/elsewhere/rule/R-01\n"+
				"      - not-a-uri\n")

		err := b.Build()
		if err == nil {
			t.Fatal("build passed with successors that do not resolve")
		}
		for _, want := range []string{
			"example R-01/EX-01: superseded_by livt://mapping/demo/rule/R-01/example/EX-09",
			"question Q-01: superseded_by livt://mapping/elsewhere/rule/R-01",
			"question Q-01: superseded_by not-a-uri",
		} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error does not report %q:\n%v", want, err)
			}
		}
	})

	t.Run("a successor in another mapping, itself retired, resolves", func(t *testing.T) {
		b := emptyDirsBuilder(t)
		writeFile(t, filepath.Join(b.MappingsDir, "demo.yaml"),
			"rules:\n"+
				"  - id: R-01\n"+
				"    name: old\n"+
				"    status: rejected\n"+
				"    superseded_by:\n"+
				"      - livt://mapping/other/rule/R-02\n")
		writeFile(t, filepath.Join(b.MappingsDir, "other.yaml"),
			"rules:\n"+
				"  - id: R-02\n"+
				"    name: also closed\n"+
				"    status: retired\n")

		if err := b.Build(); err != nil {
			t.Fatal(err)
		}
	})
}
