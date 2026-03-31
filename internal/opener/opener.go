package opener

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
)

type Opener struct {
	StoriesDir string
	USMDir     string
}

func (o *Opener) Open(key string) error {
	storyKey := domain.StoryKey{Value: key}

	card, err := o.findStoryCard(storyKey)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(o.StoriesDir, 0755); err != nil {
		return err
	}

	path := filepath.Join(o.StoriesDir, card.Key.Value+".md")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("story %q is already opened (%s)", key, path)
	}

	content := fmt.Sprintf("---\nname: %s\n---\n\nAs a\nI want to\nSo that\n", card.Name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Printf("Opened %s\n", path)
	return nil
}

func (o *Opener) findStoryCard(key domain.StoryKey) (*domain.StoryCard, error) {
	maps, err := parser.ParseAllStoryMaps(o.USMDir)
	if err != nil {
		return nil, err
	}

	for _, sm := range maps {
		for _, a := range sm.Activities {
			for _, s := range a.Steps {
				for i, sc := range s.Stories {
					if sc.Key.Value == key.Value {
						return &s.Stories[i], nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("story card %q not found in any USM YAML", key.Value)
}
