package mcp

import "github.com/boykush/livt/internal/domain"

// versioned is embedded in every tool and resource payload so each result
// records the master version it reflects. Its json field is promoted to the top
// level of the result.
type versioned struct {
	SpecVersion string `json:"spec_version"`
}

// --- list_stories tool ---

type listStoriesInput struct{}

type listStoriesOutput struct {
	versioned
	Stories []storySummaryJSON `json:"stories"`
}

type storySummaryJSON struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	// ExampleMappingURI is the resource URI of the story's example mapping
	// (livt://mapping/{key}), present only when one exists. Its presence replaces
	// a boolean flag: absent means no mapping, present gives a handle the client
	// can read via resources/read.
	ExampleMappingURI string `json:"example_mapping_uri,omitempty"`
}

// --- resource payloads ---

// exampleMappingResult is the body of the livt://mapping/{story_key} resource.
type exampleMappingResult struct {
	versioned
	Mapping exampleMappingJSON `json:"mapping"`
}

// ruleResult is the body of the livt://mapping/{story_key}/rule/{rule_id} resource.
type ruleResult struct {
	versioned
	Rule ruleJSON `json:"rule"`
}

// --- JSON projections of the domain (kept stable; internal shapes like
// StoryKey{Value} are flattened to plain fields) ---

type exampleJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ruleJSON struct {
	ID string `json:"id"`
	// URI is the rule's own resource (livt://mapping/{story_key}/rule/{id}), so a
	// rule listed inside a mapping links to its addressable form.
	URI      string        `json:"uri"`
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

func toRuleJSON(storyKey string, r domain.Rule) ruleJSON {
	examples := make([]exampleJSON, 0, len(r.Examples))
	for _, e := range r.Examples {
		examples = append(examples, exampleJSON{ID: e.ID, Name: e.Name})
	}
	return ruleJSON{ID: r.ID, URI: ruleURI(storyKey, r.ID), Name: r.Name, Examples: examples}
}

func toExampleMappingJSON(em *domain.ExampleMapping) exampleMappingJSON {
	rules := make([]ruleJSON, 0, len(em.Rules))
	for _, r := range em.Rules {
		rules = append(rules, toRuleJSON(em.StoryKey.Value, r))
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
