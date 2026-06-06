package parser

import (
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
)

func ParseTerm(path string) (*domain.Term, error) {
	key, name, body, err := parseMarkdownDoc(path)
	if err != nil {
		return nil, err
	}

	return &domain.Term{
		Key:  key,
		Name: name,
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
