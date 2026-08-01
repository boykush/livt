// Package uri builds and parses livt URIs — the deployment-independent way to
// address one point of the spec: a story map, a story, an example mapping, a
// rule, an example, a question, or a ubiquitous language term. The MCP server,
// the CLI, and the site build all have to agree on the form, so it lives here
// rather than inside any one of them.
package uri

import (
	"net/url"
	"strings"
)

// URI templates in RFC 6570 form, as the MCP server advertises them.
const (
	MappingTemplate  = "livt://mapping/{story_key}"
	RuleTemplate     = "livt://mapping/{story_key}/rule/{rule_id}"
	ExampleTemplate  = "livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}"
	QuestionTemplate = "livt://mapping/{story_key}/question/{question_id}"
	StoryMapTemplate = "livt://story-map/{map_name}"
	StoryTemplate    = "livt://story/{story_key}"
	TermTemplate     = "livt://ubiquitous/{term_key}"
)

const (
	mappingPrefix  = "livt://mapping/"
	storyMapPrefix = "livt://story-map/"
	storyPrefix    = "livt://story/"
	termPrefix     = "livt://ubiquitous/"
	ruleInfix      = "/rule/"
	exampleInfix   = "/example/"
	questionInfix  = "/question/"
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

// Story builds the URI for a story's name, body, and meta.
func Story(storyKey string) string {
	return storyPrefix + storyKey
}

// Term builds the URI for a ubiquitous language term.
func Term(termKey string) string {
	return termPrefix + termKey
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

// ParseStory extracts the story key from a story URI.
func ParseStory(s string) (storyKey string, ok bool) {
	key, found := strings.CutPrefix(s, storyPrefix)
	if !found || !ValidSegment(key) {
		return "", false
	}
	return key, true
}

// ParseTerm extracts the term key from a ubiquitous term URI.
func ParseTerm(s string) (termKey string, ok bool) {
	key, found := strings.CutPrefix(s, termPrefix)
	if !found || !ValidSegment(key) {
		return "", false
	}
	return key, true
}

// ValidSegment guards externally-supplied URI/key segments so they resolve to a
// file inside the master's directories and cannot traverse out of them. It is
// also what keeps the URI shapes apart: a segment can hold no "/", so a longer
// shape never parses as a shorter one.
func ValidSegment(s string) bool {
	return s != "" && s != "." && s != ".." &&
		!strings.ContainsAny(s, `/\`) && !strings.Contains(s, "..")
}
