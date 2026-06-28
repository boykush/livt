package mcp

import (
	"context"
	"encoding/json"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestEndToEnd drives the server through a real MCP session over in-memory
// transports: the list_stories tool for discovery, then reading a mapping and a
// rule as resources by URI (story -> mapping -> rule), plus the not-found paths.
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

	// Discovery: the only tool is list_stories, and it links to the mapping resource.
	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(tools.Tools) != 1 || tools.Tools[0].Name != "list_stories" {
		t.Fatalf("tools = %v, want [list_stories]", toolNames(tools.Tools))
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

	// Resources: both templates are advertised.
	tmpls, err := cs.ListResourceTemplates(ctx, nil)
	if err != nil {
		t.Fatalf("list resource templates: %v", err)
	}
	for _, want := range []string{mappingURITemplate, ruleURITemplate} {
		if !hasTemplate(tmpls.ResourceTemplates, want) {
			t.Errorf("resource template %q not advertised", want)
		}
	}

	// Read the mapping resource; each rule links to its own rule resource.
	mr, err := cs.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: "livt://mapping/demo"})
	if err != nil {
		t.Fatalf("read mapping: %v", err)
	}
	var em exampleMappingResult
	if err := json.Unmarshal([]byte(resourceText(t, mr)), &em); err != nil {
		t.Fatalf("decode mapping: %v", err)
	}
	if em.Mapping.StoryKey != "demo" || len(em.Mapping.Rules) != 1 {
		t.Fatalf("mapping = %+v, want demo with one rule", em.Mapping)
	}
	if em.Mapping.Rules[0].URI != "livt://mapping/demo/rule/R-01" {
		t.Errorf("rule uri = %q, want livt://mapping/demo/rule/R-01", em.Mapping.Rules[0].URI)
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

	// Unknown URIs are resource errors.
	for _, uri := range []string{"livt://mapping/nope", "livt://mapping/demo/rule/R-99"} {
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

func mappingURIOf(stories []storySummaryJSON, key string) string {
	for _, st := range stories {
		if st.Key == key {
			return st.ExampleMappingURI
		}
	}
	return ""
}
