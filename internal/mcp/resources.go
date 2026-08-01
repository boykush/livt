package mcp

import (
	"context"
	"encoding/json"

	"github.com/boykush/livt/internal/uri"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerResources exposes the master as addressable resources, so a client
// reads the spec by URI (story map -> story -> mapping -> rule -> example, with
// questions and ubiquitous terms linked alongside) rather than calling a tool.
// Only resource templates are advertised — no concrete resources and no
// subscribe — keeping the server stateless with every read served fresh.
func (s *Server) registerResources(srv *mcpsdk.Server) {
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "example-mapping",
		Title:       "Example mapping",
		Description: "An example mapping (rules, examples, questions, ubiquitous terms) for a story, addressed by story key.",
		MIMEType:    "application/json",
		URITemplate: uri.MappingTemplate,
	}, s.readMapping)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "rule",
		Title:       "Rule",
		Description: "A single rule and its examples from a story's example mapping. A retired rule resolves too, carrying retired: true.",
		MIMEType:    "application/json",
		URITemplate: uri.RuleTemplate,
	}, s.readRule)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "example",
		Title:       "Example",
		Description: "A single example of a rule. Example ids are numbered within their rule, so the address carries the rule id.",
		MIMEType:    "application/json",
		URITemplate: uri.ExampleTemplate,
	}, s.readExample)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "question",
		Title:       "Question",
		Description: "A single question from a story's example mapping. A retired question resolves too, carrying retired: true.",
		MIMEType:    "application/json",
		URITemplate: uri.QuestionTemplate,
	}, s.readQuestion)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "story-map",
		Title:       "Story map",
		Description: "A user story map (activities, steps, story cards, releases, ubiquitous terms), addressed by its display name (percent-encoded).",
		MIMEType:    "application/json",
		URITemplate: uri.StoryMapTemplate,
	}, s.readStoryMap)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "story",
		Title:       "Story",
		Description: "A story's name, body, and frontmatter meta, addressed by story key.",
		MIMEType:    "application/json",
		URITemplate: uri.StoryTemplate,
	}, s.readStory)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "ubiquitous-term",
		Title:       "Ubiquitous language term",
		Description: "A ubiquitous language term's name and definition, addressed by term key.",
		MIMEType:    "application/json",
		URITemplate: uri.TermTemplate,
	}, s.readTerm)
}

// readMapping serves livt://mapping/{story_key}.
func (s *Server) readMapping(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	resURI := req.Params.URI
	key, ok := uri.ParseMapping(resURI)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	em, err := s.cfg.exampleMapping(key)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	return jsonResource(resURI, exampleMappingResult{versioned: s.versioned(), Mapping: s.cfg.toExampleMappingJSON(em)})
}

// readRule serves livt://mapping/{story_key}/rule/{rule_id}.
func (s *Server) readRule(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	resURI := req.Params.URI
	key, ruleID, ok := uri.ParseRule(resURI)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	rule, err := s.cfg.rule(key, ruleID)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	return jsonResource(resURI, ruleResult{versioned: s.versioned(), Rule: toRuleJSON(key, rule)})
}

// readExample serves livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}.
func (s *Server) readExample(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	resURI := req.Params.URI
	key, ruleID, exampleID, ok := uri.ParseExample(resURI)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	example, err := s.cfg.example(key, ruleID, exampleID)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	return jsonResource(resURI, exampleResult{versioned: s.versioned(), Example: toExampleJSON(key, ruleID, example)})
}

// readQuestion serves livt://mapping/{story_key}/question/{question_id}.
func (s *Server) readQuestion(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	resURI := req.Params.URI
	key, questionID, ok := uri.ParseQuestion(resURI)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	question, err := s.cfg.question(key, questionID)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	return jsonResource(resURI, questionResult{versioned: s.versioned(), Question: toQuestionJSON(key, question)})
}

// readStoryMap serves livt://story-map/{map_name}.
func (s *Server) readStoryMap(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	resURI := req.Params.URI
	name, ok := uri.ParseStoryMap(resURI)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	sm, err := s.cfg.storyMap(name)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	return jsonResource(resURI, storyMapResult{versioned: s.versioned(), StoryMap: s.cfg.toStoryMapJSON(sm)})
}

// readStory serves livt://story/{story_key}.
func (s *Server) readStory(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	resURI := req.Params.URI
	key, ok := uri.ParseStory(resURI)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	story, err := s.cfg.story(key)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	out, err := s.cfg.toStoryJSON(story)
	if err != nil {
		return nil, err
	}
	return jsonResource(resURI, storyResult{versioned: s.versioned(), Story: out})
}

// readTerm serves livt://ubiquitous/{term_key}.
func (s *Server) readTerm(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	resURI := req.Params.URI
	key, ok := uri.ParseTerm(resURI)
	if !ok {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	term, err := s.cfg.term(key)
	if err != nil {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	return jsonResource(resURI, termResult{versioned: s.versioned(), Term: toTermJSON(term)})
}

// jsonResource marshals v as the single application/json content of a resource.
func jsonResource(resURI string, v any) (*mcpsdk.ReadResourceResult, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{{
			URI:      resURI,
			MIMEType: "application/json",
			Text:     string(body),
		}},
	}, nil
}
