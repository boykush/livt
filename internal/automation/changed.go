package automation

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Changed reports whether the diff between two revisions adds or removes a
// line holding the marker — the only change that can alter what a report
// claims. A rename or a line that merely moved produces neither, and a deleted
// file shows its marker lines as removals, so one question covers both.
func Changed(root, base, head string) (bool, error) {
	out, err := exec.Command("git", "-C", root, "diff", "--name-only", "-G", Marker, base, head).Output()
	if err != nil {
		return false, fmt.Errorf("git diff %s %s: %w%s", base, head, err, gitMessage(err))
	}
	return strings.TrimSpace(string(out)) != "", nil
}

// gitMessage surfaces git's own words: "unknown revision" is the reason the
// caller needs, and the exec error alone does not carry it.
func gitMessage(err error) string {
	var exit *exec.ExitError
	if errors.As(err, &exit) && len(exit.Stderr) > 0 {
		return ": " + strings.TrimSpace(string(exit.Stderr))
	}
	return ""
}
