package cmd

import "testing"

func TestResolveRootPrecedence(t *testing.T) {
	t.Setenv("LIVT_ROOT", "/from/env")

	if got := resolveRoot("/from/flag"); got != "/from/flag" {
		t.Errorf("flag should win: got %q", got)
	}
	if got := resolveRoot(""); got != "/from/env" {
		t.Errorf("env should be used when no flag: got %q", got)
	}

	t.Setenv("LIVT_ROOT", "")
	if got := resolveRoot(""); got != "." {
		t.Errorf("default should be current dir: got %q", got)
	}
}
