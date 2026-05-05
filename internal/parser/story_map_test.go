package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseStoryMapAllowsStoryCardsWithoutKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "map.yaml")
	data := []byte(`name: discovery
activities:
  - key: activity
    name: Activity
    steps:
      - key: step
        name: Step
        stories:
          - name: Lightweight card
          - key: detailed-card
            name: Detailed card
releases:
  - name: First release
    stories:
      - detailed-card
`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	storyMap, err := ParseStoryMap(path)
	if err != nil {
		t.Fatal(err)
	}

	stories := storyMap.Activities[0].Steps[0].Stories
	if len(stories) != 2 {
		t.Fatalf("got %d stories, want 2", len(stories))
	}
	if stories[0].HasKey() {
		t.Fatalf("expected first story to be keyless, got key %q", stories[0].Key.Value)
	}
	if stories[0].Name != "Lightweight card" {
		t.Fatalf("got keyless story name %q", stories[0].Name)
	}

	releaseStories := storyMap.Releases[0].Stories
	if len(releaseStories) != 1 {
		t.Fatalf("got %d release stories, want 1", len(releaseStories))
	}
	if releaseStories[0].Name != "Detailed card" {
		t.Fatalf("got release story name %q", releaseStories[0].Name)
	}
}
