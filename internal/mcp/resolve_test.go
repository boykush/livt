package mcp

import (
	"context"
	"encoding/json"
	"errors"

	"testing"

	"github.com/boykush/livt/internal/uri"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// everyKind is one URI of each shape the livt repository can address, against the
// fixture newTestServer builds.
var everyKind = []string{
	"livt://mapping/demo",
	"livt://mapping/demo/rule/R-01",
	"livt://mapping/demo/rule/R-01/example/EX-01",
	"livt://mapping/demo/question/Q-01",
	"livt://story-map/%E3%83%87%E3%83%A2%E3%83%9E%E3%83%83%E3%83%97",
	"livt://story/demo",
	"livt://ubiquitous/story",
}

// livt://mapping/trace-test-to-rule/rule/R-04/example/EX-02: what the CLI
// prints is what an MCP resources/read serves. Both go through Resolve, and
// this drives a real session to hold that to the byte.
func TestResolveMatchesTheMCPResourceRead(t *testing.T) {
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

	for _, raw := range everyKind {
		read, err := cs.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: raw})
		if err != nil {
			t.Errorf("resources/read %s: %v", raw, err)
			continue
		}
		p, ok := uri.Parse(raw)
		if !ok {
			t.Errorf("Parse(%q) failed", raw)
			continue
		}
		payload, err := s.Resolve(p)
		if err != nil {
			t.Errorf("Resolve(%s): %v", raw, err)
			continue
		}
		body, err := json.Marshal(payload)
		if err != nil {
			t.Errorf("marshal %s: %v", raw, err)
			continue
		}
		if got, want := string(body), read.Contents[0].Text; got != want {
			t.Errorf("%s\n CLI: %s\n MCP: %s", raw, got, want)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-04/example/EX-01: every kind the
// livt repository holds resolves, not just the rules.
func TestResolveEveryKind(t *testing.T) {
	s := newTestServer(t)
	for _, raw := range everyKind {
		p, ok := uri.Parse(raw)
		if !ok {
			t.Errorf("Parse(%q) failed", raw)
			continue
		}
		if _, err := s.Resolve(p); err != nil {
			t.Errorf("Resolve(%s): %v", raw, err)
		}
	}
}

// A URI whose shape is fine but whose item is absent is not the same failure as
// a livt repository that will not read, so only the first is ErrNotFound.
func TestResolveReportsMissingItemsAsNotFound(t *testing.T) {
	s := newTestServer(t)
	for _, raw := range []string{
		"livt://mapping/nope",
		"livt://mapping/demo/rule/R-99",
		"livt://mapping/demo/rule/R-01/example/EX-99",
		"livt://mapping/demo/rule/R-02/example/EX-01", // right example id, wrong rule
		"livt://mapping/demo/question/Q-99",
		"livt://story-map/%E5%AD%98%E5%9C%A8%E3%81%97%E3%81%AA%E3%81%84",
		"livt://story/nope",
		"livt://ubiquitous/nope",
	} {
		p, ok := uri.Parse(raw)
		if !ok {
			t.Errorf("Parse(%q) failed", raw)
			continue
		}
		if _, err := s.Resolve(p); !errors.Is(err, ErrNotFound) {
			t.Errorf("Resolve(%s) error = %v, want ErrNotFound", raw, err)
		}
	}
}
