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

// buildMappings builds example mapping HTML pages and returns a preview tile per
// mapping for the Example Mappings overview page, plus everything the mappings
// leave unfinished — open questions and un-automated rules — for the Tasks page.
func (b *Builder) buildMappings() ([]mappingTile, taskSet, error) {
	files, err := filepath.Glob(filepath.Join(b.MappingsDir, "*.yaml"))
	if err != nil {
		return nil, taskSet{}, err
	}

	var tiles []mappingTile
	var open taskSet
	for _, f := range files {
		em, err := parser.ParseExampleMapping(f)
		if err != nil {
			return nil, taskSet{}, fmt.Errorf("parse %s: %w", f, err)
		}

		storyName := b.resolveStoryName(em.StoryKey)
		storyPath := ""
		if b.hasStoryPage(em.StoryKey) {
			storyPath = "../" + uri.StoryPage(em.StoryKey.Value)
		}

		ubiquitous := b.resolveTermCards(em.Ubiquitous)

		outPath := filepath.Join(b.OutDir, "mapping", em.StoryKey.Value+".html")
		if err := b.buildMapping(outPath, em, storyName, storyPath, ubiquitous); err != nil {
			return nil, taskSet{}, err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))

		tiles = append(tiles, mappingTile{Key: em.StoryKey.Value, StoryName: storyName})
		// The Tasks page renders at the output root, so its links resolve from
		// there, not from the mapping/ directory.
		open.add(collectTasks(em, storyName, strings.TrimPrefix(storyPath, "../")))
	}

	return tiles, open, nil
}

// taskSet holds what the mappings leave unfinished, split by how it gets closed:
// a question by a conversation, an un-automated rule by a test.
type taskSet struct {
	Questions        []taskItem
	UnautomatedRules []taskItem
}

func (o *taskSet) add(other taskSet) {
	o.Questions = append(o.Questions, other.Questions...)
	o.UnautomatedRules = append(o.UnautomatedRules, other.UnautomatedRules...)
}

// collectTasks lifts one mapping's open questions and un-automated rules onto
// the Tasks page, each deep-linked to its own sticky on the board through the
// same derivation the board itself anchors by.
func collectTasks(em *domain.ExampleMapping, storyName, storyPath string) taskSet {
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
			StoryName:   storyName,
			StoryPath:   storyPath,
			MappingPath: mappingPath,
		}
	}

	var out taskSet
	for _, q := range active.Questions {
		out.Questions = append(out.Questions, item("question", q.ID, q.Text, sticky(q.ID, uri.QuestionPage)))
	}
	for _, r := range active.Rules {
		if r.Automated {
			continue
		}
		out.UnautomatedRules = append(out.UnautomatedRules, item("rule", r.ID, r.Name, sticky(r.ID, uri.RulePage)))
	}
	return out
}

func (b *Builder) resolveStoryName(key domain.StoryKey) string {
	return parser.FindStoryByKey(b.StoriesDir, key).DisplayName()
}

func (b *Builder) buildMapping(path string, em *domain.ExampleMapping, storyName, storyPath string, ubiquitous []termCard) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderMapping(f, b.Lang, em, storyName, storyPath, ubiquitous)
}
