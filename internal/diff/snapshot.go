// Package diff says what changed between two revisions of a livt repository,
// one livt URI at a time. The unit is the URI rather than the file or the line
// because that is the address a reader already has for a rule, an example, or a
// term: a reviewer reading a YAML diff has to rebuild the spec from indentation
// and quoting before the change means anything.
package diff

import (
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
)

// Dirs are the livt repository's input directories, relative to the root they
// are read from. One revision is snapshotted by rebasing them onto the tree git
// exported for it, so the same layout serves the working tree and a revision
// without either side knowing which it is.
type Dirs struct {
	Opportunities string
	Canvases      string
	Mappings      string
	Stories       string
	USM           string
	Ubiquitous    string
}

func (d Dirs) under(root string) Dirs {
	return Dirs{
		Opportunities: filepath.Join(root, d.Opportunities),
		Canvases:      filepath.Join(root, d.Canvases),
		Mappings:      filepath.Join(root, d.Mappings),
		Stories:       filepath.Join(root, d.Stories),
		USM:           filepath.Join(root, d.USM),
		Ubiquitous:    filepath.Join(root, d.Ubiquitous),
	}
}

// paths lists the directories as the repository holds them, for git to export.
func (d Dirs) paths() []string {
	return []string{d.Opportunities, d.Canvases, d.Mappings, d.Stories, d.USM, d.Ubiquitous}
}

// Entry is one addressable point of the livt repository at one revision: the
// URI it answers to, what to call it on a page, and its own fields as lines.
// Only its own — a rule's text is no part of its mapping's entry, so a reworded
// rule is one change rather than two.
type Entry struct {
	URI  string
	Kind uri.Kind
	// Parent is the mapping a rule, example, or question hangs off, so the diff
	// can gather them under it. Empty on everything addressed in its own right.
	Parent string
	Title  string
	Lines  []string
}

// Snapshot is every entry of one revision, in the order the livt repository
// reads.
type Snapshot struct {
	Entries []Entry
	byURI   map[string]int
}

func (s *Snapshot) add(e Entry) {
	s.byURI[e.URI] = len(s.Entries)
	s.Entries = append(s.Entries, e)
}

// index is where the URI sits in reading order, and whether it is held at all.
func (s *Snapshot) index(uri string) (int, bool) {
	i, ok := s.byURI[uri]
	return i, ok
}

func (s *Snapshot) get(uri string) (Entry, bool) {
	i, ok := s.byURI[uri]
	if !ok {
		return Entry{}, false
	}
	return s.Entries[i], true
}

// Scan reads every input directory into one snapshot. A directory that does not
// exist reads as empty rather than failing: a revision from before a resource
// type existed is a legitimate side of a diff, and the whole point is to show
// what it did not have.
func Scan(dirs Dirs) (*Snapshot, error) {
	s := &Snapshot{byURI: make(map[string]int)}
	for _, scan := range []func(*Snapshot, Dirs) error{
		scanMappings,
		scanOpportunities,
		scanStoryMaps,
		scanStories,
		scanTerms,
	} {
		if err := scan(s, dirs); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func scanMappings(s *Snapshot, dirs Dirs) error {
	mappings, err := parser.ParseAllExampleMappings(dirs.Mappings)
	if err != nil {
		return err
	}
	for _, em := range mappings {
		key := em.StoryKey.Value
		mappingURI := uri.Mapping(key)
		// Named by its story, which is what a reader calls the board — the key
		// alone is the filename, and the mapping heads a group of rules.
		s.add(Entry{
			URI:   mappingURI,
			Kind:  uri.KindMapping,
			Title: parser.FindStoryByKey(dirs.Stories, em.StoryKey).DisplayName(),
			Lines: prefixed("ubiquitous", em.Ubiquitous),
		})
		// Scanned as recorded rather than as the board shows it: a retirement
		// and an agreement are exactly the changes a reviewer came for, and
		// Active() is what hides them.
		for _, r := range em.Rules {
			s.add(Entry{
				URI:    uri.Rule(key, r.ID),
				Kind:   uri.KindRule,
				Parent: mappingURI,
				Title:  r.ID,
				Lines:  ruleLines(r),
			})
			for _, ex := range r.Examples {
				s.add(Entry{
					URI:    uri.Example(key, r.ID, ex.ID),
					Kind:   uri.KindExample,
					Parent: mappingURI,
					Title:  r.ID + " " + ex.ID,
					Lines:  itemLines("name", ex.Name, ex.Retired, ex.SupersededBy),
				})
			}
		}
		for _, q := range em.Questions {
			s.add(Entry{
				URI:    uri.Question(key, q.ID),
				Kind:   uri.KindQuestion,
				Parent: mappingURI,
				Title:  q.ID,
				Lines:  itemLines("text", q.Text, q.Retired, q.SupersededBy),
			})
		}
	}
	return nil
}

// ruleLines always spells the status, so a proposal being agreed reads as the
// one-line change it is rather than as a line appearing out of nowhere — and so
// a rule that writes its default explicitly diffs against one that omits it as
// no change at all.
func ruleLines(r domain.Rule) []string {
	lines := []string{field("name", r.Name), field("status", string(r.Status.OrDefault()))}
	if r.Automated {
		lines = append(lines, field("automated", "true"))
	}
	lines = append(lines, prefixed("issue", r.Issues)...)
	return append(lines, prefixed("superseded_by", r.SupersededBy)...)
}

// itemLines renders an example or a question, which differ only in what their
// text field is called.
func itemLines(textKey, text string, retired bool, supersededBy []string) []string {
	lines := []string{field(textKey, text)}
	if retired {
		lines = append(lines, field("retired", "true"))
	}
	return append(lines, prefixed("superseded_by", supersededBy)...)
}

func scanOpportunities(s *Snapshot, dirs Dirs) error {
	opportunities, err := parser.ParseAllOpportunities(dirs.Opportunities)
	if err != nil {
		return err
	}
	for _, o := range opportunities {
		key := o.Key.Value
		lines := append([]string{field("name", o.Name)}, bodyLines(o.Body)...)
		s.add(Entry{
			URI:   uri.Opportunity(key),
			Kind:  uri.KindOpportunity,
			Title: o.DisplayName(),
			Lines: append(lines, metaLines(o.Meta)...),
		})
		// A canvas is read through its opportunity's key, so it is scanned here
		// rather than from a directory walk of its own: the two are joined by
		// that key, and only one of them names it.
		canvas, err := parser.ParseOpportunityCanvas(filepath.Join(dirs.Canvases, key+".yaml"))
		if err != nil {
			continue
		}
		s.add(Entry{
			URI:   uri.OpportunityCanvas(key),
			Kind:  uri.KindOpportunityCanvas,
			Title: o.DisplayName(),
			Lines: canvasLines(canvas),
		})
	}
	return nil
}

func canvasLines(c *domain.OpportunityCanvas) []string {
	var lines []string
	for _, box := range c.Boxes() {
		lines = append(lines, prefixed(box.Key, box.Items)...)
	}
	return append(lines, prefixed("ubiquitous", c.Ubiquitous)...)
}

func scanStoryMaps(s *Snapshot, dirs Dirs) error {
	maps, err := parser.ParseAllStoryMaps(dirs.USM)
	if err != nil {
		return err
	}
	for _, m := range maps {
		s.add(Entry{
			URI:   uri.StoryMap(m.Name),
			Kind:  uri.KindStoryMap,
			Title: m.Name,
			Lines: storyMapLines(m),
		})
	}
	return nil
}

// storyMapLines walks the backbone in the order the board is read. The nesting
// is spelled with indentation so a card moving between steps reads as the move
// it is, rather than as two unrelated lines.
func storyMapLines(m *domain.StoryMap) []string {
	lines := []string{field("name", m.Name)}
	for _, r := range m.Releases {
		lines = append(lines, field("release "+r.ID, r.Name))
	}
	for _, a := range m.Activities {
		lines = append(lines, field("activity "+a.Key, a.Name))
		for _, step := range a.Steps {
			lines = append(lines, "  "+field("step "+step.Key, step.Name))
			for _, card := range step.Stories {
				lines = append(lines, "    "+field("story "+card.Key.Value, card.Name)+releaseSuffix(card))
			}
		}
	}
	return append(lines, prefixed("ubiquitous", m.Ubiquitous)...)
}

func releaseSuffix(card domain.StoryCard) string {
	if card.Release == "" {
		return ""
	}
	return " (" + card.Release + ")"
}

func scanStories(s *Snapshot, dirs Dirs) error {
	stories, err := parser.ParseAllStories(dirs.Stories)
	if err != nil {
		return err
	}
	for _, story := range stories {
		lines := append([]string{field("name", story.Name)}, bodyLines(story.Body)...)
		s.add(Entry{
			URI:   uri.Story(story.Key.Value),
			Kind:  uri.KindStory,
			Title: story.DisplayName(),
			Lines: append(lines, metaLines(story.Meta)...),
		})
	}
	return nil
}

func scanTerms(s *Snapshot, dirs Dirs) error {
	terms, err := parser.ParseAllTerms(dirs.Ubiquitous)
	if err != nil {
		return err
	}
	for _, t := range terms {
		s.add(Entry{
			URI:   uri.Term(t.Ctx, t.Key),
			Kind:  uri.KindTerm,
			Title: t.Name,
			Lines: append([]string{field("name", t.Name)}, bodyLines(t.Body)...),
		})
	}
	return nil
}

func field(key, value string) string {
	return key + ": " + value
}

// prefixed renders a list field one line per item, so adding an item to it is
// one added line rather than a rewrite of the whole list.
func prefixed(key string, values []string) []string {
	lines := make([]string, 0, len(values))
	for _, v := range values {
		lines = append(lines, field(key, v))
	}
	return lines
}

func metaLines(meta []domain.MetaField) []string {
	lines := make([]string, 0, len(meta))
	for _, m := range meta {
		lines = append(lines, field("meta "+m.Key, m.Value))
	}
	return lines
}

// bodyLines keeps prose line by line, so a reworded sentence diffs to that
// sentence. Blank lines at either end are dropped: they are how the file is
// laid out, not anything the spec says.
func bodyLines(body string) []string {
	trimmed := strings.Trim(body, "\n")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "\n")
}
