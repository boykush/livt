package domain

// Automation is one test claiming, in the implementation repository, that it
// automates this point of the spec. It is derived from the collected reports
// at build time and never written into a mapping: which tests cover what is
// the implementation's state, and the livt repository records decisions.
type Automation struct {
	Repo string
	Rev  string
	File string
	Line int
	// URL is where the line can be read on the forge, empty when the host
	// could not be identified. The link is decoration; its absence never
	// changes whether the citation counts.
	URL string
}

// Label names the citation the way a reader scans it — the repository and the
// file, without the path that only a machine follows.
func (a Automation) Label() string {
	name := a.File
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}
	if a.Repo == "" {
		return name
	}
	return shortRepo(a.Repo) + ": " + name
}

func shortRepo(repo string) string {
	for i := len(repo) - 1; i >= 0; i-- {
		if repo[i] == '/' {
			return repo[i+1:]
		}
	}
	return repo
}
