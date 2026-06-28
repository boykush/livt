package mcp

import "github.com/boykush/livt/internal/domain"

// versioned is embedded in every tool output so each result records the master
// version it reflects. Its json field is promoted to the top level of the result.
type versioned struct {
	SpecVersion string `json:"spec_version"`
}

// --- tool inputs ---

type getExampleMappingInput struct {
	StoryKey string `json:"story_key" jsonschema:"the story key identifying the example mapping (its YAML filename without extension)"`
}

type getRuleInput struct {
	StoryKey string `json:"story_key" jsonschema:"the story key identifying the example mapping"`
	RuleID   string `json:"rule_id" jsonschema:"the rule id to fetch, e.g. R-01"`
}

type listStoriesInput struct{}

// --- tool outputs ---

type getExampleMappingOutput struct {
	versioned
	Mapping exampleMappingJSON `json:"mapping"`
}

type getRuleOutput struct {
	versioned
	Rule ruleJSON `json:"rule"`
}

type listStoriesOutput struct {
	versioned
	Stories []storySummaryJSON `json:"stories"`
}

// --- JSON projections of the domain (kept stable; internal shapes like
// StoryKey{Value} are flattened to plain fields) ---

type exampleJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ruleJSON struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Examples []exampleJSON `json:"examples,omitempty"`
}

type questionJSON struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type exampleMappingJSON struct {
	StoryKey   string         `json:"story_key"`
	Rules      []ruleJSON     `json:"rules,omitempty"`
	Questions  []questionJSON `json:"questions,omitempty"`
	Ubiquitous []string       `json:"ubiquitous,omitempty"`
}

type storySummaryJSON struct {
	Key               string `json:"key"`
	Name              string `json:"name"`
	HasExampleMapping bool   `json:"has_example_mapping"`
}

func toRuleJSON(r domain.Rule) ruleJSON {
	examples := make([]exampleJSON, 0, len(r.Examples))
	for _, e := range r.Examples {
		examples = append(examples, exampleJSON{ID: e.ID, Name: e.Name})
	}
	return ruleJSON{ID: r.ID, Name: r.Name, Examples: examples}
}

func toExampleMappingJSON(em *domain.ExampleMapping) exampleMappingJSON {
	rules := make([]ruleJSON, 0, len(em.Rules))
	for _, r := range em.Rules {
		rules = append(rules, toRuleJSON(r))
	}
	questions := make([]questionJSON, 0, len(em.Questions))
	for _, q := range em.Questions {
		questions = append(questions, questionJSON{ID: q.ID, Text: q.Text})
	}
	return exampleMappingJSON{
		StoryKey:   em.StoryKey.Value,
		Rules:      rules,
		Questions:  questions,
		Ubiquitous: em.Ubiquitous,
	}
}
