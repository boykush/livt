package mcp

import (
	"errors"
	"fmt"

	"github.com/boykush/livt/internal/uri"
)

// ErrNotFound marks a well-formed livt URI naming an item the livt repository does not
// hold. Callers tell it apart from a malformed URI, which never reaches here,
// and from a livt repository that fails to read, which surfaces as itself.
//
// A retired item is not one of these: it resolves like any other, carrying
// retired, so a reference filed against it keeps working
// (livt://mapping/trace-test-to-rule/rule/R-05).
var ErrNotFound = errors.New("not found")

type notFoundError struct{ err error }

func (e notFoundError) Error() string        { return e.err.Error() }
func (e notFoundError) Unwrap() error        { return e.err }
func (e notFoundError) Is(target error) bool { return target == ErrNotFound }
func notFound(err error) error               { return notFoundError{err} }

// Resolve returns the payload a resources/read of this URI serves. The CLI
// resolves through here too, so the two surfaces cannot drift into different
// JSON for the same URI.
func (s *Server) Resolve(p uri.Parsed) (any, error) {
	switch p.Kind {
	case uri.KindMapping:
		em, err := s.cfg.exampleMapping(p.StoryKey)
		if err != nil {
			return nil, notFound(err)
		}
		return exampleMappingResult{versioned: s.versioned(), Mapping: s.cfg.toExampleMappingJSON(em)}, nil
	case uri.KindRule:
		rule, err := s.cfg.rule(p.StoryKey, p.RuleID)
		if err != nil {
			return nil, notFound(err)
		}
		return ruleResult{versioned: s.versioned(), Rule: toRuleJSON(p.StoryKey, rule)}, nil
	case uri.KindExample:
		example, err := s.cfg.example(p.StoryKey, p.RuleID, p.ExampleID)
		if err != nil {
			return nil, notFound(err)
		}
		return exampleResult{versioned: s.versioned(), Example: toExampleJSON(p.StoryKey, p.RuleID, example)}, nil
	case uri.KindQuestion:
		question, err := s.cfg.question(p.StoryKey, p.QuestionID)
		if err != nil {
			return nil, notFound(err)
		}
		return questionResult{versioned: s.versioned(), Question: toQuestionJSON(p.StoryKey, question)}, nil
	case uri.KindOpportunity:
		o, err := s.cfg.opportunity(p.OpportunityKey)
		if err != nil {
			return nil, notFound(err)
		}
		// A malformed story map is a read failure, not a missing opportunity:
		// dropping the maps would misreport it as never taken on.
		out, err := s.cfg.toOpportunityJSON(o)
		if err != nil {
			return nil, err
		}
		return opportunityResult{versioned: s.versioned(), Opportunity: out}, nil
	case uri.KindOpportunityCanvas:
		canvas, err := s.cfg.opportunityCanvas(p.OpportunityKey)
		if err != nil {
			return nil, notFound(err)
		}
		return opportunityCanvasResult{versioned: s.versioned(), Canvas: s.cfg.toOpportunityCanvasJSON(canvas)}, nil
	case uri.KindStoryMap:
		sm, err := s.cfg.storyMap(p.MapName)
		if err != nil {
			return nil, notFound(err)
		}
		return storyMapResult{versioned: s.versioned(), StoryMap: s.cfg.toStoryMapJSON(sm)}, nil
	case uri.KindStory:
		story, err := s.cfg.story(p.StoryKey)
		if err != nil {
			return nil, notFound(err)
		}
		// A malformed story map is a read failure, not a missing story: dropping
		// the opportunities would misreport the story as unattached.
		out, err := s.cfg.toStoryJSON(story)
		if err != nil {
			return nil, err
		}
		return storyResult{versioned: s.versioned(), Story: out}, nil
	case uri.KindTerm:
		term, err := s.cfg.term(uri.TermRef(p.TermCtx, p.TermKey))
		if err != nil {
			return nil, notFound(err)
		}
		return termResult{versioned: s.versioned(), Term: toTermJSON(term)}, nil
	}
	return nil, fmt.Errorf("cannot resolve livt URI kind %q", p.Kind)
}

// Verify reports whether the livt repository holds the item a URI names, without
// building its payload. The URL form needs this: a deployed link derives from
// the URI alone, so nothing else would stop a mistyped reference resolving to a
// confident link to a page that was never built.
func (c Config) Verify(p uri.Parsed) error {
	var err error
	switch p.Kind {
	case uri.KindMapping:
		_, err = c.exampleMapping(p.StoryKey)
	case uri.KindRule:
		_, err = c.rule(p.StoryKey, p.RuleID)
	case uri.KindExample:
		_, err = c.example(p.StoryKey, p.RuleID, p.ExampleID)
	case uri.KindQuestion:
		_, err = c.question(p.StoryKey, p.QuestionID)
	case uri.KindOpportunity:
		_, err = c.opportunity(p.OpportunityKey)
	case uri.KindOpportunityCanvas:
		_, err = c.opportunityCanvas(p.OpportunityKey)
	case uri.KindStoryMap:
		_, err = c.storyMap(p.MapName)
	case uri.KindStory:
		_, err = c.story(p.StoryKey)
	case uri.KindTerm:
		_, err = c.term(uri.TermRef(p.TermCtx, p.TermKey))
	default:
		return fmt.Errorf("cannot verify livt URI kind %q", p.Kind)
	}
	if err != nil {
		return notFound(err)
	}
	return nil
}
