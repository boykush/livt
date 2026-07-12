package mcp

import "github.com/boykush/livt/internal/domain"

// versioned is embedded in every tool and resource payload so each result
// records the master version it reflects. Its json field is promoted to the top
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

// termResult is the body of the livt://ubiquitous/{term_key} resource.
type termResult struct {
	versioned
	Term termJSON `json:"term"`
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
	// Issues are the rule's automation Issue URLs as recorded on the master.
	Issues []string `json:"issues,omitempty"`
	// Automated is always present: consumers read the recorded judgment
	// without distinguishing absent from false.
	Automated bool `json:"automated"`
}

type questionJSON struct {
	ID   string `json:"id"`
	Text string `json:"text"`
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
	Body string `json:"body"`
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

type termJSON struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Body string `json:"body"`
}

// termRefJSON points at a ubiquitous language term referenced by key from a
// mapping or story map. Name and URI are filled when the term file exists; a
// bare key means the term is not committed yet (mirroring how the site renders
// unresolved term cards).
type termRefJSON struct {
	Key  string `json:"key"`
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
		examples = append(examples, exampleJSON{ID: e.ID, Name: e.Name})
	}
	return ruleJSON{ID: r.ID, URI: ruleURI(storyKey, r.ID), Name: r.Name, Examples: examples, Issues: r.Issues, Automated: r.Automated}
}

// toExampleMappingJSON is a Config method because resolving the referenced
// ubiquitous terms reads the master's ubiquitous directory.
func (c Config) toExampleMappingJSON(em *domain.ExampleMapping) exampleMappingJSON {
	rules := make([]ruleJSON, 0, len(em.Rules))
	for _, r := range em.Rules {
		rules = append(rules, toRuleJSON(em.StoryKey.Value, r))
	}
	questions := make([]questionJSON, 0, len(em.Questions))
	for _, q := range em.Questions {
		questions = append(questions, questionJSON{ID: q.ID, Text: q.Text})
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
					card.URI = storyURI(sc.Key.Value)
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
		Ubiquitous:      sm.Ubiquitous,
		UbiquitousTerms: c.toTermRefs(sm.Ubiquitous),
	}
}

// toStoryJSON is a Config method because linking the example mapping and
// resolving the story's opportunities read the master's mappings and usm
// directories. Unlike term refs, a malformed story map is an error, not a bare
// ref: silently dropping opportunities would misreport the story as unattached.
func (c Config) toStoryJSON(story *domain.Story) (storyJSON, error) {
	meta := make([]metaFieldJSON, 0, len(story.Meta))
	for _, m := range story.Meta {
		meta = append(meta, metaFieldJSON{Key: m.Key, Value: m.Value})
	}
	out := storyJSON{Key: story.Key.Value, Name: story.Name, Body: story.Body, Meta: meta}
	if c.hasExampleMapping(story.Key.Value) {
		out.ExampleMappingURI = mappingURI(story.Key.Value)
	}
	index, err := c.storyOpportunities()
	if err != nil {
		return storyJSON{}, err
	}
	out.Opportunities = index[story.Key.Value]
	return out, nil
}

func toTermJSON(term *domain.Term) termJSON {
	return termJSON{Key: term.Key, Name: term.Name, Body: term.Body}
}

// toTermRefs resolves referenced term keys against the ubiquitous directory: a
// committed term gets its display name and resource URI, an unknown or
// malformed one stays a bare key.
func (c Config) toTermRefs(keys []string) []termRefJSON {
	refs := make([]termRefJSON, 0, len(keys))
	for _, key := range keys {
		ref := termRefJSON{Key: key}
		if term, err := c.term(key); err == nil {
			ref.Name = term.Name
			ref.URI = termURI(key)
		}
		refs = append(refs, ref)
	}
	return refs
}
