package mcp

import (
	"context"
	"encoding/json"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestToolsEndToEnd drives the server through a real MCP session over in-memory
// transports: handshake, tools/list, and tools/call for both a success and an
// error case. This exercises schema validation and result packing that the
// direct handler tests skip.
func TestToolsEndToEnd(t *testing.T) {
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

	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(tools.Tools) != 3 {
		t.Errorf("tools advertised = %d, want 3", len(tools.Tools))
	}

	// Success: get_rule returns the requested rule.
	res, err := cs.CallTool(ctx, &mcpsdk.CallToolParams{
		Name:      "get_rule",
		Arguments: map[string]any{"story_key": "demo", "rule_id": "R-01"},
	})
	if err != nil {
		t.Fatalf("call get_rule: %v", err)
	}
	if res.IsError {
		t.Fatalf("get_rule reported error: %s", contentText(res))
	}
	var out getRuleOutput
	if err := json.Unmarshal([]byte(contentText(res)), &out); err != nil {
		t.Fatalf("decode get_rule result: %v", err)
	}
	if out.Rule.ID != "R-01" || out.Rule.Name != "ルール1" {
		t.Errorf("rule = %+v, want R-01 ルール1", out.Rule)
	}

	// Error: an unknown rule is surfaced as a tool error, not a transport error.
	res, err = cs.CallTool(ctx, &mcpsdk.CallToolParams{
		Name:      "get_rule",
		Arguments: map[string]any{"story_key": "demo", "rule_id": "R-99"},
	})
	if err != nil {
		t.Fatalf("call get_rule (unknown): %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true for unknown rule id")
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
