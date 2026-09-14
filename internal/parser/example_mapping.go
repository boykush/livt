package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"gopkg.in/yaml.v3"
)

type exampleMappingYAML struct {
	Rules      []ruleYAML     `yaml:"rules"`
	Questions  []questionYAML `yaml:"questions"`
	Ubiquitous []string       `yaml:"ubiquitous"`
}

type ruleYAML struct {
	ID        string        `yaml:"id"`
	Name      string        `yaml:"name"`
	Examples  []exampleYAML `yaml:"examples"`
	Issues    []string      `yaml:"issues"`
	Automated bool          `yaml:"automated"`
	// Status reads as an ADR's does: proposed while the rule awaits agreement,
	// accepted once it has it, rejected or retired once it has closed. Omitted
	// means accepted.
	Status string `yaml:"status"`
	// Retired is the superseded spelling of the two closed statuses, still read
	// so a livt repository written against it keeps building. It is folded onto
	// Status at parse time and will stop being accepted.
	Retired bool `yaml:"retired"`
	// SupersededBy carries the retirement's other half: where the spec went.
	// It is livt URIs rather than bare ids so a successor in another mapping is
	// sayable, and it is a list so a rule that split into two can name both.
	SupersededBy []string `yaml:"superseded_by"`
}

// status reads the rule's status, defaulting an omitted one. An unknown
// value fails the parse: read as the default, a mistyped "proposed" would put
// an unagreed rule on the board as spec, with nothing there to say so.
func (r ruleYAML) status() (domain.RuleStatus, error) {
	status := domain.RuleStatusDefault
	if r.Status != "" {
		status = domain.RuleStatus(r.Status)
		if !status.Valid() {
			return "", fmt.Errorf("rule %q: unknown status %q (supported: %s)", r.ID, r.Status, domain.RuleStatusList())
		}
	}
	if r.Retired {
		return closed(status), nil
	}
	return status, nil
}

// closed folds the superseded retired: flag onto the status it was written
// beside. A proposal closed that way was turned down; anything else was spec
// that stopped holding — the distinction the flag left to whoever read both
// lines together.
func closed(status domain.RuleStatus) domain.RuleStatus {
	switch status {
	case domain.RuleProposed:
		return domain.RuleRejected
	case domain.RuleRejected, domain.RuleRetired:
		return status
	default:
		return domain.RuleRetired
	}
}

type exampleYAML struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Retired      bool     `yaml:"retired"`
	SupersededBy []string `yaml:"superseded_by"`
}

type questionYAML struct {
	ID           string   `yaml:"id"`
	Text         string   `yaml:"text"`
	Retired      bool     `yaml:"retired"`
	SupersededBy []string `yaml:"superseded_by"`
}

func ParseExampleMapping(path string) (*domain.ExampleMapping, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw exampleMappingYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	storyKey := domain.StoryKey{Value: strings.TrimSuffix(filepath.Base(path), ".yaml")}

	var rules []domain.Rule
	for _, r := range raw.Rules {
		status, err := r.status()
		if err != nil {
			return nil, err
		}
		var examples []domain.Example
		for _, ex := range r.Examples {
			examples = append(examples, domain.Example{ID: ex.ID, Name: ex.Name, Retired: ex.Retired, SupersededBy: ex.SupersededBy})
		}
		rules = append(rules, domain.Rule{ID: r.ID, Name: r.Name, Examples: examples, Status: status, Issues: r.Issues, Automated: r.Automated, SupersededBy: r.SupersededBy})
	}

	var questions []domain.Question
	for _, q := range raw.Questions {
		questions = append(questions, domain.Question{ID: q.ID, Text: q.Text, Retired: q.Retired, SupersededBy: q.SupersededBy})
	}

	return &domain.ExampleMapping{
		StoryKey:   storyKey,
		Rules:      rules,
		Questions:  questions,
		Ubiquitous: raw.Ubiquitous,
	}, nil
}

func ParseAllExampleMappings(mappingsDir string) ([]*domain.ExampleMapping, error) {
	files, err := filepath.Glob(filepath.Join(mappingsDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	var mappings []*domain.ExampleMapping
	for _, f := range files {
		em, err := ParseExampleMapping(f)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, em)
	}

	return mappings, nil
}
