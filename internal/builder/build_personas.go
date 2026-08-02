package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
)

// buildPersonas renders every persona as a row at personas.html, each carrying
// an id={persona-key} anchor so a story can link back to the actor it is written
// for. storyIndex is the reverse link — the stories that named this persona —
// which is what turns the list from a second glossary into the answer to "whose
// stories are these".
func (b *Builder) buildPersonas(storyIndex map[string][]personaStoryRef) error {
	personas, err := parser.ParseAllPersonas(b.PersonasDir)
	if err != nil {
		return err
	}

	var rows []personaRow
	for _, p := range personas {
		rows = append(rows, personaRow{
			Anchor:      uri.PersonaAnchor(p.Key),
			Key:         p.Key,
			Name:        p.Name,
			Description: p.Body,
			Stories:     storyIndex[p.Key],
		})
	}

	sb, err := b.sidebar("persona", "")
	if err != nil {
		return err
	}

	outPath := filepath.Join(b.OutDir, "personas.html")
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := renderPersonas(f, personasView{Sidebar: sb, Personas: rows}); err != nil {
		return err
	}
	fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	return nil
}

// resolvePersonaCard turns a story's persona key into the chip its pages show.
// An empty key yields no chip at all: a story may name no actor, which is not
// the same as naming one that is not committed yet.
func (b *Builder) resolvePersonaCard(key, prefix string) *personaCard {
	if key == "" {
		return nil
	}
	card := b.personaCardFor(key, prefix)
	return &card
}

// resolvePersonaCards turns the keys a story map declares into stickies for its
// board, the way resolveTermCards does for referenced terms.
func (b *Builder) resolvePersonaCards(keys []string, prefix string) []personaCard {
	var cards []personaCard
	for _, key := range keys {
		cards = append(cards, b.personaCardFor(key, prefix))
	}
	return cards
}

// personaCardFor resolves one key. A key with a matching file links to its row;
// an unresolved one renders as a plain card showing the key, mirroring how term
// cards degrade. prefix is the relative path back to the output root from the
// page doing the rendering.
func (b *Builder) personaCardFor(key, prefix string) personaCard {
	if !uri.ValidSegment(key) {
		return personaCard{Name: key}
	}
	card := personaCard{Name: key}
	persona, err := parser.ParsePersona(b.PersonasDir, key)
	if err != nil {
		return card
	}
	if persona.Name != "" {
		card.Name = persona.Name
	}
	card.Href = prefix + uri.PersonaPage(key)
	return card
}
