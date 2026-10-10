// Package agreement checks that a rule naming the people whose agreement it
// needs is only marked accepted once they have agreed, and keeps the record of
// who did. It reads files and git and nothing else: what a forge knows about a
// change arrives as Facts, written by whatever speaks to that forge.
package agreement

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Dir is where the records live, beside the collected reports and apart from
// the mappings people write.
const Dir = "agreements"

// Facts is what a forge knows about the change being checked, in a shape no
// forge owns. Author is empty when no person opened the change.
type Facts struct {
	Change     string   `json:"change"`
	Author     string   `json:"author"`
	ApprovedBy []string `json:"approved_by"`
}

// LoadFacts reads the facts a forge adapter wrote. No path means nothing is
// known, which is a state the check has an answer for.
func LoadFacts(path string) (Facts, error) {
	if path == "" {
		return Facts{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Facts{}, err
	}
	var f Facts
	if err := json.Unmarshal(data, &f); err != nil {
		return Facts{}, err
	}
	return f, nil
}

// The two ways a person stands behind a rule.
const (
	AsAuthor   = "author"
	AsApprover = "approver"
)

// Record is one person's agreement to one rule, and the change it was given in.
type Record struct {
	Who string `json:"who"`
	As  string `json:"as"`
	In  string `json:"in"`
}

// Entry is everyone on record as agreeing to one rule.
type Entry struct {
	URI        string   `json:"uri"`
	AcceptedBy []Record `json:"accepted_by"`
}

// File is one example mapping's record.
type File struct {
	Agreements []Entry `json:"agreements"`
}

// records returns who is on record for the rule.
func (f File) records(ruleURI string) []Record {
	for _, e := range f.Agreements {
		if e.URI == ruleURI {
			return e.AcceptedBy
		}
	}
	return nil
}

func recordPath(root, storyKey string) string {
	return filepath.Join(root, Dir, storyKey+".json")
}

func parseFile(data []byte) (File, error) {
	var f File
	if len(data) == 0 {
		return f, nil
	}
	err := json.Unmarshal(data, &f)
	return f, err
}

// readFile reads a mapping's record from the working tree. A mapping nobody
// has agreed anything on has no file, and that is not an error.
func readFile(root, storyKey string) (File, error) {
	data, err := os.ReadFile(recordPath(root, storyKey))
	if errors.Is(err, fs.ErrNotExist) {
		return File{}, nil
	}
	if err != nil {
		return File{}, err
	}
	return parseFile(data)
}

// same compares two people as a forge does: a login is not case-sensitive.
func same(a, b string) bool { return a != "" && strings.EqualFold(a, b) }

func contains(people []string, who string) bool {
	for _, p := range people {
		if same(p, who) {
			return true
		}
	}
	return false
}
