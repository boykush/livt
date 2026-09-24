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
// are read from — which is what lets the same layout name a path inside the
// repository for git to export and a path inside the export to read it back
// from. One revision is snapshotted by rebasing them onto the tree git wrote,
// so neither side knows which of the two it is looking at.
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

// Field is one line of an entry: a value, and what the thing holding it is
// called. Label is empty when the value is the item's own text — a rule is its
// text, not the value of a `name` key, and a diff that spelled the key would be
// a prettier YAML diff rather than a different reading of one.
type Field struct {
	Label string
	// Translate marks Label as a message catalog key. What livt names, livt
	// translates; what the livt repository named — a frontmatter key — is its
	// own word and is carried through as written.
	Translate bool
	Value     string
}

// Labels livt puts on its own fields. Where the site already has a word for
// something, that word is reused rather than a second one minted here.
const (
	LabelStatus       = "diff.field.status"
	LabelIssue        = "diff.field.issue"
	LabelSupersededBy = "diff.field.superseded-by"
	LabelRetired      = "diff.field.retired"
	LabelRelease      = "diff.field.release"
	LabelTerms        = "nav.ubiquitous"
	LabelActivity     = "label.activity"
	LabelStep         = "label.step"
	LabelStory        = "label.story"
	// canvasLabelPrefix addresses a canvas box's heading, which the sheet
	// already reads off the same key.
	canvasLabelPrefix = "canvas."
	// statusPrefix addresses a rule status as a word rather than as the value
	// written in the file.
	statusPrefix = "diff.status."
)

// Entry is one addressable point of the livt repository at one revision: the
// URI it answers to, what to call it on a page, and its own fields. Only its
// own — a rule's text is no part of its mapping's entry, so a reworded rule is
// one change rather than two.
type Entry struct {
	URI  string
	Kind uri.Kind
	// Parent is the mapping a rule, example, or question hangs off, so the diff
	// can gather them under it. Empty on everything addressed in its own right.
	Parent string
	Title  string
	Fields []Field
	// Live is whether the entry is spec in this revision. A rule closed and an
	// example or question retired are still on file, and still have to be read
	// back by their URI, but the spec has stopped asking for them.
	Live bool
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
		// Named as its board is, since the mapping heads a group of rules. Its own
		// name is a line only when it has one: every nameless mapping carrying a
		// blank line would add one to each mapping a diff adds.
		fields := listed(LabelTerms, em.Ubiquitous)
		if em.Name != "" {
			fields = append([]Field{text(em.Name)}, fields...)
		}
		s.add(Entry{
			URI:    mappingURI,
			Kind:   uri.KindMapping,
			Live:   true,
			Title:  em.DisplayName(parser.FindStoryByKey(dirs.Stories, em.StoryKey)),
			Fields: fields,
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
				Fields: ruleFields(r),
				Live:   r.Status.Active(),
			})
			for _, ex := range r.Examples {
				s.add(Entry{
					URI:    uri.Example(key, r.ID, ex.ID),
					Kind:   uri.KindExample,
					Parent: mappingURI,
					Title:  r.ID + " " + ex.ID,
					Fields: itemFields(ex.Name, ex.Retired, ex.SupersededBy),
					// An example under a closed rule went with it: the statement
					// it illustrates is no longer one the spec makes.
					Live: !ex.Retired && r.Status.Active(),
				})
			}
		}
		for _, q := range em.Questions {
			s.add(Entry{
				URI:    uri.Question(key, q.ID),
				Kind:   uri.KindQuestion,
				Parent: mappingURI,
				Title:  q.ID,
				Fields: itemFields(q.Text, q.Retired, q.SupersededBy),
				Live:   !q.Retired,
			})
		}
	}
	return nil
}

// ruleFields always spells the status, so a proposal being agreed reads as the
// one-line change it is rather than as a line appearing out of nowhere — and so
// a rule that writes its default explicitly diffs against one that omits it as
// no change at all. The status is carried as a message key, since the site has
// words for these and "accepted" is the spelling in the file, not the reading.
func ruleFields(r domain.Rule) []Field {
	fields := []Field{text(r.Name), labelled(LabelStatus, statusPrefix+string(r.Status.OrDefault()))}
	fields = append(fields, listed(LabelIssue, r.Issues)...)
	return append(fields, listed(LabelSupersededBy, r.SupersededBy)...)
}

// itemFields renders an example or a question, which are the same shape: their
// own text, and what became of them.
func itemFields(body string, retired bool, supersededBy []string) []Field {
	fields := []Field{text(body)}
	if retired {
		fields = append(fields, flag(LabelRetired))
	}
	return append(fields, listed(LabelSupersededBy, supersededBy)...)
}

func scanOpportunities(s *Snapshot, dirs Dirs) error {
	opportunities, err := parser.ParseAllOpportunities(dirs.Opportunities)
	if err != nil {
		return err
	}
	for _, o := range opportunities {
		key := o.Key.Value
		fields := append([]Field{text(o.Name)}, bodyFields(o.Body)...)
		s.add(Entry{
			URI:    uri.Opportunity(key),
			Kind:   uri.KindOpportunity,
			Live:   true,
			Title:  o.DisplayName(),
			Fields: append(fields, metaFields(o.Meta)...),
		})
		// A canvas is read through its opportunity's key, so it is scanned here
		// rather than from a directory walk of its own: the two are joined by
		// that key, and only one of them names it.
		canvas, err := parser.ParseOpportunityCanvas(filepath.Join(dirs.Canvases, key+".yaml"))
		if err != nil {
			continue
		}
		s.add(Entry{
			URI:    uri.OpportunityCanvas(key),
			Kind:   uri.KindOpportunityCanvas,
			Live:   true,
			Title:  o.DisplayName(),
			Fields: canvasFields(canvas),
		})
	}
	return nil
}

// canvasFields names each sticky by the box it sits in, using the heading the
// sheet itself prints over that box.
func canvasFields(c *domain.OpportunityCanvas) []Field {
	var fields []Field
	for _, box := range c.Boxes() {
		fields = append(fields, listed(canvasLabelPrefix+box.Key, box.Items)...)
	}
	return append(fields, listed(LabelTerms, c.Ubiquitous)...)
}

func scanStoryMaps(s *Snapshot, dirs Dirs) error {
	maps, err := parser.ParseAllStoryMaps(dirs.USM)
	if err != nil {
		return err
	}
	for _, m := range maps {
		s.add(Entry{
			URI:    uri.StoryMap(m.Name),
			Kind:   uri.KindStoryMap,
			Live:   true,
			Title:  m.Name,
			Fields: storyMapFields(m),
		})
	}
	return nil
}

// storyMapFields walks the backbone in the order the board is read, each card
// named by the kind of sticky it is. A card that moved between steps reads as
// the move it is, because the step above it moved with it.
func storyMapFields(m *domain.StoryMap) []Field {
	fields := []Field{text(m.Name)}
	for _, r := range m.Releases {
		fields = append(fields, labelled(LabelRelease, r.DisplayName(0)))
	}
	for _, a := range m.Activities {
		fields = append(fields, labelled(LabelActivity, a.Name))
		for _, step := range a.Steps {
			fields = append(fields, labelled(LabelStep, step.Name))
			for _, card := range step.Stories {
				fields = append(fields, labelled(LabelStory, card.Name+releaseSuffix(card)))
			}
		}
	}
	return append(fields, listed(LabelTerms, m.Ubiquitous)...)
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
		fields := append([]Field{text(story.Name)}, bodyFields(story.Body)...)
		s.add(Entry{
			URI:    uri.Story(story.Key.Value),
			Kind:   uri.KindStory,
			Live:   true,
			Title:  story.DisplayName(),
			Fields: append(fields, metaFields(story.Meta)...),
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
			URI:    uri.Term(t.Ctx, t.Key),
			Kind:   uri.KindTerm,
			Live:   true,
			Title:  t.Name,
			Fields: append([]Field{text(t.Name)}, bodyFields(t.Body)...),
		})
	}
	return nil
}

// text is the item speaking for itself, with nothing named in front of it.
func text(value string) Field {
	return Field{Value: value}
}

// labelled is a field livt names and so translates.
func labelled(label, value string) Field {
	return Field{Label: label, Translate: true, Value: value}
}

// flag is a field that says something by being there at all. It carries no
// value, so the line reads as the word itself: gaining it is the change.
func flag(label string) Field {
	return Field{Label: label, Translate: true}
}

// authored is a field the livt repository named, which is not livt's to
// translate — a frontmatter key is the repository's own word.
func authored(label, value string) Field {
	return Field{Label: label, Value: value}
}

// listed renders a list field one line per item, so adding an item to it is one
// added line rather than a rewrite of the whole list.
func listed(label string, values []string) []Field {
	fields := make([]Field, 0, len(values))
	for _, v := range values {
		fields = append(fields, labelled(label, v))
	}
	return fields
}

func metaFields(meta []domain.MetaField) []Field {
	fields := make([]Field, 0, len(meta))
	for _, m := range meta {
		fields = append(fields, authored(m.Key, m.Value))
	}
	return fields
}

// bodyFields keeps prose line by line, so a reworded sentence diffs to that
// sentence. Blank lines at either end are dropped: they are how the file is
// laid out, not anything the spec says.
func bodyFields(body string) []Field {
	trimmed := strings.Trim(body, "\n")
	if trimmed == "" {
		return nil
	}
	var fields []Field
	for _, line := range strings.Split(trimmed, "\n") {
		fields = append(fields, text(line))
	}
	return fields
}
