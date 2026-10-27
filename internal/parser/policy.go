package parser

import (
	"fmt"
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/uri"
	"gopkg.in/yaml.v3"
)

type policyFrontmatter struct {
	Name         string   `yaml:"name"`
	Status       string   `yaml:"status"`
	Input        []string `yaml:"input"`
	Output       []string `yaml:"output"`
	SupersededBy []string `yaml:"superseded_by"`
}

// ParsePolicy reads the policy a reference names, relative to the policies
// directory: "{ctx}/{policy-key}" for a policy scoped to one context,
// "{policy-key}" for one that holds across them. The reference is guarded
// segment by segment, so an externally-supplied one cannot walk out of the
// directory.
func ParsePolicy(policiesDir, ref string) (*domain.Policy, error) {
	ctx, key, ok := uri.SplitPolicyRef(ref)
	if !ok {
		return nil, fmt.Errorf("invalid policy reference %q", ref)
	}

	_, frontmatter, body, err := splitFrontmatter(filepath.Join(policiesDir, ctx, key+".md"))
	if err != nil {
		return nil, err
	}

	var fm policyFrontmatter
	if frontmatter != "" {
		if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
			return nil, err
		}
	}

	// An unknown status is refused rather than read as the default, the way a
	// rule's is: read as accepted, a mistyped retired would put a withdrawn
	// policy back in force.
	status := domain.RuleStatusDefault
	if fm.Status != "" {
		status = domain.RuleStatus(fm.Status)
		if !status.Valid() {
			return nil, fmt.Errorf("policy %q: unknown status %q (supported: %s)", ref, fm.Status, domain.RuleStatusList())
		}
	}

	return &domain.Policy{
		Ctx:          ctx,
		Key:          key,
		Name:         fm.Name,
		Status:       status,
		Input:        fm.Input,
		Output:       fm.Output,
		SupersededBy: fm.SupersededBy,
		Body:         body,
	}, nil
}
