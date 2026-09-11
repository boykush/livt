package parser

import (
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
	// Retired is a field rather than a commented-out block so a structural edit
	// of the YAML cannot drop it, and so the id stays visible to whoever numbers
	// the next one. Omitted means live.
	Retired bool `yaml:"retired"`
	// SupersededBy carries the retirement's other half: where the spec went.
	// It is livt URIs rather than bare ids so a successor in another mapping is
	// sayable, and it is a list so a rule that split into two can name both.
	SupersededBy []string `yaml:"superseded_by"`
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
		var examples []domain.Example
		for _, ex := range r.Examples {
			examples = append(examples, domain.Example{ID: ex.ID, Name: ex.Name, Retired: ex.Retired, SupersededBy: ex.SupersededBy})
		}
		rules = append(rules, domain.Rule{ID: r.ID, Name: r.Name, Examples: examples, Issues: r.Issues, Automated: r.Automated, Retired: r.Retired, SupersededBy: r.SupersededBy})
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
