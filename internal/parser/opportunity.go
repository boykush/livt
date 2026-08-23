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

// FindOpportunityByKey falls back to a placeholder named after the key, the
// same way FindStoryByKey does: a board referencing an opportunity that is not
// committed yet should still render the reference rather than fail the build.
func FindOpportunityByKey(opportunitiesDir string, key domain.OpportunityKey) (*domain.Opportunity, error) {
	path := filepath.Join(opportunitiesDir, key.Value+".md")
	opportunity, err := ParseOpportunity(path)
	if err != nil {
		return &domain.Opportunity{Key: key, Name: key.Value}, nil
	}
	return opportunity, nil
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
