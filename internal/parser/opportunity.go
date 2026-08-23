package parser

import (
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
)

func ParseOpportunity(path string) (*domain.Opportunity, error) {
	key, frontmatter, body, err := splitFrontmatter(path)
	if err != nil {
		return nil, err
	}

	name, meta, err := parseNamedFrontmatter(frontmatter)
	if err != nil {
		return nil, err
	}

	return &domain.Opportunity{
		Key:  domain.OpportunityKey{Value: key},
		Name: name,
		Body: body,
		Meta: meta,
	}, nil
}

func ParseAllOpportunities(opportunitiesDir string) ([]*domain.Opportunity, error) {
	files, err := filepath.Glob(filepath.Join(opportunitiesDir, "*.md"))
	if err != nil {
		return nil, err
	}

	var opportunities []*domain.Opportunity
	for _, f := range files {
		opportunity, err := ParseOpportunity(f)
		if err != nil {
			return nil, err
		}
		opportunities = append(opportunities, opportunity)
	}

	return opportunities, nil
}
