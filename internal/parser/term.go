package parser

import (
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
	"gopkg.in/yaml.v3"
)

type termFrontmatter struct {
	Name string `yaml:"name"`
}

func ParseTerm(path string) (*domain.Term, error) {
	key, frontmatter, body, err := splitFrontmatter(path)
	if err != nil {
		return nil, err
	}

	var fm termFrontmatter
	if frontmatter != "" {
		if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
			return nil, err
		}
	}

	return &domain.Term{
		Key:  key,
		Name: fm.Name,
		Body: body,
	}, nil
}

func ParseAllTerms(ubiquitousDir string) ([]*domain.Term, error) {
	files, err := filepath.Glob(filepath.Join(ubiquitousDir, "*.md"))
	if err != nil {
		return nil, err
	}

	var terms []*domain.Term
	for _, f := range files {
		term, err := ParseTerm(f)
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}

	return terms, nil
}
