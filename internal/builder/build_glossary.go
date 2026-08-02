package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
)

// buildGlossary renders every ubiquitous language term as a row in the glossary
// table at ubiquitous.html. Each row carries an id={term-ref} anchor so it can
// be linked as ubiquitous.html#{term-ref} (e.g. from a story map or example
// mapping), where the reference holds the term's context when it has one.
func (b *Builder) buildGlossary() error {
	terms, err := parser.ParseAllTerms(b.UbiquitousDir)
	if err != nil {
		return err
	}

	var cards []glossaryCard
	var contexts []string
	seen := make(map[string]bool)
	for _, t := range terms {
		cards = append(cards, glossaryCard{
			Anchor:     uri.TermAnchor(t.Ctx, t.Key),
			Ctx:        t.Ctx,
			Key:        t.Key,
			Name:       t.Name,
			Definition: t.Body,
		})
		// Context-free terms offer no axis to filter on: they belong to every
		// context, so no chip would narrow the table to them.
		if t.Ctx != "" && !seen[t.Ctx] {
			seen[t.Ctx] = true
			contexts = append(contexts, t.Ctx)
		}
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

	if err := renderGlossary(f, glossaryView{
		Sidebar:  sb,
		Terms:    cards,
		Contexts: contexts,
	}); err != nil {
		return err
	}
	fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	return nil
}

// resolveTermCards turns referenced term references (from a story map or example
// mapping) into pink sticky cards. A reference with a matching term file links
// to its glossary row; an unresolved one renders as a plain card showing the
// reference, mirroring how keyless story cards degrade.
func (b *Builder) resolveTermCards(refs []string) []termCard {
	var cards []termCard
	for _, ref := range refs {
		cards = append(cards, b.resolveTermCard(ref))
	}
	return cards
}

// resolveTermCard resolves one reference — "{ctx}/{term-key}" or "{term-key}".
// The context stays on the card even when the term resolves: two contexts can
// hold the same word for different things, so a sticky showing only the display
// name would not say which of them is meant.
func (b *Builder) resolveTermCard(ref string) termCard {
	ctx, key, ok := uri.SplitTermRef(ref)
	if !ok {
		return termCard{Name: ref}
	}
	path := filepath.Join(b.UbiquitousDir, ctx, key+".md")
	if _, err := os.Stat(path); err != nil {
		return termCard{Name: key, Ctx: ctx}
	}
	card := termCard{Name: key, Ctx: ctx, Href: "../" + uri.TermPage(ctx, key)}
	term, err := parser.ParseTerm(b.UbiquitousDir, ref)
	if err == nil && term.Name != "" {
		card.Name = term.Name
	}
	return card
}
