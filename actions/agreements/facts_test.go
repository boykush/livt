package agreements_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
)

// facts runs the action's own filter over what a forge would return, so what
// is tested is the file the action ships rather than a copy of its logic.
func facts(t *testing.T, pr, reviews string) (change, author string, approvedBy []string) {
	t.Helper()
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq is not installed; the action's filter cannot be run here")
	}
	dir := t.TempDir()
	prPath, reviewsPath := filepath.Join(dir, "pr.json"), filepath.Join(dir, "reviews.json")
	for path, body := range map[string]string{prPath: pr, reviewsPath: reviews} {
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out, err := exec.Command("jq", "-n", "--slurpfile", "pr", prPath, "--slurpfile", "reviews", reviewsPath, "-f", "facts.jq").Output()
	if err != nil {
		t.Fatalf("jq: %v", err)
	}
	var got struct {
		Change     string   `json:"change"`
		Author     string   `json:"author"`
		ApprovedBy []string `json:"approved_by"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("decode %s: %v", out, err)
	}
	return got.Change, got.Author, got.ApprovedBy
}

const byAlice = `{"html_url":"https://forge.example/acme/specs/pull/7","user":{"login":"alice","type":"User"}}`

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-05/example/EX-01
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-02/example/EX-04
// The action reads who wrote the pull request and who approved it, and hands
// livt nothing else. An approval counts whoever gave it: a reviewer the forge
// would not count towards its own required reviews is still someone agreeing.
func TestFactsCarryTheAuthorAndEveryApprover(t *testing.T) {
	reviews := `[
	  {"user":{"login":"bob"},"state":"APPROVED","submitted_at":"2026-10-01T00:00:00Z","author_association":"MEMBER"},
	  {"user":{"login":"carol"},"state":"APPROVED","submitted_at":"2026-10-01T00:01:00Z","author_association":"NONE"}
	]`
	change, author, approvedBy := facts(t, byAlice, reviews)

	if change != "https://forge.example/acme/specs/pull/7" || author != "@alice" {
		t.Errorf("change = %q, author = %q, want the pull request's own URL and @alice", change, author)
	}
	if want := []string{"@bob", "@carol"}; !slices.Equal(approvedBy, want) {
		t.Errorf("approved_by = %v, want %v", approvedBy, want)
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-02/example/EX-03
// A pull request a bot opened has no person as its author, so whoever it was
// written for has to approve it like anyone else.
func TestFactsNameNoAuthorForABot(t *testing.T) {
	pr := `{"html_url":"https://forge.example/acme/specs/pull/8","user":{"login":"claude[bot]","type":"Bot"}}`
	_, author, _ := facts(t, pr, `[]`)

	if author != "" {
		t.Errorf("author = %q, want none for a bot", author)
	}
}

// A person stands on their latest review that took a side: an approval taken
// back or dismissed is no approval, and a comment afterwards does not undo one.
func TestFactsReadEachPersonsLatestStanding(t *testing.T) {
	reviews := `[
	  {"user":{"login":"bob"},"state":"APPROVED","submitted_at":"2026-10-01T00:00:00Z"},
	  {"user":{"login":"bob"},"state":"CHANGES_REQUESTED","submitted_at":"2026-10-02T00:00:00Z"},
	  {"user":{"login":"carol"},"state":"DISMISSED","submitted_at":"2026-10-01T00:00:00Z"},
	  {"user":{"login":"dave"},"state":"APPROVED","submitted_at":"2026-10-01T00:00:00Z"},
	  {"user":{"login":"dave"},"state":"COMMENTED","submitted_at":"2026-10-03T00:00:00Z"},
	  {"user":null,"state":"APPROVED","submitted_at":"2026-10-01T00:00:00Z"}
	]`
	_, _, approvedBy := facts(t, byAlice, reviews)

	if want := []string{"@dave"}; !slices.Equal(approvedBy, want) {
		t.Errorf("approved_by = %v, want %v", approvedBy, want)
	}
}
