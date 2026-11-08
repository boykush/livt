package builder

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/boykush/livt/internal/mcp"
	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
)

// checkSuccessors fails the build on a superseded_by that names nothing in the
// livt repository. A retired item is off the board, so a broken pointer on it
// would otherwise go unseen until someone follows it. Every one is reported at
// once, so the writer fixes them in one change.
// livt://mapping/trace-test-to-rule/rule/R-11
func (b *Builder) checkSuccessors() error {
	mappings, err := parser.ParseAllExampleMappings(b.MappingsDir)
	if err != nil {
		return err
	}
	// Resolved the way `livt resolve` resolves, so the build cannot accept a
	// successor the CLI and the MCP server then fail to find.
	cfg := mcp.Config{Dirs: &mcp.Dirs{
		Opportunities: b.OpportunitiesDir,
		Canvases:      b.CanvasesDir,
		Mappings:      b.MappingsDir,
		Stories:       b.StoriesDir,
		USM:           b.USMDir,
		Ubiquitous:    b.UbiquitousDir,
		Automations:   b.AutomationsDir,
	}}
	var errs []error
	for _, em := range mappings {
		file := filepath.Join(b.MappingsDir, em.StoryKey.Value+".yaml")
		check := func(item string, successors []string) {
			for _, s := range successors {
				if err := verifySuccessor(cfg, s); err != nil {
					errs = append(errs, fmt.Errorf("%s: %s: superseded_by %s: %w", file, item, s, err))
				}
			}
		}
		for _, r := range em.Rules {
			check("rule "+r.ID, r.SupersededBy)
			for _, ex := range r.Examples {
				check("example "+r.ID+"/"+ex.ID, ex.SupersededBy)
			}
		}
		for _, q := range em.Questions {
			check("question "+q.ID, q.SupersededBy)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("superseded_by names what the livt repository does not hold:\n%w", errors.Join(errs...))
	}
	return nil
}

func verifySuccessor(cfg mcp.Config, raw string) error {
	p, ok := uri.Parse(raw)
	if !ok {
		return errors.New("not a livt URI")
	}
	return cfg.Verify(p)
}
