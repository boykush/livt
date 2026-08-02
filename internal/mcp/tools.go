package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerTools wires the discovery tools. The spec reads (story maps, stories,
// personas, example mappings, rules, ubiquitous terms) are exposed as resources
// instead — see registerResources; the tools only list what exists and hand out
// URIs.
func (s *Server) registerTools(srv *mcpsdk.Server) {
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_stories",
		Description: "List all stories. Each entry hands out the uris to read next: its story resource (livt://story/{key}), its example mapping resource (livt://mapping/{key}) when one exists, the persona it is written for (livt://persona/{persona_key}) when it names one, and the opportunities (story maps, livt://story-map/{map_name}) it sits on. Pass opportunity to list only the stories on that map.",
	}, s.listStories)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_story_maps",
		Description: "List all story maps. Each entry hands out the uri of its story map resource (livt://story-map/{map_name}).",
	}, s.listStoryMaps)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_terms",
		Description: "List the ubiquitous language — the team's agreed vocabulary — as every term's key, display name, and the uri of its term resource (livt://ubiquitous/{term_key}) to read the definition from. Look a word up here before naming things in code and tests, so the implementation speaks the domain's language. A term scoped to one context also carries ctx; the same key can name one term across contexts and another inside one, so the pair identifies it.",
	}, s.listTerms)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_personas",
		Description: "List the personas — the actors the stories are written for — as every persona's key, display name, and the uri of its persona resource (livt://persona/{persona_key}) to read the description from. A story names at most one persona; read it here to learn who a story serves before deciding what the behaviour should be.",
	}, s.listPersonas)
}

func (s *Server) listStories(_ context.Context, _ *mcpsdk.CallToolRequest, in listStoriesInput) (*mcpsdk.CallToolResult, listStoriesOutput, error) {
	stories, err := s.cfg.stories(in.Opportunity)
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

func (s *Server) listTerms(_ context.Context, _ *mcpsdk.CallToolRequest, _ listTermsInput) (*mcpsdk.CallToolResult, listTermsOutput, error) {
	terms, err := s.cfg.terms()
	if err != nil {
		return nil, listTermsOutput{}, err
	}
	return nil, listTermsOutput{versioned: s.versioned(), Terms: terms}, nil
}

func (s *Server) listPersonas(_ context.Context, _ *mcpsdk.CallToolRequest, _ listPersonasInput) (*mcpsdk.CallToolResult, listPersonasOutput, error) {
	personas, err := s.cfg.personas()
	if err != nil {
		return nil, listPersonasOutput{}, err
	}
	return nil, listPersonasOutput{versioned: s.versioned(), Personas: personas}, nil
}

func (s *Server) versioned() versioned {
	return versioned{SpecVersion: specVersion(s.cfg.Root)}
}

// --- data access on Config (pure; shared by the tool and the resource handlers) ---

// Retirement never turns a read into not-found: a rule, example, or question
// keeps resolving by its URI and carries retired so the caller can tell. The
// mapping keeps its retired entries too — it is the structural record their ids
// are numbered from, and dropping them would make a taken id look free. The
// list_* tools enumerate stories, story maps, terms, and personas, none of which
// have a retired concept.

// exampleMapping loads the example mapping for storyKey. It distinguishes a
// missing mapping ("not found") from a malformed one (parse error). The key is
// validated first so an externally-supplied key (tool argument or resource URI)
// cannot escape the mappings directory.
func (c Config) exampleMapping(storyKey string) (*domain.ExampleMapping, error) {
	if !uri.ValidSegment(storyKey) {
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

// example returns a single example by id from a rule. The lookup is scoped by
// rule because example ids are numbered within their rule — the same reason the
// example URI carries the rule id.
func (c Config) example(storyKey, ruleID, exampleID string) (domain.Example, error) {
	r, err := c.rule(storyKey, ruleID)
	if err != nil {
		return domain.Example{}, err
	}
	for _, e := range r.Examples {
		if e.ID == exampleID {
			return e, nil
		}
	}
	return domain.Example{}, fmt.Errorf("example %q not found in rule %q of story %q", exampleID, ruleID, storyKey)
}

// question returns a single question by id from the story's example mapping,
// retired or not. Questions hang off the mapping, not off a rule.
func (c Config) question(storyKey, questionID string) (domain.Question, error) {
	em, err := c.exampleMapping(storyKey)
	if err != nil {
		return domain.Question{}, err
	}
	for _, q := range em.Questions {
		if q.ID == questionID {
			return q, nil
		}
	}
	return domain.Question{}, fmt.Errorf("question %q not found in story %q", questionID, storyKey)
}

// stories lists every story, linking to its own story resource, to its example
// mapping resource when one exists, and to the opportunities (story maps) it
// sits on. A non-empty opportunity keeps only the stories on the map with that
// exact name; an unknown name yields an empty list, not an error — which
// opportunities exist is the livt repository's business, not the caller's. A missing
// stories directory yields an empty list, not an error.
func (c Config) stories(opportunity string) ([]storySummaryJSON, error) {
	all, err := parser.ParseAllStories(c.storiesDir())
	if err != nil {
		return nil, err
	}
	index, err := c.storyOpportunities()
	if err != nil {
		return nil, err
	}
	out := make([]storySummaryJSON, 0, len(all))
	for _, story := range all {
		refs := index[story.Key.Value]
		if opportunity != "" && !hasOpportunity(refs, opportunity) {
			continue
		}
		summary := storySummaryJSON{Key: story.Key.Value, Name: story.Name, Persona: c.toPersonaRef(story.Persona), URI: uri.Story(story.Key.Value), Opportunities: refs}
		if c.hasExampleMapping(story.Key.Value) {
			summary.ExampleMappingURI = uri.Mapping(story.Key.Value)
		}
		out = append(out, summary)
	}
	return out, nil
}

// hasOpportunity reports whether refs include a story map with this name.
func hasOpportunity(refs []opportunityRefJSON, name string) bool {
	for _, r := range refs {
		if r.Name == name {
			return true
		}
	}
	return false
}

// storyOpportunities indexes every story key to the opportunities (story maps)
// it sits on, one ref per map in map file order — the same story-to-maps
// derivation the site build uses for its opportunity chips. A key can recur
// across steps within one map and still gets a single ref for it.
func (c Config) storyOpportunities() (map[string][]opportunityRefJSON, error) {
	maps, err := parser.ParseAllStoryMaps(c.usmDir())
	if err != nil {
		return nil, err
	}
	index := make(map[string][]opportunityRefJSON)
	for _, sm := range maps {
		ref := opportunityRefJSON{Name: sm.Name, URI: uri.StoryMap(sm.Name)}
		seen := make(map[string]bool)
		for _, a := range sm.Activities {
			for _, st := range a.Steps {
				for _, sc := range st.Stories {
					if !sc.HasKey() || seen[sc.Key.Value] {
						continue
					}
					seen[sc.Key.Value] = true
					index[sc.Key.Value] = append(index[sc.Key.Value], ref)
				}
			}
		}
	}
	return index, nil
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
		out = append(out, storyMapSummaryJSON{Name: sm.Name, URI: uri.StoryMap(sm.Name)})
	}
	return out, nil
}

// terms lists every ubiquitous language term with its resource URI, including
// the ones no board references yet: the glossary is the livt repository's
// vocabulary, not a projection of what the boards happen to cite. Terms come
// out context-free first, then scoped, each group in filesystem order — the
// order the glossary page renders. A missing ubiquitous directory yields an
// empty list, not an error.
func (c Config) terms() ([]termSummaryJSON, error) {
	all, err := parser.ParseAllTerms(c.ubiquitousDir())
	if err != nil {
		return nil, err
	}
	out := make([]termSummaryJSON, 0, len(all))
	for _, term := range all {
		out = append(out, termSummaryJSON{Key: term.Key, Ctx: term.Ctx, Name: term.Name, URI: uri.Term(term.Ctx, term.Key)})
	}
	return out, nil
}

// personas lists every persona with its resource URI, including the ones no
// story names yet: the list is the livt repository's cast, not a projection of
// what the stories happen to cite. A missing personas directory yields an empty
// list, not an error.
func (c Config) personas() ([]personaSummaryJSON, error) {
	all, err := parser.ParseAllPersonas(c.personasDir())
	if err != nil {
		return nil, err
	}
	out := make([]personaSummaryJSON, 0, len(all))
	for _, p := range all {
		out = append(out, personaSummaryJSON{Key: p.Key, Name: p.Name, URI: uri.Persona(p.Key)})
	}
	return out, nil
}

// persona loads a persona by key.
func (c Config) persona(key string) (*domain.Persona, error) {
	if !uri.ValidSegment(key) {
		return nil, fmt.Errorf("persona %q not found", key)
	}
	if _, err := os.Stat(filepath.Join(c.personasDir(), key+".md")); err != nil {
		return nil, fmt.Errorf("persona %q not found", key)
	}
	persona, err := parser.ParsePersona(c.personasDir(), key)
	if err != nil {
		return nil, fmt.Errorf("parse persona %q: %w", key, err)
	}
	return persona, nil
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
	if !uri.ValidSegment(storyKey) {
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

// term loads a ubiquitous language term by reference — "{ctx}/{term-key}" for a
// term scoped to one context, "{term-key}" for one that holds across them. The
// two are distinct terms even when the key matches, so the reference is what is
// looked up rather than the key alone.
func (c Config) term(ref string) (*domain.Term, error) {
	ctx, key, ok := uri.SplitTermRef(ref)
	if !ok {
		return nil, fmt.Errorf("term %q not found", ref)
	}
	if _, err := os.Stat(filepath.Join(c.ubiquitousDir(), ctx, key+".md")); err != nil {
		return nil, fmt.Errorf("term %q not found", ref)
	}
	term, err := parser.ParseTerm(c.ubiquitousDir(), ref)
	if err != nil {
		return nil, fmt.Errorf("parse term %q: %w", ref, err)
	}
	return term, nil
}
