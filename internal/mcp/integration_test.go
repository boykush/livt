package mcp

import (
	"context"
	"encoding/json"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestEndToEnd drives the server through a real MCP session over in-memory
// transports: the list tools for discovery, then reading the master as
// resources by URI (story map -> story -> mapping -> rule, plus a ubiquitous
// term), plus the not-found paths.
func TestEndToEnd(t *testing.T) {
	ctx := context.Background()
	s := newTestServer(t)

	serverTransport, clientTransport := mcpsdk.NewInMemoryTransports()
	serverSession, err := s.mcpServer().Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer serverSession.Close()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "test-client", Version: "0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer cs.Close()

	// Discovery: the list tools hand out resource URIs.
	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if names := toolNames(tools.Tools); len(names) != 2 || !hasName(names, "list_stories") || !hasName(names, "list_story_maps") {
		t.Fatalf("tools = %v, want [list_stories list_story_maps]", names)
	}

	res, err := cs.CallTool(ctx, &mcpsdk.CallToolParams{Name: "list_stories", Arguments: map[string]any{}})
	if err != nil {
		t.Fatalf("call list_stories: %v", err)
	}
	var list listStoriesOutput
	if err := json.Unmarshal([]byte(contentText(res)), &list); err != nil {
		t.Fatalf("decode list_stories: %v", err)
	}
	if uri := mappingURIOf(list.Stories, "demo"); uri != "livt://mapping/demo" {
		t.Errorf("demo example_mapping_uri = %q, want livt://mapping/demo", uri)
	}
	if uri := storyURIOf(list.Stories, "demo"); uri != "livt://story/demo" {
		t.Errorf("demo uri = %q, want livt://story/demo", uri)
	}

	// Resources: every template is advertised, and only templates (no concrete
	// resources, no subscribe) — the server stays stateless.
	tmpls, err := cs.ListResourceTemplates(ctx, nil)
	if err != nil {
		t.Fatalf("list resource templates: %v", err)
	}
	for _, want := range []string{mappingURITemplate, ruleURITemplate, storyMapURITemplate, storyURITemplate, termURITemplate} {
		if !hasTemplate(tmpls.ResourceTemplates, want) {
			t.Errorf("resource template %q not advertised", want)
		}
	}

	// Read the mapping resource; each rule links to its own rule resource and
	// each referenced ubiquitous term resolves to its resource.
	mr, err := cs.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: "livt://mapping/demo"})
	if err != nil {
		t.Fatalf("read mapping: %v", err)
	}
	var em exampleMappingResult
	if err := json.Unmarshal([]byte(resourceText(t, mr)), &em); err != nil {
		t.Fatalf("decode mapping: %v", err)
	}
	if em.Mapping.StoryKey != "demo" || len(em.Mapping.Rules) != 2 {
		t.Fatalf("mapping = %+v, want demo with two rules", em.Mapping)
	}
	if em.Mapping.Rules[0].URI != "livt://mapping/demo/rule/R-01" {
		t.Errorf("rule uri = %q, want livt://mapping/demo/rule/R-01", em.Mapping.Rules[0].URI)
	}
	if len(em.Mapping.UbiquitousTerms) != 2 || em.Mapping.UbiquitousTerms[0].URI != "livt://ubiquitous/story" {
		t.Errorf("ubiquitous_terms = %+v, want story resolved to livt://ubiquitous/story", em.Mapping.UbiquitousTerms)
	}

	// Read the rule resource directly.
	rr, err := cs.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: "livt://mapping/demo/rule/R-01"})
	if err != nil {
		t.Fatalf("read rule: %v", err)
	}
	var rule ruleResult
	if err := json.Unmarshal([]byte(resourceText(t, rr)), &rule); err != nil {
		t.Fatalf("decode rule: %v", err)
	}
	if rule.Rule.ID != "R-01" || rule.Rule.Name != "ルール1" {
		t.Errorf("rule = %+v, want R-01 ルール1", rule.Rule)
	}
	if len(rule.Rule.Issues) != 1 || !rule.Rule.Automated {
		t.Errorf("rule = %+v, want the recorded issue URL and automated=true over the wire", rule.Rule)
	}

	// Discover story maps, then read one through the URI the tool handed out —
	// the display name survives the percent-encoded round trip.
	res, err = cs.CallTool(ctx, &mcpsdk.CallToolParams{Name: "list_story_maps", Arguments: map[string]any{}})
	if err != nil {
		t.Fatalf("call list_story_maps: %v", err)
	}
	var mapsList listStoryMapsOutput
	if err := json.Unmarshal([]byte(contentText(res)), &mapsList); err != nil {
		t.Fatalf("decode list_story_maps: %v", err)
	}
	if len(mapsList.StoryMaps) != 1 || mapsList.StoryMaps[0].Name != "デモマップ" {
		t.Fatalf("story maps = %+v, want one デモマップ", mapsList.StoryMaps)
	}
	smr, err := cs.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: mapsList.StoryMaps[0].URI})
	if err != nil {
		t.Fatalf("read story map %q: %v", mapsList.StoryMaps[0].URI, err)
	}
	var sm storyMapResult
	if err := json.Unmarshal([]byte(resourceText(t, smr)), &sm); err != nil {
		t.Fatalf("decode story map: %v", err)
	}
	if sm.StoryMap.Name != "デモマップ" {
		t.Errorf("story map name = %q, want デモマップ", sm.StoryMap.Name)
	}
	cards := sm.StoryMap.Activities[0].Steps[0].Stories
	if len(cards) != 2 || cards[0].URI != "livt://story/demo" || cards[1].URI != "" {
		t.Errorf("cards = %+v, want committed demo linked and candidate bare", cards)
	}

	// Follow the card's URI to the story resource: name, body, and meta.
	sr, err := cs.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: "livt://story/demo"})
	if err != nil {
		t.Fatalf("read story: %v", err)
	}
	var story storyResult
	if err := json.Unmarshal([]byte(resourceText(t, sr)), &story); err != nil {
		t.Fatalf("decode story: %v", err)
	}
	if story.Story.Name != "デモストーリー" || story.Story.Body != "本文" {
		t.Errorf("story = %+v, want デモストーリー with body 本文", story.Story)
	}
	if len(story.Story.Meta) != 1 || story.Story.Meta[0].Key != "issue" {
		t.Errorf("story meta = %+v, want one issue field", story.Story.Meta)
	}
	if story.Story.ExampleMappingURI != "livt://mapping/demo" {
		t.Errorf("story example_mapping_uri = %q, want livt://mapping/demo", story.Story.ExampleMappingURI)
	}

	// Follow a resolved term ref to the ubiquitous term resource.
	tr, err := cs.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: "livt://ubiquitous/story"})
	if err != nil {
		t.Fatalf("read term: %v", err)
	}
	var term termResult
	if err := json.Unmarshal([]byte(resourceText(t, tr)), &term); err != nil {
		t.Fatalf("decode term: %v", err)
	}
	if term.Term.Name != "ストーリー" || term.Term.Body != "ストーリーの定義" {
		t.Errorf("term = %+v, want ストーリー with its definition", term.Term)
	}

	// Unknown URIs are resource errors.
	for _, uri := range []string{
		"livt://mapping/nope",
		"livt://mapping/demo/rule/R-99",
		storyMapURI("なし"),
		"livt://story/nope",
		"livt://ubiquitous/nope",
	} {
		if _, err := cs.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: uri}); err == nil {
			t.Errorf("expected error reading %q", uri)
		}
	}
}

func contentText(res *mcpsdk.CallToolResult) string {
	if len(res.Content) == 0 {
		return ""
	}
	if tc, ok := res.Content[0].(*mcpsdk.TextContent); ok {
		return tc.Text
	}
	return ""
}

func resourceText(t *testing.T, res *mcpsdk.ReadResourceResult) string {
	t.Helper()
	if len(res.Contents) != 1 || res.Contents[0].MIMEType != "application/json" {
		t.Fatalf("resource contents = %+v, want one application/json item", res.Contents)
	}
	return res.Contents[0].Text
}

func toolNames(tools []*mcpsdk.Tool) []string {
	names := make([]string, len(tools))
	for i, tool := range tools {
		names[i] = tool.Name
	}
	return names
}

func hasTemplate(tmpls []*mcpsdk.ResourceTemplate, uriTemplate string) bool {
	for _, tm := range tmpls {
		if tm.URITemplate == uriTemplate {
			return true
		}
	}
	return false
}

func hasName(names []string, want string) bool {
	for _, n := range names {
		if n == want {
			return true
		}
	}
	return false
}

func mappingURIOf(stories []storySummaryJSON, key string) string {
	for _, st := range stories {
		if st.Key == key {
			return st.ExampleMappingURI
		}
	}
	return ""
}

func storyURIOf(stories []storySummaryJSON, key string) string {
	for _, st := range stories {
		if st.Key == key {
			return st.URI
		}
	}
	return ""
}
