// Package mcp serves the discovery master (stories, example mappings) over the
// Model Context Protocol so implementation repos can fetch the spec for a story
// or rule without reading livt's source. The master usually lives in a separate
// checkout from the consumer, so Config.Root locates it explicitly.
package mcp

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// httpPath is where the Streamable HTTP transport is mounted. Consumers point
// their MCP client at this path, e.g. http://localhost:5488/mcp.
const httpPath = "/mcp"

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

// RunHTTP serves the master over Streamable HTTP at addr, blocking until ctx is
// cancelled. The master is read-only and identical for every client, so the
// handler is stateless: each request is served from a temporary session with no
// retained per-client state, which lets one local server back many repos.
// Responses are plain JSON -- the spec server never pushes server-initiated
// notifications, so it needs no event stream.
func (s *Server) RunHTTP(ctx context.Context, addr string) error {
	httpSrv := &http.Server{Addr: addr, Handler: s.httpHandler()}
	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpSrv.Shutdown(shutCtx)
	}()
	if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// httpHandler builds the Streamable HTTP handler mounted at httpPath. Split out
// from RunHTTP so tests can drive it via httptest without binding a port.
func (s *Server) httpHandler() http.Handler {
	srv := s.mcpServer()
	handler := mcpsdk.NewStreamableHTTPHandler(
		func(*http.Request) *mcpsdk.Server { return srv },
		&mcpsdk.StreamableHTTPOptions{Stateless: true, JSONResponse: true},
	)
	mux := http.NewServeMux()
	mux.Handle(httpPath, handler)
	return mux
}

// mcpServer builds the configured SDK server. Split out so tests can assert tool
// registration without standing up a transport.
func (s *Server) mcpServer() *mcpsdk.Server {
	srv := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "livt", Version: s.version}, nil)
	s.registerTools(srv)
	s.registerResources(srv)
	return srv
}
