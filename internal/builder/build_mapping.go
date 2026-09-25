package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
)

// mappingTally is one mapping's live counts, keyed by story so an opportunity
// can sum the stories it took on. Taken off the active view, like the Tasks
// page and the sidebar badge: a retired rule is not spec anyone is still
// waiting on, so counting it would make an opportunity read as less finished
// than it is.
type mappingTally struct {
	// Name is what the board is called, so a row naming a story with no story
	// file reads as the board it leads to does rather than as its bare key.
	Name      string
	Rules     int
	Automated int
	Proposed  int
	// Unautomated counts the rules waiting for a test the way the Tasks page
	// lists them: neither automated nor proposed. It is not Rules less the
	// other two, because a proposal can carry a test ahead of its agreement
	// and then sits in both counts — subtracting both takes it out twice.
	Unautomated int
	Questions   int
	// QuestionsAsked counts the retired ones too — the whole an open question
	// is a part of. Without it the count of what is open has no size: five open
	// questions reads differently on a board that has settled ten than on one
	// that has settled none.
	QuestionsAsked int
}

// buildMappings builds example mapping HTML pages and returns a preview tile per
// mapping for the Example Mappings overview page, everything the mappings leave
// unfinished — open questions, proposed rules, and un-automated rules — for the
// Tasks page, and a tally per story for the opportunity pages.
func (b *Builder) buildMappings() ([]mappingTile, taskSet, map[string]mappingTally, error) {
	files, err := filepath.Glob(filepath.Join(b.MappingsDir, "*.yaml"))
	if err != nil {
		return nil, taskSet{}, nil, err
	}

	automations, err := b.automationIndex()
	if err != nil {
		return nil, taskSet{}, nil, err
	}
	b.automationKnown = !automations.Empty()

	var tiles []mappingTile
	var open taskSet
	tallies := make(map[string]mappingTally, len(files))
	for _, f := range files {
		em, err := parser.ParseExampleMapping(f)
		if err != nil {
			return nil, taskSet{}, nil, fmt.Errorf("parse %s: %w", f, err)
		}
		automations.Attach(em)

		name := em.DisplayName(parser.FindStoryByKey(b.StoriesDir, em.StoryKey))
		storyPath := ""
		if b.hasStoryPage(em.StoryKey) {
			storyPath = "../" + uri.StoryPage(em.StoryKey.Value)
		}

		ubiquitous := b.resolveTermCards(em.Ubiquitous)

		outPath := filepath.Join(b.OutDir, "mapping", em.StoryKey.Value+".html")
		if err := b.buildMapping(outPath, em, name, storyPath, ubiquitous, b.automationKnown); err != nil {
			return nil, taskSet{}, nil, err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))

		tiles = append(tiles, mappingTile{Key: em.StoryKey.Value, Name: name})
		// The Tasks page renders at the output root, so its links resolve from
		// there, not from the mapping/ directory.
		open.add(collectTasks(em, name, strings.TrimPrefix(storyPath, "../"), b.automationKnown))
		t := tally(em, b.automationKnown)
		t.Name = name
		tallies[em.StoryKey.Value] = t
	}

	return tiles, open, tallies, nil
}

// tally counts one mapping's live rules, questions, and the two standings an
// opportunity reports on. With nothing known about automation the two stay at
// zero, and the surfaces that read them are not drawn.
func tally(em *domain.ExampleMapping, known bool) mappingTally {
	active := em.Active()
	t := mappingTally{
		Rules:          len(active.Rules),
		Questions:      len(active.Questions),
		QuestionsAsked: len(em.Questions),
	}
	for _, r := range active.Rules {
		if r.Proposed() {
			t.Proposed++
		}
		if !known {
			continue
		}
		if r.Automated() {
			t.Automated++
		}
		if !r.Proposed() && !r.Automated() {
			t.Unautomated++
		}
	}
	return t
}

// taskSet holds what the mappings leave unfinished, split by how it gets closed:
// a question by a conversation, a proposed rule by agreement, an un-automated
// rule by a test.
type taskSet struct {
	Questions        []taskItem
	ProposedRules    []taskItem
	UnautomatedRules []taskItem
}

func (o *taskSet) add(other taskSet) {
	o.Questions = append(o.Questions, other.Questions...)
	o.ProposedRules = append(o.ProposedRules, other.ProposedRules...)
	o.UnautomatedRules = append(o.UnautomatedRules, other.UnautomatedRules...)
}

// collectTasks lifts one mapping's open questions, proposed rules, and
// un-automated rules onto the Tasks page, each deep-linked to its own sticky on
// the board through the same derivation the board itself anchors by.
func collectTasks(em *domain.ExampleMapping, mappingName, storyPath string, known bool) taskSet {
	// Retired items are not unfinished work: a retired question is not open, and
	// a retired rule is not waiting for a test. Left in, they would sit on this
	// page forever, since nothing can happen to close them.
	active := em.Active()

	key := active.StoryKey.Value
	// An item without an ID has no anchor to aim at (the board only renders one
	// for keyed stickies), so it links to the board itself.
	sticky := func(id string, locate func(string, string) string) string {
		if id == "" {
			return uri.MappingPage(key)
		}
		return locate(key, id)
	}
	item := func(kind, id, text, mappingPath string) taskItem {
		return taskItem{
			Kind:        kind,
			ID:          id,
			Text:        text,
			StoryKey:    key,
			MappingName: mappingName,
			StoryPath:   storyPath,
			MappingPath: mappingPath,
		}
	}

	var out taskSet
	for _, q := range active.Questions {
		out.Questions = append(out.Questions, item("question", q.ID, q.Text, sticky(q.ID, uri.QuestionPage)))
	}
	for _, r := range active.Rules {
		// Agreement is what closes a proposed rule, so it waits for one even when
		// a test already covers it.
		switch {
		case r.Proposed():
			out.ProposedRules = append(out.ProposedRules, item("proposed-rule", r.ID, r.Name, sticky(r.ID, uri.RulePage)))
		case known && !r.Automated():
			out.UnautomatedRules = append(out.UnautomatedRules, item("rule", r.ID, r.Name, sticky(r.ID, uri.RulePage)))
		}
	}
	return out
}

func (b *Builder) resolveStoryName(key domain.StoryKey) string {
	return parser.FindStoryByKey(b.StoriesDir, key).DisplayName()
}

func (b *Builder) buildMapping(path string, em *domain.ExampleMapping, name, storyPath string, ubiquitous []termCard, known bool) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	// The board is one level down, so the diff it points back to is too.
	return renderMapping(f, b.Lang, b.boardFor(em, known), name, storyPath, ubiquitous,
		b.diffMark("../", uri.Mapping(em.StoryKey.Value)), b.diffMarks("../"))
}
