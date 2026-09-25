package gitrev

import (
	"os/exec"
	"testing"
)

func TestShortOfANonGitDirIsEmpty(t *testing.T) {
	if got := Short(t.TempDir()); got != "" {
		t.Fatalf("Short of a non-git dir = %q, want empty", got)
	}
}

func TestShortInAGitRepoReturnsTheRevision(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init"},
		{"-c", "user.email=t@example.com", "-c", "user.name=t", "commit", "--allow-empty", "-m", "init"},
	} {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Skipf("git setup failed (%v): %s", err, out)
		}
	}
	if got := Short(dir); got == "" {
		t.Fatal("Short in a git repo = empty, want a revision")
	}
}
