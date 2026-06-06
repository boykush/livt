package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// emptyDirsBuilder returns a Builder whose resource directories are all empty so
// sidebar counts are deterministic and only the passed-in tiles drive output.
func emptyDirsBuilder(t *testing.T) Builder {
	t.Helper()
	return Builder{
		MappingsDir:   t.TempDir(),
		StoriesDir:    t.TempDir(),
		USMDir:        t.TempDir(),
		UbiquitousDir: t.TempDir(),
		OutDir:        t.TempDir(),
	}
}

func TestBuildMappingsIndexRendersPreviewCards(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.buildMappingsIndex([]mappingTile{{Key: "checkout", StoryName: "Checkout flow"}}); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(filepath.Join(b.OutDir, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)

	for _, want := range []string{
		"Example Mappings",            // section / sidebar label
		"Checkout flow",               // tile title
		`src="mapping/checkout.html"`, // iframe preview source
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("index.html missing %q", want)
		}
	}
}

func TestBuildMappingsIndexEmptyState(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.buildMappingsIndex(nil); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(filepath.Join(b.OutDir, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "No example mappings yet.") {
		t.Fatal("expected English empty-state message when no mappings exist")
	}
}

func TestBuildStoryMapsIndexRendersPreviewCards(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.buildStoryMapsIndex([]storyMapTile{{Name: "discovery"}}); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(filepath.Join(b.OutDir, "story-maps.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, `src="story-map/discovery.html"`) {
		t.Fatal("expected story map preview iframe source")
	}
	if !strings.Contains(html, "Story Maps") {
		t.Fatal("expected Story Maps label")
	}
}

func TestBuildStoriesIndexRendersList(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.buildStoriesIndex([]storyItem{{Key: "first-story", Name: "First story"}}); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(filepath.Join(b.OutDir, "stories.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "First story") {
		t.Fatal("expected story name in list")
	}
	if !strings.Contains(html, `href="story/first-story.html"`) {
		t.Fatal("expected link to story page")
	}
}
