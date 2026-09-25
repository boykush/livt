package builder

import (
	"github.com/boykush/livt/internal/diff"
	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/uri"
)

// The board draws the active view: what the spec asks for today. During a diff
// it also draws what the spec stopped asking in that range, because that is the
// one reading a review is for — and a count of what went says nothing about
// what it said.
//
// Such a sticky goes back in the language the board already has for one that is
// not spec: a proposal is pale, dashed and stamped because it is not spec yet,
// and this is the same statement inverted.

// board is what a mapping page renders, and which of its stickies are only
// there because they were withdrawn in this diff. What became of one is on its
// own mark, so nothing here has to say it a second time.
type board struct {
	Mapping *domain.ExampleMapping
	Ghosts  map[string]bool
	// AutomationKnown says whether to draw the automation axis at all. False is
	// a board nothing has spoken about, where the ✓ legend would promise a mark
	// no sticky can carry.
	AutomationKnown bool
}

// boardFor builds the board for one mapping. With no diff it is the active view
// and nothing more, which is every build that was given no revisions.
func (b *Builder) boardFor(em *domain.ExampleMapping, automationKnown bool) board {
	if b.diffResult == nil {
		return board{Mapping: em.Active(), AutomationKnown: automationKnown}
	}
	out := &domain.ExampleMapping{StoryKey: em.StoryKey, Name: em.Name, Ubiquitous: em.Ubiquitous}
	ghosts := make(map[string]bool)
	key := em.StoryKey.Value

	for _, r := range em.Rules {
		live := r.Status.Active()
		// A rule closed before this range left the board before this range, and
		// putting it back would answer a question nobody asked of this diff.
		if !live && !b.changed(uri.Rule(key, r.ID)) {
			continue
		}
		if !live {
			ghosts[uri.Rule(key, r.ID)] = true
		}
		var examples []domain.Example
		for _, ex := range r.Examples {
			exURI := uri.Example(key, r.ID, ex.ID)
			if ex.Retired && !b.changed(exURI) {
				continue
			}
			// An example under a closed rule went with the rule: the statement
			// it illustrates is no longer one the spec makes.
			if ex.Retired || !live {
				ghosts[exURI] = true
			}
			examples = append(examples, ex)
		}
		r.Examples = examples
		out.Rules = append(out.Rules, r)
	}

	for _, q := range em.Questions {
		qURI := uri.Question(key, q.ID)
		if q.Retired && !b.changed(qURI) {
			continue
		}
		if q.Retired {
			ghosts[qURI] = true
		}
		out.Questions = append(out.Questions, q)
	}

	b.addDeleted(out, ghosts, em.StoryKey)
	return board{Mapping: out, Ghosts: ghosts, AutomationKnown: automationKnown}
}

func (b *Builder) changed(u string) bool {
	_, ok := b.diffByURI[u]
	return ok
}

// addDeleted puts back what the working tree does not hold at all. Its text
// comes off the diff, which carries the base revision's own fields — the only
// copy of it left anywhere the build can reach.
func (b *Builder) addDeleted(out *domain.ExampleMapping, ghosts map[string]bool, key domain.StoryKey) {
	mappingURI := uri.Mapping(key.Value)
	for _, c := range b.diffResult.Changes {
		if c.Status != diff.StatusRemoved || c.Parent != mappingURI {
			continue
		}
		parsed, ok := uri.Parse(c.URI)
		if !ok {
			continue
		}
		ghosts[c.URI] = true
		switch parsed.Kind {
		case uri.KindRule:
			out.Rules = append(out.Rules, domain.Rule{ID: parsed.RuleID, Name: deletedText(c), Status: domain.RuleRetired})
		case uri.KindExample:
			addDeletedExample(out, parsed, deletedText(c))
		case uri.KindQuestion:
			out.Questions = append(out.Questions, domain.Question{ID: parsed.QuestionID, Text: deletedText(c), Retired: true})
		}
	}
}

// addDeletedExample hangs a deleted example off its rule, adding the rule as a
// ghost of its own when it went too. An example under no rule is not a sticky
// the board can place.
func addDeletedExample(out *domain.ExampleMapping, parsed uri.Parsed, text string) {
	example := domain.Example{ID: parsed.ExampleID, Name: text, Retired: true}
	for i := range out.Rules {
		if out.Rules[i].ID == parsed.RuleID {
			out.Rules[i].Examples = append(out.Rules[i].Examples, example)
			return
		}
	}
	out.Rules = append(out.Rules, domain.Rule{ID: parsed.RuleID, Status: domain.RuleRetired, Examples: []domain.Example{example}})
}

// deletedText is the item's own text as the base revision had it: the first
// line the diff carries with no field name in front of it.
func deletedText(c diff.Change) string {
	for _, l := range c.Lines {
		if l.Field.Label == "" {
			return l.Field.Value
		}
	}
	return ""
}
