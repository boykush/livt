package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"gopkg.in/yaml.v3"
)

type exampleMappingYAML struct {
	Rules     []ruleYAML     `yaml:"rules"`
	Questions []questionYAML `yaml:"questions"`
}

type ruleYAML struct {
	ID       string        `yaml:"id"`
	Name     string        `yaml:"name"`
	Examples []exampleYAML `yaml:"examples"`
}

type exampleYAML struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type questionYAML struct {
	ID   string `yaml:"id"`
	Text string `yaml:"text"`
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

	storyKey := strings.TrimSuffix(filepath.Base(path), ".yaml")

	var rules []domain.Rule
	for _, r := range raw.Rules {
		var examples []domain.Example
		for _, ex := range r.Examples {
			examples = append(examples, domain.Example{ID: ex.ID, Name: ex.Name})
		}
		rules = append(rules, domain.Rule{ID: r.ID, Name: r.Name, Examples: examples})
	}

	var questions []domain.Question
	for _, q := range raw.Questions {
		questions = append(questions, domain.Question{ID: q.ID, Text: q.Text})
	}

	return &domain.ExampleMapping{
		Story:     storyKey,
		Rules:     rules,
		Questions: questions,
	}, nil
}
