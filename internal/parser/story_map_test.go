package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseStoryMapReadsReferencedTerms(t *testing.T) {
	path := filepath.Join(t.TempDir(), "map.yaml")
	data := []byte("name: discovery\nubiquitous:\n  - story-map\n  - story\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	storyMap, err := ParseStoryMap(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(storyMap.Ubiquitous) != 2 || storyMap.Ubiquitous[0] != "story-map" || storyMap.Ubiquitous[1] != "story" {
		t.Fatalf("got ubiquitous %v, want [story-map story]", storyMap.Ubiquitous)
	}
}

func TestParseStoryMapAllowsStoryCardsWithoutKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "map.yaml")
	data := []byte(`name: discovery
releases:
  - id: first-release
    name: First release
activities:
  - key: activity
    name: Activity
    steps:
      - key: step
        name: Step
        stories:
          - name: Lightweight card
            release: first-release
          - key: detailed-card
            name: Detailed card
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

	if storyMap.Releases[0].ID != "first-release" {
		t.Fatalf("got release id %q", storyMap.Releases[0].ID)
	}
	if stories[0].Release != "first-release" {
		t.Fatalf("got keyless story release %q", stories[0].Release)
	}
}

func TestParseStoryMapRejectsDuplicateReleaseIDs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "map.yaml")
	data := []byte(`name: discovery
releases:
  - id: first-release
    name: First release
  - id: first-release
    name: Duplicate release
`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := ParseStoryMap(path); err == nil {
		t.Fatal("expected duplicate release id error")
	}
}

func TestParseStoryMapRejectsUnknownStoryRelease(t *testing.T) {
	path := filepath.Join(t.TempDir(), "map.yaml")
	data := []byte(`name: discovery
releases:
  - id: first-release
    name: First release
activities:
  - key: activity
    name: Activity
    steps:
      - key: step
        name: Step
        stories:
          - name: Lightweight card
            release: missing-release
`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := ParseStoryMap(path); err == nil {
		t.Fatal("expected unknown release error")
	}
}

func TestParseStoryMapRejectsDuplicateStoryCardNames(t *testing.T) {
	path := filepath.Join(t.TempDir(), "map.yaml")
	data := []byte(`name: discovery
activities:
  - key: activity
    name: Activity
    steps:
      - key: first-step
        name: First step
        stories:
          - name: Duplicate card
      - key: second-step
        name: Second step
        stories:
          - name: Duplicate card
`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := ParseStoryMap(path); err == nil {
		t.Fatal("expected duplicate story card name error")
	}
}
