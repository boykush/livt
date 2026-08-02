package parser

import (
	"fmt"
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/uri"
	"gopkg.in/yaml.v3"
)

type personaFrontmatter struct {
	Name string `yaml:"name"`
}

// ParsePersona reads the persona a key names, relative to the personas
// directory. The key is guarded so an externally-supplied one — a story's
// frontmatter, a livt URI — cannot walk out of the directory.
func ParsePersona(personasDir, key string) (*domain.Persona, error) {
	if !uri.ValidSegment(key) {
		return nil, fmt.Errorf("invalid persona key %q", key)
	}
	return parsePersonaFile(filepath.Join(personasDir, key+".md"))
}

func parsePersonaFile(path string) (*domain.Persona, error) {
	key, frontmatter, body, err := splitFrontmatter(path)
	if err != nil {
		return nil, err
	}

	var fm personaFrontmatter
	if frontmatter != "" {
		if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
			return nil, err
		}
	}

	return &domain.Persona{
		Key:  key,
		Name: fm.Name,
		Body: body,
	}, nil
}

// ParseAllPersonas reads every persona in the personas directory. The glob is
// one level deep and sorts, so the order a build renders is the order the
// filesystem lays out.
func ParseAllPersonas(personasDir string) ([]*domain.Persona, error) {
	files, err := filepath.Glob(filepath.Join(personasDir, "*.md"))
	if err != nil {
		return nil, err
	}

	var personas []*domain.Persona
	for _, f := range files {
		persona, err := parsePersonaFile(f)
		if err != nil {
			return nil, err
		}
		personas = append(personas, persona)
	}

	return personas, nil
}
