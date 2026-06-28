package mcp

import (
	"context"
	"encoding/json"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	mappingURIPrefix   = "livt://mapping/"
	ruleURIInfix       = "/rule/"
	mappingURITemplate = "livt://mapping/{story_key}"
	ruleURITemplate    = "livt://mapping/{story_key}/rule/{rule_id}"
)

// mappingURI builds the resource URI for a story's example mapping.
func mappingURI(storyKey string) string {
	return mappingURIPrefix + storyKey
}

// ruleURI builds the resource URI for a single rule within a story's mapping.
func ruleURI(storyKey, ruleID string) string {
	return mappingURIPrefix + storyKey + ruleURIInfix + ruleID
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

// validSegment guards externally-supplied URI/key segments so they resolve to a
// file inside the mappings directory and cannot traverse out of it.
func validSegment(s string) bool {
	return s != "" && s != "." && s != ".." &&
		!strings.ContainsAny(s, `/\`) && !strings.Contains(s, "..")
}

// registerResources exposes example mappings and their rules as addressable
// resources, so a client reads a story's spec by URI (story -> mapping -> rule)
// rather than calling a tool.
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
	return jsonResource(uri, exampleMappingResult{versioned: s.versioned(), Mapping: toExampleMappingJSON(em)})
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
