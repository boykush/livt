package mcp

import (
	"os/exec"
	"strings"
)

// specVersion returns the git revision (short HEAD) of the livt repository at root, so
// every tool result records which version of the spec it reflects and consumers
// can detect drift. It returns an empty string when root is not a git repository
// or git is unavailable.
func specVersion(root string) string {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
