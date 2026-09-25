// Package gitrev reads the livt repository's own revision. Two surfaces name
// it — the MCP payload's spec_version and the site's build page — and one
// reading keeps them from disagreeing about what the spec was.
package gitrev

import (
	"os/exec"
	"strings"
)

// Short returns the short HEAD of the repository at root, or an empty string
// when root is not a git repository or git is unavailable. Both callers treat
// that absence as "cannot say" rather than as an error: a livt repository
// works without git, and neither surface is worth failing a build over.
func Short(root string) string {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
