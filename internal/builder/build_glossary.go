package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/parser"
)

// buildGlossary renders every ubiquitous language term as a row in the glossary
// table at ubiquitous.html. Each row carries an id={key} anchor so it can be
// linked as ubiquitous.html#{key} (e.g. from a story map or example mapping).
func (b *Builder) buildGlossary() error {
	terms, err := parser.ParseAllTerms(b.UbiquitousDir)
	if err != nil {
		return err
	}

	var cards []glossaryCard
	for _, t := range terms {
		cards = append(cards, glossaryCard{
			Key:        t.Key,
			Name:       t.Name,
			Definition: t.Body,
		})
	}

	sb, err := b.sidebar("ubiquitous", "")
	if err != nil {
		return err
	}

	outPath := filepath.Join(b.OutDir, "ubiquitous.html")
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := renderGlossary(f, glossaryView{Sidebar: sb, Terms: cards}); err != nil {
		return err
	}
	fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	return nil
}
