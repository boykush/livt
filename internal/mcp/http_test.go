package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestHTTPTransport drives the server over a real Streamable HTTP endpoint, the
// way a consumer repo's client connects to a shared local server: a tool call
// and a resource read, plus the mount path. DisableStandaloneSSE mirrors a
// request/response client, since the stateless server pushes no notifications.
func TestHTTPTransport(t *testing.T) {
	ctx := context.Background()
	s := newTestServer(t)

	httpServer := httptest.NewServer(s.httpHandler())
	defer httpServer.Close()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "test-client", Version: "0"}, nil)
	cs, err := client.Connect(ctx, &mcpsdk.StreamableClientTransport{
		Endpoint:             httpServer.URL + httpPath,
		DisableStandaloneSSE: true,
	}, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer cs.Close()

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

	// A resource read works over HTTP too.
	mr, err := cs.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: "livt://mapping/demo"})
	if err != nil {
		t.Fatalf("read mapping: %v", err)
	}
	var em exampleMappingResult
	if err := json.Unmarshal([]byte(resourceText(t, mr)), &em); err != nil {
		t.Fatalf("decode mapping: %v", err)
	}
	if em.Mapping.StoryKey != "demo" {
		t.Errorf("mapping story key = %q, want demo", em.Mapping.StoryKey)
	}

	// The handler is mounted only at httpPath; other paths are not served.
	resp, err := http.Get(httpServer.URL + "/")
	if err != nil {
		t.Fatalf("get /: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET / status = %d, want 404", resp.StatusCode)
	}
}
