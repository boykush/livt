package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

func TestStoryMapViewRendersKeylessStoriesAsUnscopedPlainCards(t *testing.T) {
	storiesDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(storiesDir, "detailed-card.md"), []byte("---\nname: Detailed card\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Builder{StoriesDir: storiesDir}
	view := b.toStoryMapView(&domain.StoryMap{
		Name: "discovery",
		Activities: []domain.Activity{
			{
				Key:  "activity",
				Name: "Activity",
				Steps: []domain.Step{
					{
						Key:  "step",
						Name: "Step",
						Stories: []domain.StoryCard{
							{Name: "Lightweight card"},
							{Key: domain.StoryKey{Value: "detailed-card"}, Name: "Detailed card"},
						},
					},
				},
			},
		},
		Releases: []domain.Release{
			{
				Name: "First release",
				Stories: []domain.StoryCard{
					{Key: domain.StoryKey{Value: "detailed-card"}, Name: "Detailed card"},
				},
			},
		},
	})

	releaseStories := view.StoryMap.ReleaseRows[0].Activities[0].StepStories[0]
	if len(releaseStories) != 1 {
		t.Fatalf("got %d release stories, want 1", len(releaseStories))
	}
	if !releaseStories[0].Opened {
		t.Fatal("expected keyed story with markdown file to be opened")
	}

	if view.StoryMap.UnscopedStories == nil {
		t.Fatal("expected keyless story to appear in unscoped stories")
	}
	unscopedStories := view.StoryMap.UnscopedStories.Activities[0].StepStories[0]
	if len(unscopedStories) != 1 {
		t.Fatalf("got %d unscoped stories, want 1", len(unscopedStories))
	}
	if unscopedStories[0].Key != "" {
		t.Fatalf("got keyless story key %q", unscopedStories[0].Key)
	}
	if unscopedStories[0].Opened {
		t.Fatal("expected keyless story to render as a plain card")
	}
}
