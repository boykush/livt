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

	return parsePolicyFile(filepath.Join(policiesDir, ctx, key+".md"), ctx)
}

// parsePolicyFile reads one policy file, taking its context from the caller for
// the reason parseTermFile does.
func parsePolicyFile(path, ctx string) (*domain.Policy, error) {
	key, frontmatter, body, err := splitFrontmatter(path)
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
			return nil, fmt.Errorf("policy %q: unknown status %q (supported: %s)", uri.PolicyRef(ctx, key), fm.Status, domain.RuleStatusList())
		}
	}

	return &domain.Policy{
		Ctx:          ctx,
		Key:          key,
		Name:         fm.Name,
		Status:       status,
		SupersededBy: fm.SupersededBy,
		Body:         body,
	}, nil
}

// ParseAllPolicies reads every policy in the policies directory, bounded and
// ordered the way ParseAllTerms bounds and orders terms. A missing directory
// holds no policies rather than failing: most livt repositories have none.
func ParseAllPolicies(policiesDir string) ([]*domain.Policy, error) {
	rooted, err := filepath.Glob(filepath.Join(policiesDir, "*.md"))
	if err != nil {
		return nil, err
	}
	scoped, err := filepath.Glob(filepath.Join(policiesDir, "*", "*.md"))
	if err != nil {
		return nil, err
	}

	var policies []*domain.Policy
	for _, f := range rooted {
		policy, err := parsePolicyFile(f, "")
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}
	for _, f := range scoped {
		policy, err := parsePolicyFile(f, filepath.Base(filepath.Dir(f)))
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}
	return policies, nil
}
