package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
)

// buildPolicies renders every policy still in force as a row of policies.html,
// anchored by its reference the way a glossary row is, so livt://policy/… lands
// on it. A closed policy is off the table for the reason a closed rule is off
// its board: nobody is asked to follow it any more.
func (b *Builder) buildPolicies() error {
	policies, err := parser.ParseAllPolicies(b.PoliciesDir)
	if err != nil {
		return err
	}

	var rows []policyRow
	var contexts []string
	seen := make(map[string]bool)
	for _, p := range policies {
		if !p.Status.Active() {
			continue
		}
		rows = append(rows, policyRow{
			Anchor:   uri.PolicyAnchor(p.Ctx, p.Key),
			Ctx:      p.Ctx,
			Key:      p.Key,
			Name:     orKey(p.Name, p.Key),
			Proposed: p.Status.OrDefault() == domain.RuleProposed,
			Body:     b.markdown(p.Body, filepath.Join(b.PoliciesDir, p.Ctx), ""),
		})
		if p.Ctx != "" && !seen[p.Ctx] {
			seen[p.Ctx] = true
			contexts = append(contexts, p.Ctx)
		}
	}

	sb, err := b.sidebar("policies", "")
	if err != nil {
		return err
	}

	outPath := filepath.Join(b.OutDir, "policies.html")
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := renderPolicies(f, b.Lang, policiesView{
		Sidebar:  sb,
		Policies: rows,
		Contexts: contexts,
	}); err != nil {
		return err
	}
	fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	return nil
}
