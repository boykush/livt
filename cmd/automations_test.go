package cmd

import (
	"bytes"
	"strings"
	"testing"
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
