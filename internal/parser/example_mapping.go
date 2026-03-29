package parser

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ExampleMapping struct {
	Story     string     `yaml:"story"`
	Rules     []Rule     `yaml:"rules"`
	Questions []Question `yaml:"questions"`
}

type Rule struct {
	ID       string    `yaml:"id"`
	Name     string    `yaml:"name"`
	Examples []Example `yaml:"examples"`
}

type Example struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type Question struct {
	ID   string `yaml:"id"`
	Text string `yaml:"text"`
}

func ParseExampleMapping(path string) (*ExampleMapping, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var em ExampleMapping
	if err := yaml.Unmarshal(data, &em); err != nil {
		return nil, err
	}

	return &em, nil
}
