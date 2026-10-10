package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const proposedMapping = "rules:\n  - id: R-01\n    name: first\n    status: proposed\n    decision_makers: [\"@alice\", \"@bob\"]\n"

// agreementsRepo commits one mapping as the base a change is read against.
func agreementsRepo(t *testing.T, mapping string) string {
	t.Helper()
	root := t.TempDir()
	writeTo(t, filepath.Join(root, "discoveries", "example-mappings", "checkout.yaml"), mapping)
	for _, args := range [][]string{
		{"init", "-q"},
		{"add", "."},
		{"-c", "user.name=t", "-c", "user.email=t@example.com", "-c", "commit.gpgsign=false", "commit", "-q", "-m", "base"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return root
}

func writeTo(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func approvals(t *testing.T, json string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "approvals.json")
	writeTo(t, path, json)
	return path
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-01/example/EX-07
// A change that accepts no rule naming anyone has nothing to be held to, and
// the command says so and passes rather than staying silent.
func TestCheckAgreementsPassesAChangeWithNothingToCheck(t *testing.T) {
	root := agreementsRepo(t, proposedMapping)

	var out bytes.Buffer
	if err := checkAgreements(&out, root, "HEAD", "", "", false); err != nil {
		t.Fatalf("check failed on an unchanged tree: %v", err)
	}
	if !strings.Contains(out.String(), "nothing to check") {
		t.Errorf("output = %q, want it to say there was nothing to check", out.String())
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-04/example/EX-04
// A person turning the proposal accepted by hand is held to the same check as
// a rewrite would be: refused while someone is missing, passed once they are not.
func TestCheckAgreementsHoldsAHandWrittenAcceptance(t *testing.T) {
	root := agreementsRepo(t, proposedMapping)
	writeTo(t, filepath.Join(root, "discoveries", "example-mappings", "checkout.yaml"),
		strings.Replace(proposedMapping, "status: proposed", "status: accepted", 1))

	var out bytes.Buffer
	err := checkAgreements(&out, root, "HEAD", "", approvals(t, `{"change":"pr/2","author":"@alice"}`), false)
	if err == nil || !strings.Contains(out.String(), "waiting for: @bob") {
		t.Fatalf("err = %v, output = %q, want it refused while @bob is missing", err, out.String())
	}

	out.Reset()
	if err := checkAgreements(&out, root, "HEAD", "", approvals(t, `{"change":"pr/2","author":"@alice","approved_by":["@bob"]}`), false); err != nil {
		t.Fatalf("check failed with everyone agreed: %v\n%s", err, out.String())
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-04/example/EX-01
// Rewriting happens only when asked for. Without --write the proposal on the
// working tree is left as it is however many people have approved. It needs a
// working tree to rewrite, so it is refused against a revision.
func TestCheckAgreementsRewritesOnlyWhenAsked(t *testing.T) {
	root := agreementsRepo(t, "rules: []\n")
	path := filepath.Join(root, "discoveries", "example-mappings", "checkout.yaml")
	writeTo(t, path, proposedMapping)
	everyone := approvals(t, `{"change":"pr/1","author":"@alice","approved_by":["@bob"]}`)

	var out bytes.Buffer
	if err := checkAgreements(&out, root, "HEAD", "", everyone, false); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(path); string(data) != proposedMapping {
		t.Fatalf("mapping changed without --write: %q", data)
	}

	out.Reset()
	if err := checkAgreements(&out, root, "HEAD", "", everyone, true); err != nil {
		t.Fatalf("write failed: %v\n%s", err, out.String())
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "status: accepted") {
		t.Errorf("mapping = %q, want the proposal accepted", data)
	}
	if _, err := os.Stat(filepath.Join(root, "agreements", "checkout.json")); err != nil {
		t.Errorf("no record written: %v", err)
	}

	if err := checkAgreements(&out, root, "HEAD", "HEAD", everyone, true); err == nil {
		t.Error("--write accepted a head revision, want it refused: there is no working tree to rewrite")
	}
}
