package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerTools wires the tool handlers. One tool per rule of the
// automate-from-master-in-impl-repos example mapping.
func (s *Server) registerTools(srv *mcpsdk.Server) {
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "get_example_mapping",
		Description: "Get the example mapping (rules, examples, questions, ubiquitous terms) for a story by its key.",
	}, s.getExampleMapping)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "get_rule",
		Description: "Get a single rule and its examples from a story's example mapping, by story key and rule id.",
	}, s.getRule)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_stories",
		Description: "List all stories, each flagged with whether it has an example mapping.",
	}, s.listStories)
}

// --- handlers (thin adapters: query the master, then stamp the spec version) ---

func (s *Server) getExampleMapping(_ context.Context, _ *mcpsdk.CallToolRequest, in getExampleMappingInput) (*mcpsdk.CallToolResult, getExampleMappingOutput, error) {
	em, err := s.cfg.exampleMapping(in.StoryKey)
	if err != nil {
		return nil, getExampleMappingOutput{}, err
	}
	return nil, getExampleMappingOutput{versioned: s.versioned(), Mapping: toExampleMappingJSON(em)}, nil
}

func (s *Server) getRule(_ context.Context, _ *mcpsdk.CallToolRequest, in getRuleInput) (*mcpsdk.CallToolResult, getRuleOutput, error) {
	rule, err := s.cfg.rule(in.StoryKey, in.RuleID)
	if err != nil {
		return nil, getRuleOutput{}, err
	}
	return nil, getRuleOutput{versioned: s.versioned(), Rule: toRuleJSON(rule)}, nil
}

func (s *Server) listStories(_ context.Context, _ *mcpsdk.CallToolRequest, _ listStoriesInput) (*mcpsdk.CallToolResult, listStoriesOutput, error) {
	stories, err := s.cfg.stories()
	if err != nil {
		return nil, listStoriesOutput{}, err
	}
	return nil, listStoriesOutput{versioned: s.versioned(), Stories: stories}, nil
}

func (s *Server) versioned() versioned {
	return versioned{SpecVersion: specVersion(s.cfg.Root)}
}

// --- data access on Config (pure; unit-tested directly with temp dirs) ---

// exampleMapping loads the example mapping for storyKey. It distinguishes a
// missing mapping ("not found") from a malformed one (parse error).
func (c Config) exampleMapping(storyKey string) (*domain.ExampleMapping, error) {
	path := filepath.Join(c.mappingsDir(), storyKey+".yaml")
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("example mapping for story %q not found", storyKey)
	}
	em, err := parser.ParseExampleMapping(path)
	if err != nil {
		return nil, fmt.Errorf("parse example mapping for story %q: %w", storyKey, err)
	}
	return em, nil
}

// rule returns a single rule by id from the story's example mapping.
func (c Config) rule(storyKey, ruleID string) (domain.Rule, error) {
	em, err := c.exampleMapping(storyKey)
	if err != nil {
		return domain.Rule{}, err
	}
	for _, r := range em.Rules {
		if r.ID == ruleID {
			return r, nil
		}
	}
	return domain.Rule{}, fmt.Errorf("rule %q not found in story %q", ruleID, storyKey)
}

// stories lists every story, flagging whether each has an example mapping. A
// missing stories directory yields an empty list, not an error.
func (c Config) stories() ([]storySummaryJSON, error) {
	all, err := parser.ParseAllStories(c.storiesDir())
	if err != nil {
		return nil, err
	}
	out := make([]storySummaryJSON, 0, len(all))
	for _, story := range all {
		out = append(out, storySummaryJSON{
			Key:               story.Key.Value,
			Name:              story.Name,
			HasExampleMapping: c.hasExampleMapping(story.Key.Value),
		})
	}
	return out, nil
}

func (c Config) hasExampleMapping(storyKey string) bool {
	_, err := os.Stat(filepath.Join(c.mappingsDir(), storyKey+".yaml"))
	return err == nil
}
