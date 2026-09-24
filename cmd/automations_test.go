package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/automation"
)

// livt:automates livt://mapping/collect-automations/rule/R-07/example/EX-05
func TestReportChangedAnswersTrueWhenTheRangeCannotBeRead(t *testing.T) {
	var out, warn bytes.Buffer
	// Not a git repository, which is the shape a new branch or a force-push
	// also arrives in: the range is unreadable rather than empty.
	reportChanged(&out, &warn, t.TempDir(), "0000000000000000000000000000000000000000", "HEAD")

	if out.String() != "true\n" {
		t.Errorf("stdout was %q, want %q: an unreadable range must not read as no change", out.String(), "true\n")
	}
	if !strings.Contains(warn.String(), "warning:") {
		t.Errorf("the run answered true without saying why on stderr: %q", warn.String())
	}
}

// newReportedRepo is the resolve tests' livt repository -- one mapping holding
// a live rule and a retired one -- with a collected report added, so a
// citation has somewhere to resolve to and somewhere to miss.
func newReportedRepo(t *testing.T, citations ...automation.Citation) string {
	t.Helper()
	root := newTestRepo(t)
	data, err := json.MarshalIndent(automation.Report{Repo: "acme/app", Citations: citations}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "automations", "acme", "app.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// livt:automates livt://mapping/collect-automations/rule/R-05/example/EX-12
// livt:automates livt://mapping/collect-automations/rule/R-05/example/EX-13
// A citation naming a point of the spec that is not there fails the check, and
// is named with the file and line that wrote it -- the report says which test
// to fix, and a mistyped URI is only ever fixed in the test.
func TestVerifyCitationsFailsOnACitationThatResolvesToNothing(t *testing.T) {
	root := newReportedRepo(t,
		automation.Citation{URI: "livt://mapping/demo/rule/R-01", File: "order_test.go", Line: 7},
		automation.Citation{URI: "livt://mapping/demo/rule/R-99", File: "cart_test.go", Line: 12},
		// Nothing livt collects writes this, and the check is here because the
		// report is trusted rather than re-derived.
		automation.Citation{URI: "R-01", File: "cart_test.go", Line: 30},
	)

	var out bytes.Buffer
	err := verifyCitations(&out, root, filepath.Join(root, "automations"))
	if err == nil {
		t.Fatal("the check passed a report citing a rule the livt repository does not hold")
	}
	if !strings.Contains(err.Error(), "2 of 3") {
		t.Errorf("error was %q, want it to count the unresolvable citations against the whole", err)
	}
	for _, want := range []string{"acme/app", "cart_test.go:12", "livt://mapping/demo/rule/R-99", "cart_test.go:30"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output %q does not name %q, so nobody can find the citation to fix", out.String(), want)
		}
	}
	if strings.Contains(out.String(), "order_test.go") {
		t.Errorf("output %q names a citation that resolves, burying the two that do not", out.String())
	}
}

// livt:automates livt://mapping/derive-automation-status/rule/R-07/example/EX-02
// A citation of a retired rule resolves, so it passes. The rule closed on the
// spec's side while the test naming it goes on running and passing, and
// failing its build over that would reach further than the decision did.
func TestVerifyCitationsPassesACitationOfARetiredRule(t *testing.T) {
	root := newReportedRepo(t,
		automation.Citation{URI: "livt://mapping/demo/rule/R-02", File: "legacy_test.go", Line: 3},
	)

	var out bytes.Buffer
	if err := verifyCitations(&out, root, filepath.Join(root, "automations")); err != nil {
		t.Fatalf("the check failed a citation of a retired rule: %v (%s)", err, out.String())
	}
}
