package parser

import (
	"fmt"
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/uri"
	"gopkg.in/yaml.v3"
)

type termFrontmatter struct {
	Name string `yaml:"name"`
}

// ParseTerm reads the term a reference names, relative to the ubiquitous
// directory: "{ctx}/{term-key}" for a term scoped to one context,
// "{term-key}" for one that holds across them. The reference is guarded segment
// by segment, so an externally-supplied one cannot walk out of the directory.
func ParseTerm(ubiquitousDir, ref string) (*domain.Term, error) {
	ctx, key, ok := uri.SplitTermRef(ref)
	if !ok {
		return nil, fmt.Errorf("invalid term reference %q", ref)
	}
	return parseTermFile(filepath.Join(ubiquitousDir, ctx, key+".md"), ctx)
}

// parseTermFile reads one term file. The context comes from the caller rather
// than the path: only the caller knows which directory is the ubiquitous root,
// and so which leading segment is a context rather than part of it.
func parseTermFile(path, ctx string) (*domain.Term, error) {
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
		Ctx:  ctx,
		Key:  key,
		Name: fm.Name,
		Body: body,
	}, nil
}

// ParseAllTerms reads every term in the ubiquitous directory: first the
// context-free ones at its root, then the ones under each context directory.
// The two globs are what bounds a context to a single directory level — a
// context is one path segment, the same shape its URI carries — and each glob
// sorts, so the order a build renders is the order the filesystem lays out.
func ParseAllTerms(ubiquitousDir string) ([]*domain.Term, error) {
	rooted, err := filepath.Glob(filepath.Join(ubiquitousDir, "*.md"))
	if err != nil {
		return nil, err
	}
	scoped, err := filepath.Glob(filepath.Join(ubiquitousDir, "*", "*.md"))
	if err != nil {
		return nil, err
	}

	var terms []*domain.Term
	for _, f := range rooted {
		term, err := parseTermFile(f, "")
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}
	for _, f := range scoped {
		term, err := parseTermFile(f, filepath.Base(filepath.Dir(f)))
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}

	return terms, nil
}
