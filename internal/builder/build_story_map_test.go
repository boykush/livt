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
							{Key: domain.StoryKey{Value: "detailed-card"}, Name: "Detailed card", Release: "first-release"},
						},
					},
				},
			},
		},
		Releases: []domain.Release{
			{
				ID:   "first-release",
				Name: "First release",
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

func TestStoryMapViewUsesStoryReleaseForReleaseRows(t *testing.T) {
	b := Builder{}
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
							{Name: "Keyless release card", Release: "first-release"},
							{Key: domain.StoryKey{Value: "later-card"}, Name: "Later card"},
						},
					},
				},
			},
		},
		Releases: []domain.Release{
			{
				ID:   "first-release",
				Name: "First release",
			},
		},
	})

	releaseStories := view.StoryMap.ReleaseRows[0].Activities[0].StepStories[0]
	if len(releaseStories) != 1 {
		t.Fatalf("got %d release stories, want 1", len(releaseStories))
	}
	if releaseStories[0].Name != "Keyless release card" {
		t.Fatalf("got release story name %q", releaseStories[0].Name)
	}
	if releaseStories[0].Opened {
		t.Fatal("expected keyless release story to render as a plain card")
	}

	if view.StoryMap.UnscopedStories == nil {
		t.Fatal("expected unscoped story row")
	}
	unscopedStories := view.StoryMap.UnscopedStories.Activities[0].StepStories[0]
	if len(unscopedStories) != 1 {
		t.Fatalf("got %d unscoped stories, want 1", len(unscopedStories))
	}
	if unscopedStories[0].Name != "Later card" {
		t.Fatalf("got unscoped story name %q", unscopedStories[0].Name)
	}
}

func TestStoryMapViewWithoutReleasesHasNoReleaseRows(t *testing.T) {
	b := Builder{}
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
						},
					},
				},
			},
		},
	})

	if len(view.StoryMap.ReleaseRows) != 0 {
		t.Fatalf("got %d release rows, want 0", len(view.StoryMap.ReleaseRows))
	}
	if view.StoryMap.UnscopedStories == nil {
		t.Fatal("expected stories to render without release dividers")
	}
}
