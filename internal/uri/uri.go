// Package uri builds and parses livt URIs — the deployment-independent way to
// address one point of the spec: an opportunity, its canvas, a story map, a
// story, an example mapping, a rule, an example, a question, or a ubiquitous
// language term. The MCP server,
// the CLI, and the site build all have to agree on the form, so it lives here
// rather than inside any one of them.
package uri

import (
	"net/url"
	"strings"
)

// URI templates in RFC 6570 form, as the MCP server advertises them.
const (
	MappingTemplate     = "livt://mapping/{story_key}"
	RuleTemplate        = "livt://mapping/{story_key}/rule/{rule_id}"
	ExampleTemplate     = "livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}"
	QuestionTemplate    = "livt://mapping/{story_key}/question/{question_id}"
	StoryMapTemplate    = "livt://story-map/{map_name}"
	OpportunityTemplate = "livt://opportunity/{opportunity_key}"
	// OpportunityCanvasTemplate addresses the canvas filled in for an
	// opportunity. It sits beside the opportunity rather than under it, the way
	// a mapping sits beside its story: the two are joined by key, and either can
	// exist without the other.
	OpportunityCanvasTemplate = "livt://opportunity-canvas/{opportunity_key}"
	StoryTemplate             = "livt://story/{story_key}"
	TermTemplate              = "livt://ubiquitous/{term_key}"
	// ScopedTermTemplate addresses a term that belongs to one context. The two
	// term shapes both stand: a context is optional, so a term that holds across
	// contexts keeps the address it had before contexts existed.
	ScopedTermTemplate = "livt://ubiquitous/{ctx}/{term_key}"
)

const (
	mappingPrefix           = "livt://mapping/"
	storyMapPrefix          = "livt://story-map/"
	opportunityPrefix       = "livt://opportunity/"
	opportunityCanvasPrefix = "livt://opportunity-canvas/"
	storyPrefix             = "livt://story/"
	termPrefix              = "livt://ubiquitous/"
	ruleInfix               = "/rule/"
	exampleInfix            = "/example/"
	questionInfix           = "/question/"
	// termCtxSep joins a term's context to its key. It is the same separator the
	// filesystem uses, so ubiquitous/{ctx}/{key}.md, livt://ubiquitous/{ctx}/{key}
	// and the reference a board authors all read alike.
	termCtxSep = "/"
)

// Mapping builds the URI for a story's example mapping.
func Mapping(storyKey string) string {
	return mappingPrefix + storyKey
}

// Rule builds the URI for a single rule within a story's mapping.
func Rule(storyKey, ruleID string) string {
	return mappingPrefix + storyKey + ruleInfix + ruleID
}

// Example builds the URI for a single example. Example ids are numbered within
// their rule, so the rule is part of the address: EX-01 alone does not identify
// an example even inside one mapping.
func Example(storyKey, ruleID, exampleID string) string {
	return mappingPrefix + storyKey + ruleInfix + ruleID + exampleInfix + exampleID
}

// Question builds the URI for a single open question. Questions hang off the
// mapping, not off a rule, so the address stops at the story key.
func Question(storyKey, questionID string) string {
	return mappingPrefix + storyKey + questionInfix + questionID
}

// StoryMap builds the URI for a story map. Maps are addressed by display name —
// the same identifier the build output uses for story-map/{name}.html —
// percent-encoded because names are free text (often non-ASCII) and the URI
// needs a single valid path segment.
func StoryMap(name string) string {
	return storyMapPrefix + url.PathEscape(name)
}

// Opportunity builds the URI for an opportunity's name, statement, and meta.
// Unlike a story map, an opportunity is addressed by key rather than display
// name: the key is what the filename and the canvas join on, so it is already
// the identifier the livt repository carries.
func Opportunity(opportunityKey string) string {
	return opportunityPrefix + opportunityKey
}

// OpportunityCanvas builds the URI for an opportunity's canvas.
func OpportunityCanvas(opportunityKey string) string {
	return opportunityCanvasPrefix + opportunityKey
}

// Story builds the URI for a story's name, body, and meta.
func Story(storyKey string) string {
	return storyPrefix + storyKey
}

// Term builds the URI for a ubiquitous language term. An empty ctx addresses a
// term that holds across contexts — the whole point of making the context
// optional is that such a term is not forced into one.
func Term(ctx, termKey string) string {
	return termPrefix + TermRef(ctx, termKey)
}

// TermRef is how a term is named where a single string is wanted: the tail of
// its URI, the reference a board's ubiquitous list authors, and its anchor in
// the glossary table.
func TermRef(ctx, termKey string) string {
	if ctx == "" {
		return termKey
	}
	return ctx + termCtxSep + termKey
}

// SplitTermRef takes a term reference apart. Both segments are guarded, so a
// reference can neither escape the ubiquitous directory nor nest a context
// inside another: contexts are one level deep, the same depth the URI carries.
func SplitTermRef(ref string) (ctx, termKey string, ok bool) {
	if c, key, found := strings.Cut(ref, termCtxSep); found {
		if !ValidSegment(c) || !ValidSegment(key) {
			return "", "", false
		}
		return c, key, true
	}
	if !ValidSegment(ref) {
		return "", "", false
	}
	return "", ref, true
}

// ParseMapping extracts the story key from a mapping URI. It does not match
// rule or question URIs (their extra segments fail ValidSegment).
func ParseMapping(s string) (storyKey string, ok bool) {
	key, found := strings.CutPrefix(s, mappingPrefix)
	if !found || !ValidSegment(key) {
		return "", false
	}
	return key, true
}

// ParseRule extracts the story key and rule id from a rule URI. An example URI
// carries its own "/example/{id}" tail, which lands in the rule id and fails
// ValidSegment — so the two shapes never resolve to each other.
func ParseRule(s string) (storyKey, ruleID string, ok bool) {
	rest, found := strings.CutPrefix(s, mappingPrefix)
	if !found {
		return "", "", false
	}
	key, id, found := strings.Cut(rest, ruleInfix)
	if !found || !ValidSegment(key) || !ValidSegment(id) {
		return "", "", false
	}
	return key, id, true
}

// ParseExample extracts the story key, rule id, and example id from an example
// URI. A rule URI has no "/example/" segment, so it never parses as one.
func ParseExample(s string) (storyKey, ruleID, exampleID string, ok bool) {
	rest, found := strings.CutPrefix(s, mappingPrefix)
	if !found {
		return "", "", "", false
	}
	key, rest, found := strings.Cut(rest, ruleInfix)
	if !found {
		return "", "", "", false
	}
	id, exID, found := strings.Cut(rest, exampleInfix)
	if !found || !ValidSegment(key) || !ValidSegment(id) || !ValidSegment(exID) {
		return "", "", "", false
	}
	return key, id, exID, true
}

// ParseQuestion extracts the story key and question id from a question URI.
func ParseQuestion(s string) (storyKey, questionID string, ok bool) {
	rest, found := strings.CutPrefix(s, mappingPrefix)
	if !found {
		return "", "", false
	}
	key, id, found := strings.Cut(rest, questionInfix)
	if !found || !ValidSegment(key) || !ValidSegment(id) {
		return "", "", false
	}
	return key, id, true
}

// ParseStoryMap extracts the map name from a story map URI, undoing the
// percent-encoding applied by StoryMap. The name is only ever compared against
// parsed maps — never used as a file path — so it needs no segment guard.
func ParseStoryMap(s string) (name string, ok bool) {
	escaped, found := strings.CutPrefix(s, storyMapPrefix)
	if !found || escaped == "" {
		return "", false
	}
	name, err := url.PathUnescape(escaped)
	if err != nil || name == "" {
		return "", false
	}
	return name, true
}

// ParseOpportunity extracts the opportunity key from an opportunity URI. The
// canvas prefix is a different string, not a longer one, so a canvas URI can
// never be cut down to an opportunity URI here.
func ParseOpportunity(s string) (opportunityKey string, ok bool) {
	key, found := strings.CutPrefix(s, opportunityPrefix)
	if !found || !ValidSegment(key) {
		return "", false
	}
	return key, true
}

// ParseOpportunityCanvas extracts the opportunity key from a canvas URI.
func ParseOpportunityCanvas(s string) (opportunityKey string, ok bool) {
	key, found := strings.CutPrefix(s, opportunityCanvasPrefix)
	if !found || !ValidSegment(key) {
		return "", false
	}
	return key, true
}

// ParseStory extracts the story key from a story URI.
func ParseStory(s string) (storyKey string, ok bool) {
	key, found := strings.CutPrefix(s, storyPrefix)
	if !found || !ValidSegment(key) {
		return "", false
	}
	return key, true
}

// ParseTerm extracts the context and term key from a ubiquitous term URI. A
// context-free term yields an empty ctx. Both term shapes live under one prefix
// no other shape uses, so accepting two segments here cannot make a term URI
// parse as anything else.
func ParseTerm(s string) (ctx, termKey string, ok bool) {
	ref, found := strings.CutPrefix(s, termPrefix)
	if !found {
		return "", "", false
	}
	return SplitTermRef(ref)
}

// ValidSegment guards externally-supplied URI/key segments so they resolve to a
// file inside the livt repository's directories and cannot traverse out of them. It is
// also what keeps the URI shapes apart: a segment can hold no "/", so a longer
// shape never parses as a shorter one.
func ValidSegment(s string) bool {
	return s != "" && s != "." && s != ".." &&
		!strings.ContainsAny(s, `/\`) && !strings.Contains(s, "..")
}
