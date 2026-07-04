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

// registerTools wires the discovery tools. The spec reads (story maps, stories,
// example mappings, rules, ubiquitous terms) are exposed as resources instead —
// see registerResources; the tools only list what exists and hand out URIs.
func (s *Server) registerTools(srv *mcpsdk.Server) {
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_stories",
		Description: "List all stories. Each entry links to its story resource (livt://story/{key}) and, when one exists, its example mapping resource (livt://mapping/{key}).",
	}, s.listStories)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_story_maps",
		Description: "List all story maps. Each entry links to its story map resource (livt://story-map/{map_name}).",
	}, s.listStoryMaps)
}

func (s *Server) listStories(_ context.Context, _ *mcpsdk.CallToolRequest, _ listStoriesInput) (*mcpsdk.CallToolResult, listStoriesOutput, error) {
	stories, err := s.cfg.stories()
	if err != nil {
		return nil, listStoriesOutput{}, err
	}
	return nil, listStoriesOutput{versioned: s.versioned(), Stories: stories}, nil
}

func (s *Server) listStoryMaps(_ context.Context, _ *mcpsdk.CallToolRequest, _ listStoryMapsInput) (*mcpsdk.CallToolResult, listStoryMapsOutput, error) {
	maps, err := s.cfg.storyMaps()
	if err != nil {
		return nil, listStoryMapsOutput{}, err
	}
	return nil, listStoryMapsOutput{versioned: s.versioned(), StoryMaps: maps}, nil
}

func (s *Server) versioned() versioned {
	return versioned{SpecVersion: specVersion(s.cfg.Root)}
}

// --- data access on Config (pure; shared by the tool and the resource handlers) ---

// exampleMapping loads the example mapping for storyKey. It distinguishes a
// missing mapping ("not found") from a malformed one (parse error). The key is
// validated first so an externally-supplied key (tool argument or resource URI)
// cannot escape the mappings directory.
func (c Config) exampleMapping(storyKey string) (*domain.ExampleMapping, error) {
	if !validSegment(storyKey) {
		return nil, fmt.Errorf("example mapping for story %q not found", storyKey)
	}
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

// stories lists every story, linking to its own story resource and to its
// example mapping resource when one exists. A missing stories directory yields
// an empty list, not an error.
func (c Config) stories() ([]storySummaryJSON, error) {
	all, err := parser.ParseAllStories(c.storiesDir())
	if err != nil {
		return nil, err
	}
	out := make([]storySummaryJSON, 0, len(all))
	for _, story := range all {
		summary := storySummaryJSON{Key: story.Key.Value, Name: story.Name, URI: storyURI(story.Key.Value)}
		if c.hasExampleMapping(story.Key.Value) {
			summary.ExampleMappingURI = mappingURI(story.Key.Value)
		}
		out = append(out, summary)
	}
	return out, nil
}

func (c Config) hasExampleMapping(storyKey string) bool {
	_, err := os.Stat(filepath.Join(c.mappingsDir(), storyKey+".yaml"))
	return err == nil
}

// storyMaps lists every story map with its resource URI. A missing usm
// directory yields an empty list, not an error.
func (c Config) storyMaps() ([]storyMapSummaryJSON, error) {
	all, err := parser.ParseAllStoryMaps(c.usmDir())
	if err != nil {
		return nil, err
	}
	out := make([]storyMapSummaryJSON, 0, len(all))
	for _, sm := range all {
		out = append(out, storyMapSummaryJSON{Name: sm.Name, URI: storyMapURI(sm.Name)})
	}
	return out, nil
}

// storyMap loads the story map with the given display name. Maps live in
// discoveries/usm/*.yaml but are addressed by name — the identifier the build
// output uses for story-map/{name}.html — so the lookup scans all maps. The
// name never touches the filesystem, so it needs no segment validation.
func (c Config) storyMap(name string) (*domain.StoryMap, error) {
	all, err := parser.ParseAllStoryMaps(c.usmDir())
	if err != nil {
		return nil, err
	}
	for _, sm := range all {
		if sm.Name == name {
			return sm, nil
		}
	}
	return nil, fmt.Errorf("story map %q not found", name)
}

// story loads a story by key. Unlike parser.FindStoryByKey it reports a
// missing story as not found instead of fabricating a placeholder.
func (c Config) story(storyKey string) (*domain.Story, error) {
	if !validSegment(storyKey) {
		return nil, fmt.Errorf("story %q not found", storyKey)
	}
	path := filepath.Join(c.storiesDir(), storyKey+".md")
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("story %q not found", storyKey)
	}
	story, err := parser.ParseStory(path)
	if err != nil {
		return nil, fmt.Errorf("parse story %q: %w", storyKey, err)
	}
	return story, nil
}

func (c Config) hasStory(storyKey string) bool {
	_, err := os.Stat(filepath.Join(c.storiesDir(), storyKey+".md"))
	return err == nil
}

// term loads a ubiquitous language term by key.
func (c Config) term(termKey string) (*domain.Term, error) {
	if !validSegment(termKey) {
		return nil, fmt.Errorf("term %q not found", termKey)
	}
	path := filepath.Join(c.ubiquitousDir(), termKey+".md")
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("term %q not found", termKey)
	}
	term, err := parser.ParseTerm(path)
	if err != nil {
		return nil, fmt.Errorf("parse term %q: %w", termKey, err)
	}
	return term, nil
}
