package automation

import (
	"fmt"
	"os/exec"
	"strings"
)

// Forge names a code host livt knows how to build a line URL for. Two is the
// whole table on purpose: a self-hosted GitLab and a self-hosted Gitea are the
// same hostname to us, so anything else declares its template instead of being
// guessed at.
const (
	ForgeGitHub = "github"
	ForgeGitLab = "gitlab"
)

// URL templates take {base}, {rev}, {path} and {line}.
const (
	githubURLTemplate = "{base}/blob/{rev}/{path}#L{line}"
	gitlabURLTemplate = "{base}/-/blob/{rev}/{path}#L{line}"
)

// Origin is what a checkout says about itself: which repository it is, at
// which revision, and where its lines can be read. Every field is optional —
// a scan of a directory that is not a git repository still reports citations.
type Origin struct {
	Repo string
	Rev  string
	Base string
	Host string
}

// ReadOrigin asks git, so a scan needs no arguments in the common case.
func ReadOrigin(root string) Origin {
	o := Origin{Rev: gitOutput(root, "rev-parse", "HEAD")}
	remote := gitOutput(root, "remote", "get-url", "origin")
	if remote == "" {
		return o
	}
	o.Base, o.Host, o.Repo = parseRemote(remote)
	return o
}

func gitOutput(root string, args ...string) string {
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// parseRemote reduces the shapes a remote can take — https, scp-like ssh,
// ssh:// — to the browsable base, since the form a checkout happens to use is
// not something the person writing a test chose.
func parseRemote(remote string) (base, host, repo string) {
	s := strings.TrimSuffix(strings.TrimSpace(remote), ".git")
	switch {
	case strings.HasPrefix(s, "git@"):
		s = strings.Replace(strings.TrimPrefix(s, "git@"), ":", "/", 1)
	case strings.HasPrefix(s, "ssh://git@"):
		s = strings.TrimPrefix(s, "ssh://git@")
	case strings.HasPrefix(s, "https://"):
		s = strings.TrimPrefix(s, "https://")
	case strings.HasPrefix(s, "http://"):
		s = strings.TrimPrefix(s, "http://")
	default:
		return "", "", ""
	}
	parts := strings.Split(s, "/")
	if len(parts) < 3 {
		return "", "", ""
	}
	host = parts[0]
	repo = strings.Join(parts[1:], "/")
	return "https://" + host + "/" + repo, host, repo
}

// URLTemplate resolves the template to build line URLs with: the one named, or
// the one the host implies. An empty result is a supported outcome — the link
// is decoration, and a citation counts without it.
func URLTemplate(forge, custom string, o Origin) (string, error) {
	if custom != "" {
		return custom, nil
	}
	switch forge {
	case ForgeGitHub:
		return githubURLTemplate, nil
	case ForgeGitLab:
		return gitlabURLTemplate, nil
	case "":
	default:
		return "", fmt.Errorf("unknown forge %q (known: %s, %s; otherwise pass --url-template)", forge, ForgeGitHub, ForgeGitLab)
	}
	switch o.Host {
	case "github.com":
		return githubURLTemplate, nil
	case "gitlab.com":
		return gitlabURLTemplate, nil
	}
	return "", nil
}
