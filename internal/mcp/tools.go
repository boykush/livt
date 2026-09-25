package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/automation"
	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/boykush/livt/internal/gitrev"
)

// registerTools wires the discovery tools. The spec reads (opportunities and
// their canvases, story maps, stories, example mappings, rules, ubiquitous
// terms) are exposed as resources instead — see registerResources; the tools
// only list what exists and hand out URIs.
func (s *Server) registerTools(srv *mcpsdk.Server) {
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_opportunities",
		Description: "List the opportunities — what the product could take on, each a user problem together with the business benefit of solving it. Start here to see why a story map exists at all. Each entry hands out the uri of its opportunity resource (livt://opportunity/{key}), of its canvas (livt://opportunity-canvas/{key}) when one has been filled in, and its story_map (livt://story-map/{key}) when one has been drawn — an opportunity has one map, filed under its key. A missing canvas or story map is the record that the opportunity has not been taken that far, not an omission.",
	}, s.listOpportunities)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_stories",
		Description: "List all stories. Each entry hands out the uris to read next: its story resource (livt://story/{key}), its example mapping resource (livt://mapping/{key}) when one exists, and the opportunities it sits on, each with its key, name, and uri — livt://opportunity/{key}, or the story map's own when no opportunity file describes it and the map stands in. Pass opportunity, an opportunity key as list_opportunities hands it out, to list only the stories on its map.",
	}, s.listStories)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_example_mappings",
		Description: "List every example mapping, including one whose story has no card, which list_stories cannot reach. Each entry hands out the uri of its mapping resource (livt://mapping/{key}), its own name when it has one, story_uri (livt://story/{key}) when a story file exists — read that story for the name of a mapping with none of its own — and the opportunities its key sits on, as list_stories gives them. Pass opportunity, an opportunity key, to list only the mappings on its map.",
	}, s.listExampleMappings)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_story_maps",
		Description: "List all story maps, one per opportunity. Each entry hands out the opportunity_key the map is filed under, the map's name, and the uri of its story map resource (livt://story-map/{opportunity_key}).",
	}, s.listStoryMaps)
	mcpsdk.AddTool(srv, &mcpsdk.Tool{
		Name:        "list_terms",
		Description: "List the ubiquitous language — the team's agreed vocabulary — as every term's key, display name, and the uri of its term resource (livt://ubiquitous/{term_key}) to read the definition from. Look a word up here before naming things in code and tests, so the implementation speaks the domain's language. A term scoped to one context also carries ctx; the same key can name one term across contexts and another inside one, so the pair identifies it.",
	}, s.listTerms)
}

func (s *Server) listOpportunities(_ context.Context, _ *mcpsdk.CallToolRequest, _ listOpportunitiesInput) (*mcpsdk.CallToolResult, listOpportunitiesOutput, error) {
	opportunities, err := s.cfg.opportunities()
	if err != nil {
		return nil, listOpportunitiesOutput{}, err
	}
	return nil, listOpportunitiesOutput{versioned: s.versioned(), Opportunities: opportunities}, nil
}

func (s *Server) listStories(_ context.Context, _ *mcpsdk.CallToolRequest, in listStoriesInput) (*mcpsdk.CallToolResult, listStoriesOutput, error) {
	stories, err := s.cfg.stories(in.Opportunity)
	if err != nil {
		return nil, listStoriesOutput{}, err
	}
	return nil, listStoriesOutput{versioned: s.versioned(), Stories: stories}, nil
}

func (s *Server) listExampleMappings(_ context.Context, _ *mcpsdk.CallToolRequest, in listExampleMappingsInput) (*mcpsdk.CallToolResult, listExampleMappingsOutput, error) {
	mappings, err := s.cfg.exampleMappings(in.Opportunity)
	if err != nil {
		return nil, listExampleMappingsOutput{}, err
	}
	return nil, listExampleMappingsOutput{versioned: s.versioned(), ExampleMappings: mappings}, nil
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

func (s *Server) versioned() versioned {
	return versioned{SpecVersion: gitrev.Short(s.cfg.Root)}
}

// --- data access on Config (pure; shared by the tool and the resource handlers) ---

// Retirement never turns a read into not-found: a rule keeps resolving by its
// URI carrying status, an example or question carrying retired. The mapping
// keeps its closed entries too — it is the structural record their ids are
// numbered from, and dropping them would make a taken id look free. The list_*
// tools enumerate stories, story maps, and terms, which never close.

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
	idx, err := automation.Load(c.automationsDir())
	if err != nil {
		return nil, fmt.Errorf("read automations: %w", err)
	}
	idx.Attach(em)
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
// mapping resource when one exists, and to the opportunities it sits on. A
// non-empty opportunity keeps only the stories on the map with that key; an
// unknown key yields an empty list, not an error — which opportunities exist is
// the livt repository's business, not the caller's. A missing stories directory
// yields an empty list, not an error.
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
		if opportunity != "" && !sitsOn(refs, opportunity) {
			continue
		}
		summary := storySummaryJSON{Key: story.Key.Value, Name: story.Name, URI: uri.Story(story.Key.Value), Opportunities: refs}
		if c.hasExampleMapping(story.Key.Value) {
			summary.ExampleMappingURI = uri.Mapping(story.Key.Value)
		}
		out = append(out, summary)
	}
	return out, nil
}

// exampleMappings lists every mapping file, story or not — the one read that
// reaches a mapping with no story — filtered by opportunity the way stories
// are. A missing mappings directory yields an empty list, not an error.
func (c Config) exampleMappings(opportunity string) ([]exampleMappingSummaryJSON, error) {
	all, err := parser.ParseAllExampleMappings(c.mappingsDir())
	if err != nil {
		return nil, err
	}
	index, err := c.storyOpportunities()
	if err != nil {
		return nil, err
	}
	out := make([]exampleMappingSummaryJSON, 0, len(all))
	for _, em := range all {
		key := em.StoryKey.Value
		refs := index[key]
		if opportunity != "" && !sitsOn(refs, opportunity) {
			continue
		}
		summary := exampleMappingSummaryJSON{StoryKey: key, Name: em.Name, URI: uri.Mapping(key), Opportunities: refs}
		if c.hasStory(key) {
			summary.StoryURI = uri.Story(key)
		}
		out = append(out, summary)
	}
	return out, nil
}

// sitsOn reports whether refs include the opportunity with this key.
func sitsOn(refs []opportunityRefJSON, opportunityKey string) bool {
	for _, r := range refs {
		if r.Key == opportunityKey {
			return true
		}
	}
	return false
}

// storyOpportunities indexes every story key to the opportunities it sits on,
// one ref per map in map file order — the same derivation the site build uses
// for its opportunity chips, down to a map with no opportunity file standing in
// as its own. A key can recur across steps within one map and still gets a
// single ref for it.
func (c Config) storyOpportunities() (map[string][]opportunityRefJSON, error) {
	maps, err := parser.ParseAllStoryMaps(c.usmDir())
	if err != nil {
		return nil, err
	}
	opportunities, err := parser.ParseAllOpportunities(c.opportunitiesDir())
	if err != nil {
		return nil, err
	}
	described := make(map[string]*domain.Opportunity, len(opportunities))
	for _, o := range opportunities {
		described[o.Key.Value] = o
	}
	index := make(map[string][]opportunityRefJSON)
	for _, sm := range maps {
		key := sm.OpportunityKey.Value
		ref := opportunityRefJSON{Key: key, Name: sm.DisplayName(), URI: uri.StoryMap(key)}
		if o, ok := described[key]; ok {
			ref.Name, ref.URI = o.DisplayName(), uri.Opportunity(key)
		}
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

// opportunities lists every committed opportunity with its resource URI, the
// uri of its canvas when one has been filled in, and the story maps mapped for
// it. A missing opportunities directory yields an empty list, not an error.
func (c Config) opportunities() ([]opportunitySummaryJSON, error) {
	all, err := parser.ParseAllOpportunities(c.opportunitiesDir())
	if err != nil {
		return nil, err
	}
	out := make([]opportunitySummaryJSON, 0, len(all))
	for _, o := range all {
		summary := opportunitySummaryJSON{
			Key:       o.Key.Value,
			Name:      o.DisplayName(),
			URI:       uri.Opportunity(o.Key.Value),
			Statement: o.Body,
		}
		if c.hasOpportunityCanvas(o.Key.Value) {
			summary.CanvasURI = uri.OpportunityCanvas(o.Key.Value)
		}
		sm, err := c.storyMapFor(o.Key.Value)
		if err != nil {
			return nil, err
		}
		summary.StoryMap = sm
		out = append(out, summary)
	}
	return out, nil
}

// storyMapFor names the map filed under this opportunity's key, nil when none
// has been drawn. A malformed map is an error rather than nil: reading it as
// absent would misreport the opportunity as never taken on.
func (c Config) storyMapFor(opportunityKey string) (*storyMapSummaryJSON, error) {
	path := filepath.Join(c.usmDir(), opportunityKey+".yaml")
	if _, err := os.Stat(path); err != nil {
		return nil, nil
	}
	sm, err := parser.ParseStoryMap(path)
	if err != nil {
		return nil, fmt.Errorf("parse story map for %q: %w", opportunityKey, err)
	}
	summary := toStoryMapSummaryJSON(sm)
	return &summary, nil
}

func (c Config) hasOpportunityCanvas(opportunityKey string) bool {
	_, err := os.Stat(filepath.Join(c.canvasesDir(), opportunityKey+".yaml"))
	return err == nil
}

func (c Config) hasOpportunity(opportunityKey string) bool {
	_, err := os.Stat(filepath.Join(c.opportunitiesDir(), opportunityKey+".md"))
	return err == nil
}

// opportunity loads one opportunity by key. Like Config.story it reports a
// missing file as not found rather than fabricating a placeholder.
func (c Config) opportunity(opportunityKey string) (*domain.Opportunity, error) {
	if !uri.ValidSegment(opportunityKey) {
		return nil, fmt.Errorf("opportunity %q not found", opportunityKey)
	}
	path := filepath.Join(c.opportunitiesDir(), opportunityKey+".md")
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("opportunity %q not found", opportunityKey)
	}
	o, err := parser.ParseOpportunity(path)
	if err != nil {
		return nil, fmt.Errorf("parse opportunity %q: %w", opportunityKey, err)
	}
	return o, nil
}

// opportunityCanvas loads the canvas filled in for an opportunity. The canvas
// stands on its own: it resolves whether or not the opportunity file exists,
// the same way a mapping resolves for an uncommitted story.
func (c Config) opportunityCanvas(opportunityKey string) (*domain.OpportunityCanvas, error) {
	if !uri.ValidSegment(opportunityKey) {
		return nil, fmt.Errorf("opportunity canvas for %q not found", opportunityKey)
	}
	path := filepath.Join(c.canvasesDir(), opportunityKey+".yaml")
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("opportunity canvas for %q not found", opportunityKey)
	}
	canvas, err := parser.ParseOpportunityCanvas(path)
	if err != nil {
		return nil, fmt.Errorf("parse opportunity canvas for %q: %w", opportunityKey, err)
	}
	return canvas, nil
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
		out = append(out, toStoryMapSummaryJSON(sm))
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

// storyMap loads the story map drawn for an opportunity. It is filed under the
// opportunity's key, so it is read by that key the way the canvas beside it is,
// and the key is validated first so it cannot escape the usm directory.
func (c Config) storyMap(opportunityKey string) (*domain.StoryMap, error) {
	if !uri.ValidSegment(opportunityKey) {
		return nil, fmt.Errorf("story map for %q not found", opportunityKey)
	}
	path := filepath.Join(c.usmDir(), opportunityKey+".yaml")
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("story map for %q not found", opportunityKey)
	}
	sm, err := parser.ParseStoryMap(path)
	if err != nil {
		return nil, fmt.Errorf("parse story map for %q: %w", opportunityKey, err)
	}
	return sm, nil
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
