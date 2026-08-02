package mcp

import (
	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/uri"
)

// versioned is embedded in every tool and resource payload so each result
// records the livt repository version it reflects. Its json field is promoted to the top
// level of the result.
type versioned struct {
	SpecVersion string `json:"spec_version"`
}

// --- list_stories tool ---

type listStoriesInput struct {
	// Opportunity narrows the list to the stories on one story map, matched by
	// exact display name — in livt a story map is one opportunity, so the map
	// name doubles as the filter axis. Unknown names yield an empty list.
	Opportunity string `json:"opportunity,omitempty" jsonschema:"Keep only stories on this opportunity, matched exactly against a story map display name (one story map is one opportunity). Unknown names yield an empty list; omit to list every story."`
}

type listStoriesOutput struct {
	versioned
	Stories []storySummaryJSON `json:"stories"`
}

type storySummaryJSON struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	// Persona is the actor the story is written for, absent when it names none.
	Persona *personaRefJSON `json:"persona,omitempty"`
	// URI is the story's own resource (livt://story/{key}): its name, body, and
	// frontmatter meta. Always present, since listed stories come from files.
	URI string `json:"uri"`
	// ExampleMappingURI is the resource URI of the story's example mapping
	// (livt://mapping/{key}), present only when one exists. Its presence replaces
	// a boolean flag: absent means no mapping, present gives a handle the client
	// can read via resources/read.
	ExampleMappingURI string `json:"example_mapping_uri,omitempty"`
	// Opportunities are the story maps this story sits on, one ref per map in
	// map order; empty for a story on no map.
	Opportunities []opportunityRefJSON `json:"opportunities,omitempty"`
}

// --- list_story_maps tool ---

type listStoryMapsInput struct{}

type listStoryMapsOutput struct {
	versioned
	StoryMaps []storyMapSummaryJSON `json:"story_maps"`
}

type storyMapSummaryJSON struct {
	Name string `json:"name"`
	// URI is the story map's resource (livt://story-map/{map_name}).
	URI string `json:"uri"`
}

// --- list_terms tool ---

type listTermsInput struct{}

type listTermsOutput struct {
	versioned
	Terms []termSummaryJSON `json:"terms"`
}

// termSummaryJSON is one row of the glossary listing. It differs from
// termRefJSON in what it promises: a ref is what a board authored and may name
// a term nobody committed, while every entry here comes from a term file, so
// Name and URI are always filled.
type termSummaryJSON struct {
	Key string `json:"key"`
	// Ctx as on termJSON: the context the term belongs to, omitted when it holds
	// across them. A caller citing a term composes its uri from the pair.
	Ctx  string `json:"ctx,omitempty"`
	Name string `json:"name"`
	// URI is the term's own resource — livt://ubiquitous/{ctx}/{key} for a
	// scoped term, livt://ubiquitous/{key} for one that holds across contexts.
	URI string `json:"uri"`
}

// --- list_personas tool ---

type listPersonasInput struct{}

type listPersonasOutput struct {
	versioned
	Personas []personaSummaryJSON `json:"personas"`
}

// personaSummaryJSON is one row of the persona listing. As with termSummaryJSON,
// every entry comes from a persona file, so Name and URI are always filled —
// unlike personaRefJSON, which carries whatever key a story authored.
type personaSummaryJSON struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	// URI is the persona's own resource (livt://persona/{key}).
	URI string `json:"uri"`
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

// exampleResult is the body of the
// livt://mapping/{story_key}/rule/{rule_id}/example/{example_id} resource.
type exampleResult struct {
	versioned
	Example exampleJSON `json:"example"`
}

// questionResult is the body of the livt://mapping/{story_key}/question/{question_id} resource.
type questionResult struct {
	versioned
	Question questionJSON `json:"question"`
}

// storyMapResult is the body of the livt://story-map/{map_name} resource.
type storyMapResult struct {
	versioned
	StoryMap storyMapJSON `json:"story_map"`
}

// storyResult is the body of the livt://story/{story_key} resource.
type storyResult struct {
	versioned
	Story storyJSON `json:"story"`
}

// personaResult is the body of the livt://persona/{persona_key} resource.
type personaResult struct {
	versioned
	Persona personaJSON `json:"persona"`
}

// termResult is the body of the livt://ubiquitous/{term_key} resource.
type termResult struct {
	versioned
	Term termJSON `json:"term"`
}

// --- JSON projections of the domain (kept stable; internal shapes like
// StoryKey{Value} are flattened to plain fields) ---

type exampleJSON struct {
	ID string `json:"id"`
	// URI is the example's own resource
	// (livt://mapping/{story_key}/rule/{rule_id}/example/{id}). Example ids are
	// numbered within their rule, so the address carries the rule.
	URI  string `json:"uri"`
	Name string `json:"name"`
	// Retired as on ruleJSON.
	Retired bool `json:"retired,omitempty"`
}

type ruleJSON struct {
	ID string `json:"id"`
	// URI is the rule's own resource (livt://mapping/{story_key}/rule/{id}), so a
	// rule listed inside a mapping links to its addressable form.
	URI      string        `json:"uri"`
	Name     string        `json:"name"`
	Examples []exampleJSON `json:"examples,omitempty"`
	// Issues are the rule's automation Issue URLs as recorded on the livt repository.
	Issues []string `json:"issues,omitempty"`
	// Automated is always present: consumers read the recorded judgment
	// without distinguishing absent from false.
	Automated bool `json:"automated"`
	// Retired says the rule is no longer part of the spec. Unlike Automated it
	// is omitted when false: it marks the exception, and every live rule
	// carrying "retired": false would drown the flag in noise.
	Retired bool `json:"retired,omitempty"`
}

type questionJSON struct {
	ID string `json:"id"`
	// URI is the question's own resource
	// (livt://mapping/{story_key}/question/{id}).
	URI  string `json:"uri"`
	Text string `json:"text"`
	// Retired as on ruleJSON; a retired question is no longer open.
	Retired bool `json:"retired,omitempty"`
}

type exampleMappingJSON struct {
	StoryKey  string         `json:"story_key"`
	Rules     []ruleJSON     `json:"rules,omitempty"`
	Questions []questionJSON `json:"questions,omitempty"`
	// Ubiquitous keeps the raw term keys as authored in the mapping;
	// UbiquitousTerms carries the same keys resolved to term resources.
	Ubiquitous      []string      `json:"ubiquitous,omitempty"`
	UbiquitousTerms []termRefJSON `json:"ubiquitous_terms,omitempty"`
}

type storyMapJSON struct {
	Name       string         `json:"name"`
	Activities []activityJSON `json:"activities,omitempty"`
	Releases   []releaseJSON  `json:"releases,omitempty"`
	// Personas keeps the raw keys the map declares; PersonaRefs carries the same
	// keys resolved to persona resources. They are the actors whose journey this
	// map covers — read them before the activities, since a step only makes sense
	// as somebody's step.
	Personas    []string         `json:"personas,omitempty"`
	PersonaRefs []personaRefJSON `json:"persona_refs,omitempty"`
	// Ubiquitous and UbiquitousTerms mirror the same pair on exampleMappingJSON.
	Ubiquitous      []string      `json:"ubiquitous,omitempty"`
	UbiquitousTerms []termRefJSON `json:"ubiquitous_terms,omitempty"`
}

type activityJSON struct {
	Key   string     `json:"key,omitempty"`
	Name  string     `json:"name"`
	Steps []stepJSON `json:"steps,omitempty"`
}

type stepJSON struct {
	Key     string          `json:"key,omitempty"`
	Name    string          `json:"name"`
	Stories []storyCardJSON `json:"stories,omitempty"`
}

type storyCardJSON struct {
	Key     string `json:"key,omitempty"`
	Name    string `json:"name"`
	Release string `json:"release,omitempty"`
	// URI is the story resource (livt://story/{key}), present when the card has a
	// key and the story file is committed; a bare card is still a candidate.
	URI string `json:"uri,omitempty"`
}

type releaseJSON struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type storyJSON struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	// Persona as on storySummaryJSON. It is the structural form of the body's
	// "as a ..." line, so a consumer reads who the story serves without parsing
	// prose.
	Persona *personaRefJSON `json:"persona,omitempty"`
	Body    string          `json:"body"`
	// Meta carries the story's frontmatter fields beyond name (e.g. issue), in
	// source order.
	Meta []metaFieldJSON `json:"meta,omitempty"`
	// ExampleMappingURI as on storySummaryJSON: present only when a mapping exists.
	ExampleMappingURI string `json:"example_mapping_uri,omitempty"`
	// Opportunities as on storySummaryJSON: the story maps this story sits on.
	Opportunities []opportunityRefJSON `json:"opportunities,omitempty"`
}

type metaFieldJSON struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type personaJSON struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Body string `json:"body"`
}

// personaRefJSON points at the persona a story is written for. Name and URI are
// filled when the persona file exists; a bare key means the persona is not
// committed yet, mirroring termRefJSON and how the site renders an unresolved
// persona chip.
type personaRefJSON struct {
	Key  string `json:"key"`
	Name string `json:"name,omitempty"`
	URI  string `json:"uri,omitempty"`
}

type termJSON struct {
	Key string `json:"key"`
	// Ctx is the context the term belongs to, omitted when it holds across them.
	// Key alone does not identify a term once contexts are in play — the same key
	// can sit at the root and under a context — so a consumer citing this term
	// composes its uri from the pair.
	Ctx  string `json:"ctx,omitempty"`
	Name string `json:"name"`
	Body string `json:"body"`
}

// termRefJSON points at a ubiquitous language term referenced from a mapping or
// story map. Key and Ctx are the reference taken apart — Ctx empty for a term
// that holds across contexts. Name and URI are filled when the term file exists;
// a bare key means the term is not committed yet (mirroring how the site renders
// unresolved term cards).
type termRefJSON struct {
	Key  string `json:"key"`
	Ctx  string `json:"ctx,omitempty"`
	Name string `json:"name,omitempty"`
	URI  string `json:"uri,omitempty"`
}

// opportunityRefJSON points at one story map a story belongs to. In livt a
// story map is one opportunity, so Name is the opportunity name (= map name)
// and URI is the story map resource (livt://story-map/{map_name}).
type opportunityRefJSON struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

func toRuleJSON(storyKey string, r domain.Rule) ruleJSON {
	examples := make([]exampleJSON, 0, len(r.Examples))
	for _, e := range r.Examples {
		examples = append(examples, toExampleJSON(storyKey, r.ID, e))
	}
	return ruleJSON{ID: r.ID, URI: uri.Rule(storyKey, r.ID), Name: r.Name, Examples: examples, Issues: r.Issues, Automated: r.Automated, Retired: r.Retired}
}

func toExampleJSON(storyKey, ruleID string, e domain.Example) exampleJSON {
	return exampleJSON{ID: e.ID, URI: uri.Example(storyKey, ruleID, e.ID), Name: e.Name, Retired: e.Retired}
}

func toQuestionJSON(storyKey string, q domain.Question) questionJSON {
	return questionJSON{ID: q.ID, URI: uri.Question(storyKey, q.ID), Text: q.Text, Retired: q.Retired}
}

// toExampleMappingJSON is a Config method because resolving the referenced
// ubiquitous terms reads the livt repository's ubiquitous directory. Retired rules and
// questions stay in the projection, flagged — see the note in tools.go.
func (c Config) toExampleMappingJSON(em *domain.ExampleMapping) exampleMappingJSON {
	rules := make([]ruleJSON, 0, len(em.Rules))
	for _, r := range em.Rules {
		rules = append(rules, toRuleJSON(em.StoryKey.Value, r))
	}
	questions := make([]questionJSON, 0, len(em.Questions))
	for _, q := range em.Questions {
		questions = append(questions, toQuestionJSON(em.StoryKey.Value, q))
	}
	return exampleMappingJSON{
		StoryKey:        em.StoryKey.Value,
		Rules:           rules,
		Questions:       questions,
		Ubiquitous:      em.Ubiquitous,
		UbiquitousTerms: c.toTermRefs(em.Ubiquitous),
	}
}

func (c Config) toStoryMapJSON(sm *domain.StoryMap) storyMapJSON {
	activities := make([]activityJSON, 0, len(sm.Activities))
	for _, a := range sm.Activities {
		steps := make([]stepJSON, 0, len(a.Steps))
		for _, st := range a.Steps {
			stories := make([]storyCardJSON, 0, len(st.Stories))
			for _, sc := range st.Stories {
				card := storyCardJSON{Key: sc.Key.Value, Name: sc.Name, Release: sc.Release}
				if sc.HasKey() && c.hasStory(sc.Key.Value) {
					card.URI = uri.Story(sc.Key.Value)
				}
				stories = append(stories, card)
			}
			steps = append(steps, stepJSON{Key: st.Key, Name: st.Name, Stories: stories})
		}
		activities = append(activities, activityJSON{Key: a.Key, Name: a.Name, Steps: steps})
	}
	releases := make([]releaseJSON, 0, len(sm.Releases))
	for _, r := range sm.Releases {
		releases = append(releases, releaseJSON{ID: r.ID, Name: r.Name})
	}
	return storyMapJSON{
		Name:            sm.Name,
		Activities:      activities,
		Releases:        releases,
		Personas:        sm.Personas,
		PersonaRefs:     c.toPersonaRefs(sm.Personas),
		Ubiquitous:      sm.Ubiquitous,
		UbiquitousTerms: c.toTermRefs(sm.Ubiquitous),
	}
}

// toStoryJSON is a Config method because linking the example mapping and
// resolving the story's opportunities read the livt repository's mappings and usm
// directories. Unlike term refs, a malformed story map is an error, not a bare
// ref: silently dropping opportunities would misreport the story as unattached.
func (c Config) toStoryJSON(story *domain.Story) (storyJSON, error) {
	meta := make([]metaFieldJSON, 0, len(story.Meta))
	for _, m := range story.Meta {
		meta = append(meta, metaFieldJSON{Key: m.Key, Value: m.Value})
	}
	out := storyJSON{Key: story.Key.Value, Name: story.Name, Persona: c.toPersonaRef(story.Persona), Body: story.Body, Meta: meta}
	if c.hasExampleMapping(story.Key.Value) {
		out.ExampleMappingURI = uri.Mapping(story.Key.Value)
	}
	index, err := c.storyOpportunities()
	if err != nil {
		return storyJSON{}, err
	}
	out.Opportunities = index[story.Key.Value]
	return out, nil
}

func toPersonaJSON(persona *domain.Persona) personaJSON {
	return personaJSON{Key: persona.Key, Name: persona.Name, Body: persona.Body}
}

// toPersonaRef resolves the key a story authored against the personas
// directory: a committed persona gets its display name and resource URI, an
// unknown one keeps the bare key. A story naming no persona gets no ref, which
// is what keeps "no actor stated" apart from "actor not committed yet".
func (c Config) toPersonaRef(key string) *personaRefJSON {
	if key == "" {
		return nil
	}
	ref := &personaRefJSON{Key: key}
	if persona, err := c.persona(key); err == nil {
		ref.Name = persona.Name
		ref.URI = uri.Persona(key)
	}
	return ref
}

// toPersonaRefs resolves the keys a story map declares, the way toTermRefs does
// for the terms it references.
func (c Config) toPersonaRefs(keys []string) []personaRefJSON {
	refs := make([]personaRefJSON, 0, len(keys))
	for _, key := range keys {
		if ref := c.toPersonaRef(key); ref != nil {
			refs = append(refs, *ref)
		}
	}
	return refs
}

func toTermJSON(term *domain.Term) termJSON {
	return termJSON{Key: term.Key, Ctx: term.Ctx, Name: term.Name, Body: term.Body}
}

// toTermRefs resolves the references a board authors against the ubiquitous
// directory: a committed term gets its display name and resource URI, an
// unknown one keeps whatever the reference said. A reference that does not even
// split — one nesting a context inside another, say — reports as a bare key, so
// a malformed reference reads like an uncommitted one rather than vanishing.
func (c Config) toTermRefs(rawRefs []string) []termRefJSON {
	refs := make([]termRefJSON, 0, len(rawRefs))
	for _, raw := range rawRefs {
		ref := termRefJSON{Key: raw}
		if ctx, key, ok := uri.SplitTermRef(raw); ok {
			ref.Key, ref.Ctx = key, ctx
			if term, err := c.term(raw); err == nil {
				ref.Name = term.Name
				ref.URI = uri.Term(ctx, key)
			}
		}
		refs = append(refs, ref)
	}
	return refs
}
