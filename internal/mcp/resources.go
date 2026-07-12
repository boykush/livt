package mcp

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	mappingURIPrefix    = "livt://mapping/"
	ruleURIInfix        = "/rule/"
	storyMapURIPrefix   = "livt://story-map/"
	storyURIPrefix      = "livt://story/"
	termURIPrefix       = "livt://ubiquitous/"
	mappingURITemplate  = "livt://mapping/{story_key}"
	ruleURITemplate     = "livt://mapping/{story_key}/rule/{rule_id}"
	storyMapURITemplate = "livt://story-map/{map_name}"
	storyURITemplate    = "livt://story/{story_key}"
	termURITemplate     = "livt://ubiquitous/{term_key}"
)

// mappingURI builds the resource URI for a story's example mapping.
func mappingURI(storyKey string) string {
	return mappingURIPrefix + storyKey
}

// ruleURI builds the resource URI for a single rule within a story's mapping.
func ruleURI(storyKey, ruleID string) string {
	return mappingURIPrefix + storyKey + ruleURIInfix + ruleID
}

// storyMapURI builds the resource URI for a story map. Maps are addressed by
// display name — the same identifier the build output uses for
// story-map/{name}.html — percent-encoded because names are free text (often
// non-ASCII) and the URI needs a single valid path segment.
func storyMapURI(name string) string {
	return storyMapURIPrefix + url.PathEscape(name)
}

// storyURI builds the resource URI for a story's name, body, and meta.
func storyURI(storyKey string) string {
	return storyURIPrefix + storyKey
}

// termURI builds the resource URI for a ubiquitous language term.
func termURI(termKey string) string {
	return termURIPrefix + termKey
}

// parseMappingURI extracts the story key from a mapping resource URI. It does
// not match rule URIs (their extra "/rule/..." segment fails validSegment).
func parseMappingURI(uri string) (storyKey string, ok bool) {
	key, found := strings.CutPrefix(uri, mappingURIPrefix)
	if !found || !validSegment(key) {
		return "", false
	}
	return key, true
}

// parseRuleURI extracts the story key and rule id from a rule resource URI.
func parseRuleURI(uri string) (storyKey, ruleID string, ok bool) {
	rest, found := strings.CutPrefix(uri, mappingURIPrefix)
	if !found {
		return "", "", false
	}
	key, id, found := strings.Cut(rest, ruleURIInfix)
	if !found || !validSegment(key) || !validSegment(id) {
		return "", "", false
	}
	return key, id, true
}

// parseStoryMapURI extracts the map name from a story map resource URI,
// undoing the percent-encoding applied by storyMapURI. The name is only ever
// compared against parsed maps — never used as a file path — so it needs no
// segment guard.
func parseStoryMapURI(uri string) (name string, ok bool) {
	escaped, found := strings.CutPrefix(uri, storyMapURIPrefix)
	if !found || escaped == "" {
		return "", false
	}
	name, err := url.PathUnescape(escaped)
	if err != nil || name == "" {
		return "", false
	}
	return name, true
}

// parseStoryURI extracts the story key from a story resource URI.
func parseStoryURI(uri string) (storyKey string, ok bool) {
	key, found := strings.CutPrefix(uri, storyURIPrefix)
	if !found || !validSegment(key) {
		return "", false
	}
	return key, true
}

// parseTermURI extracts the term key from a ubiquitous term resource URI.
func parseTermURI(uri string) (termKey string, ok bool) {
	key, found := strings.CutPrefix(uri, termURIPrefix)
	if !found || !validSegment(key) {
		return "", false
	}
	return key, true
}

// validSegment guards externally-supplied URI/key segments so they resolve to a
// file inside the mappings directory and cannot traverse out of it.
func validSegment(s string) bool {
	return s != "" && s != "." && s != ".." &&
		!strings.ContainsAny(s, `/\`) && !strings.Contains(s, "..")
}

// registerResources exposes the master as addressable resources, so a client
// reads the spec by URI (story map -> story -> mapping -> rule, with ubiquitous
// terms linked from mappings and story maps) rather than calling a tool. Only
// resource templates are advertised — no concrete resources and no subscribe —
// keeping the server stateless with every read served fresh from disk.
func (s *Server) registerResources(srv *mcpsdk.Server) {
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "example-mapping",
		Title:       "Example mapping",
		Description: "An example mapping (rules, examples, questions, ubiquitous terms) for a story, addressed by story key.",
		MIMEType:    "application/json",
		URITemplate: mappingURITemplate,
	}, s.readMapping)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "rule",
		Title:       "Rule",
		Description: "A single rule and its examples from a story's example mapping.",
		MIMEType:    "application/json",
		URITemplate: ruleURITemplate,
	}, s.readRule)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "story-map",
		Title:       "Story map",
		Description: "A user story map (activities, steps, story cards, releases, ubiquitous terms), addressed by its display name (percent-encoded).",
		MIMEType:    "application/json",
		URITemplate: storyMapURITemplate,
	}, s.readStoryMap)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "story",
		Title:       "Story",
		Description: "A story's name, body, and frontmatter meta, addressed by story key.",
		MIMEType:    "application/json",
		URITemplate: storyURITemplate,
	}, s.readStory)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "ubiquitous-term",
		Title:       "Ubiquitous language term",
		Description: "A ubiquitous language term's name and definition, addressed by term key.",
		MIMEType:    "application/json",
		URITemplate: termURITemplate,
	}, s.readTerm)
}

// readMapping serves livt://mapping/{story_key}.
func (s *Server) readMapping(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	uri := req.Params.URI
	key, ok := parseMappingURI(uri)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	em, err := s.cfg.exampleMapping(key)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	return jsonResource(uri, exampleMappingResult{versioned: s.versioned(), Mapping: s.cfg.toExampleMappingJSON(em)})
}

// readRule serves livt://mapping/{story_key}/rule/{rule_id}.
func (s *Server) readRule(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	uri := req.Params.URI
	key, ruleID, ok := parseRuleURI(uri)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	rule, err := s.cfg.rule(key, ruleID)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	return jsonResource(uri, ruleResult{versioned: s.versioned(), Rule: toRuleJSON(key, rule)})
}

// readStoryMap serves livt://story-map/{map_name}.
func (s *Server) readStoryMap(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	uri := req.Params.URI
	name, ok := parseStoryMapURI(uri)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	sm, err := s.cfg.storyMap(name)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	return jsonResource(uri, storyMapResult{versioned: s.versioned(), StoryMap: s.cfg.toStoryMapJSON(sm)})
}

// readStory serves livt://story/{story_key}.
func (s *Server) readStory(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	uri := req.Params.URI
	key, ok := parseStoryURI(uri)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	story, err := s.cfg.story(key)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	out, err := s.cfg.toStoryJSON(story)
	if err != nil {
		return nil, err
	}
	return jsonResource(uri, storyResult{versioned: s.versioned(), Story: out})
}

// readTerm serves livt://ubiquitous/{term_key}.
func (s *Server) readTerm(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	uri := req.Params.URI
	key, ok := parseTermURI(uri)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	term, err := s.cfg.term(key)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(uri)
	}
	return jsonResource(uri, termResult{versioned: s.versioned(), Term: toTermJSON(term)})
}

// jsonResource marshals v as the single application/json content of a resource.
func jsonResource(uri string, v any) (*mcpsdk.ReadResourceResult, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{{
			URI:      uri,
			MIMEType: "application/json",
			Text:     string(body),
		}},
	}, nil
}
