package mcp

import (
	"os/exec"
	"testing"
)

func TestSpecVersionNonGitIsEmpty(t *testing.T) {
	if got := specVersion(t.TempDir()); got != "" {
		t.Fatalf("specVersion of non-git dir = %q, want empty", got)
	}
}

func TestSpecVersionInGitRepoReturnsRevision(t *testing.T) {
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
	if got := specVersion(dir); got == "" {
		t.Fatal("specVersion in git repo = empty, want a revision")
	}
}
