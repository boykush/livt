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

// resolveTermCards turns referenced term keys (from a story map or example
// mapping) into pink sticky cards. A key with a matching ubiquitous/{key}.md
// file links to its glossary row; an unresolved key renders as a plain card
// showing the key, mirroring how keyless story cards degrade.
func (b *Builder) resolveTermCards(keys []string) []termCard {
	var cards []termCard
	for _, key := range keys {
		cards = append(cards, b.resolveTermCard(key))
	}
	return cards
}

func (b *Builder) resolveTermCard(key string) termCard {
	path := filepath.Join(b.UbiquitousDir, key+".md")
	if _, err := os.Stat(path); err != nil {
		return termCard{Name: key}
	}
	term, err := parser.ParseTerm(path)
	if err != nil || term.Name == "" {
		return termCard{Name: key, Href: "../ubiquitous.html#" + key}
	}
	return termCard{Name: term.Name, Href: "../ubiquitous.html#" + key}
}
