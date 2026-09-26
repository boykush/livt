package automation

// livt:automates livt://mapping/collect-automations/rule/R-06

import "testing"

// livt:automates livt://mapping/collect-automations/rule/R-06/example/EX-01
func TestURLTemplateKnowsTheTwoHostedForges(t *testing.T) {
	for _, tc := range []struct {
		host string
		want string
	}{
		{"github.com", githubURLTemplate},
		{"gitlab.com", gitlabURLTemplate},
		{"git.acme.com", ""},
	} {
		got, err := URLTemplate("", "", Origin{Host: tc.host})
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.host, got, tc.want)
		}
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-06/example/EX-02
func TestURLTemplateTakesADeclarationForAnythingElse(t *testing.T) {
	custom := "{base}/src/commit/{rev}/{path}#L{line}"
	got, err := URLTemplate("", custom, Origin{Host: "git.acme.com"})
	if err != nil {
		t.Fatal(err)
	}
	if got != custom {
		t.Errorf("got %q, want the declared template", got)
	}

	got, err = URLTemplate(ForgeGitLab, "", Origin{Host: "git.acme.com"})
	if err != nil {
		t.Fatal(err)
	}
	if got != gitlabURLTemplate {
		t.Errorf("got %q, want the named forge's template", got)
	}

	if _, err := URLTemplate("gitea", "", Origin{}); err == nil {
		t.Error("an unknown forge name was accepted")
	}
}

func TestParseRemoteReducesEveryCheckoutShapeToTheSameBase(t *testing.T) {
	for _, remote := range []string{
		"https://github.com/acme/backend.git",
		"https://github.com/acme/backend",
		"git@github.com:acme/backend.git",
		"ssh://git@github.com/acme/backend.git",
	} {
		base, host, repo := parseRemote(remote)
		if base != "https://github.com/acme/backend" || host != "github.com" || repo != "acme/backend" {
			t.Errorf("%s: got %q %q %q", remote, base, host, repo)
		}
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-04/example/EX-04
func TestCleanDisregardsOnlyTheReportBeingWritten(t *testing.T) {
	report := "automations/acme/impl.json"
	if !clean(" M "+report+"\n", []string{report}) {
		t.Error("the scan could not name a revision because of its own output")
	}
	if clean(" M internal/thing.go\n M "+report+"\n", []string{report}) {
		t.Error("a scanned file changed and the revision was claimed anyway")
	}
	if !clean("", nil) {
		t.Error("an unchanged tree was read as dirty")
	}
}
