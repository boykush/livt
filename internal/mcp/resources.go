package mcp

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/boykush/livt/internal/uri"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerResources exposes the livt repository as addressable resources, so a client
// reads the spec by URI (opportunity -> story map -> story -> mapping -> rule ->
// example, with the canvas, questions, and ubiquitous terms linked alongside)
// rather than calling a tool.
// Only resource templates are advertised — no concrete resources and no
// subscribe — keeping the server stateless with every read served fresh.
func (s *Server) registerResources(srv *mcpsdk.Server) {
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "example-mapping",
		Title:       "Example mapping",
		Description: "An example mapping (rules, examples, questions, ubiquitous terms) for a story, addressed by story key. Every rule, example, and question inside carries its own uri — the address to read next or to cite.",
		MIMEType:    "application/json",
		URITemplate: uri.MappingTemplate,
	}, s.readMapping)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "rule",
		Title:       "Rule",
		Description: "A single rule, its examples, and its automation record (issues, automated) from a story's example mapping. Rule ids restart in every mapping, so the whole uri — story key included — is what addresses this rule. A retired rule resolves too, carrying retired: true and, when something took its place, superseded_by: the uris to read next.",
		MIMEType:    "application/json",
		URITemplate: uri.RuleTemplate,
	}, s.readRule)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "example",
		Title:       "Example",
		Description: "A single example of a rule. Example ids are numbered within their rule, so the address carries the rule id: EX-01 names nothing on its own. A retired example resolves too, carrying retired: true and any superseded_by uris.",
		MIMEType:    "application/json",
		URITemplate: uri.ExampleTemplate,
	}, s.readExample)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "question",
		Title:       "Question",
		Description: "A single question from a story's example mapping. Questions hang off the mapping rather than off a rule, so the address stops at the story key. A retired question resolves too, carrying retired: true; a question settled into a rule carries that rule in superseded_by, which is where its answer landed.",
		MIMEType:    "application/json",
		URITemplate: uri.QuestionTemplate,
	}, s.readQuestion)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "opportunity",
		Title:       "Opportunity",
		Description: "One thing the product could take on — a user problem together with the business benefit of solving it — addressed by opportunity key. Carries the statement, the uri of its canvas when one has been filled in, and the story maps mapped for it. Read this to learn why a story map exists before reading what it maps.",
		MIMEType:    "application/json",
		URITemplate: uri.OpportunityTemplate,
	}, s.readOpportunity)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "opportunity-canvas",
		Title:       "Opportunity canvas",
		Description: "The Opportunity Canvas filled in for an opportunity: ten boxes holding the problems, the users, what they do today, the business challenges, the value and metrics, the adoption strategy, and the budget. Every box is returned, empty ones included — a blank box records a question the opportunity has not answered yet.",
		MIMEType:    "application/json",
		URITemplate: uri.OpportunityCanvasTemplate,
	}, s.readOpportunityCanvas)
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
		Description: "A story's name, body, and frontmatter meta, addressed by story key, plus the uris of its example mapping and of the opportunities it sits on.",
		MIMEType:    "application/json",
		URITemplate: uri.StoryTemplate,
	}, s.readStory)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "ubiquitous-term",
		Title:       "Ubiquitous language term",
		Description: "A ubiquitous language term's name and definition, addressed by term key. This shape addresses a term that holds across contexts; one scoped to a single context is addressed by ubiquitous-term-in-context instead.",
		MIMEType:    "application/json",
		URITemplate: uri.TermTemplate,
	}, s.readTerm)
	srv.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		Name:        "ubiquitous-term-in-context",
		Title:       "Ubiquitous language term in a context",
		Description: "A ubiquitous language term scoped to one context, addressed by context and term key. A context is optional and bounds where the term's meaning holds, so the same term key can name one thing across contexts and another inside one — the ctx is part of what addresses this term, not a label on it.",
		MIMEType:    "application/json",
		URITemplate: uri.ScopedTermTemplate,
	}, s.readTerm)
}

// readMapping serves livt://mapping/{story_key}.
func (s *Server) readMapping(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return s.read(req.Params.URI, uri.KindMapping)
}

// readRule serves livt://mapping/{story_key}/rule/{rule_id}.
func (s *Server) readRule(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return s.read(req.Params.URI, uri.KindRule)
}

// readExample serves livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}.
func (s *Server) readExample(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return s.read(req.Params.URI, uri.KindExample)
}

// readQuestion serves livt://mapping/{story_key}/question/{question_id}.
func (s *Server) readQuestion(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return s.read(req.Params.URI, uri.KindQuestion)
}

// readOpportunity serves livt://opportunity/{opportunity_key}.
func (s *Server) readOpportunity(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return s.read(req.Params.URI, uri.KindOpportunity)
}

// readOpportunityCanvas serves livt://opportunity-canvas/{opportunity_key}.
func (s *Server) readOpportunityCanvas(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return s.read(req.Params.URI, uri.KindOpportunityCanvas)
}

// readStoryMap serves livt://story-map/{map_name}.
func (s *Server) readStoryMap(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return s.read(req.Params.URI, uri.KindStoryMap)
}

// readStory serves livt://story/{story_key}.
func (s *Server) readStory(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return s.read(req.Params.URI, uri.KindStory)
}

// readTerm serves both term shapes: livt://ubiquitous/{term_key} and
// livt://ubiquitous/{ctx}/{term_key}. They resolve the same way — the context is
// part of the address, not a different kind of thing — so one handler answers
// for the two templates.
func (s *Server) readTerm(_ context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return s.read(req.Params.URI, uri.KindTerm)
}

// read serves one resource. The URI has to parse as the shape its own template
// advertises, so a handler never answers for a neighbouring shape. Resolving
// goes through Resolve — the same entry the CLI uses, which is what keeps the
// two surfaces on one JSON shape.
func (s *Server) read(resURI string, want uri.Kind) (*mcpsdk.ReadResourceResult, error) {
	p, ok := uri.Parse(resURI)
	if !ok || p.Kind != want {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	payload, err := s.Resolve(p)
	if errors.Is(err, ErrNotFound) {
		return nil, mcpsdk.ResourceNotFoundError(resURI)
	}
	if err != nil {
		return nil, err
	}
	return jsonResource(resURI, payload)
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
