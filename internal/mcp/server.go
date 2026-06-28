// Package mcp serves the discovery master (stories, example mappings) over the
// Model Context Protocol so implementation repos can fetch the spec for a story
// or rule without reading livt's source. The master usually lives in a separate
// checkout from the consumer, so Config.Root locates it explicitly.
package mcp

import (
	"context"
	"path/filepath"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Config locates the discovery master. Root points at a livt project root; the
// input subdirectories are derived from it the same way build/serve lay them out.
type Config struct {
	Root string
}

func (c Config) mappingsDir() string {
	return filepath.Join(c.Root, "discoveries", "example-mappings")
}

func (c Config) storiesDir() string {
	return filepath.Join(c.Root, "stories")
}

// Server exposes the master under Config over MCP. version is the livt build
// version, reported in the MCP handshake (distinct from the per-result
// spec_version, which is the master's git revision).
type Server struct {
	cfg     Config
	version string
}

func NewServer(cfg Config, version string) *Server {
	return &Server{cfg: cfg, version: version}
}

// Run serves the tools over stdio, blocking until the client disconnects or ctx
// is cancelled.
func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer().Run(ctx, &mcpsdk.StdioTransport{})
}

// mcpServer builds the configured SDK server. Split out so tests can assert tool
// registration without standing up a transport.
func (s *Server) mcpServer() *mcpsdk.Server {
	srv := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "livt", Version: s.version}, nil)
	s.registerTools(srv)
	return srv
}
