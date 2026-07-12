package builder

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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

func TestStoryMapRendersActivitiesInARowWithStepsBeneath(t *testing.T) {
	b := Builder{}
	view := b.toStoryMapView(&domain.StoryMap{
		Name: "discovery",
		Activities: []domain.Activity{
			{
				Key:  "browse",
				Name: "Browse catalog",
				Steps: []domain.Step{
					{Key: "search", Name: "Search products"},
					{Key: "view", Name: "View product"},
				},
			},
			{
				Key:  "checkout",
				Name: "Check out",
				Steps: []domain.Step{
					{Key: "pay", Name: "Pay for order"},
				},
			},
		},
	})

	var buf bytes.Buffer
	if err := renderStoryMap(&buf, view); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	const activitiesRow = `<div class="flex gap-8 items-start">`
	if got := strings.Count(html, activitiesRow); got != 1 {
		t.Fatalf("activities row container rendered %d times, want a single shared row", got)
	}
	rowStart := strings.Index(html, activitiesRow)
	rowEnd := matchingDivEnd(t, html, rowStart)

	browse := indexOf(t, html, "Browse catalog")
	checkout := indexOf(t, html, "Check out")
	if browse < rowStart || rowEnd < browse || checkout < rowStart || rowEnd < checkout {
		t.Fatal("expected every activity card inside the single row container")
	}
	if checkout < browse {
		t.Fatal("expected activity cards to keep declaration order")
	}

	search := indexOf(t, html, "Search products")
	viewProduct := indexOf(t, html, "View product")
	pay := indexOf(t, html, "Pay for order")
	if browse >= search || search >= viewProduct || viewProduct >= checkout {
		t.Fatal("expected the first activity's steps beneath its card, before the next activity")
	}
	if checkout >= pay || pay >= rowEnd {
		t.Fatal("expected the second activity's steps beneath its card, inside the row")
	}
}

func TestRenderStoryMapKeyedCardCarriesIDAnchorAndCopyLink(t *testing.T) {
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
							{Key: domain.StoryKey{Value: "detailed-card"}, Name: "Detailed card", Release: "first-release"},
							{Key: domain.StoryKey{Value: "planned-card"}, Name: "Planned card"},
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

	var buf bytes.Buffer
	if err := renderStoryMap(&buf, view); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	// detailed-card sits in a release row as a story-page link, planned-card in
	// the unscoped row as a plain card; both keyed shapes carry the same anchor
	// and one-click copy-link trigger.
	for _, key := range []string{"detailed-card", "planned-card"} {
		if !strings.Contains(html, `id="story-`+key+`"`) {
			t.Fatalf("expected story card %s to carry an id anchor", key)
		}
		if !strings.Contains(html, `href="#story-`+key+`" data-copy-link`) {
			t.Fatalf("expected a one-click copy-link trigger on story card %s", key)
		}
	}
	if !strings.Contains(html, "🔗") {
		t.Fatal("expected a link emoji to signal the trigger is copyable")
	}
	if !strings.Contains(html, `class="story-card relative scroll-mt-24"`) {
		t.Fatal("expected anchored story cards to keep a scroll margin so deep links land on the card")
	}
}

func TestRenderStoryMapKeylessCardHasNoCopyLink(t *testing.T) {
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

	var buf bytes.Buffer
	if err := renderStoryMap(&buf, view); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if strings.Contains(html, `data-copy-link title=`) {
		t.Fatal("expected no copy-link trigger on a keyless story candidate")
	}
	if strings.Contains(html, `id="story-`) {
		t.Fatal("expected no id anchor on a keyless story candidate")
	}
}

func indexOf(t *testing.T, html, substr string) int {
	t.Helper()
	i := strings.Index(html, substr)
	if i == -1 {
		t.Fatalf("expected rendered story map to contain %q", substr)
	}
	return i
}

// matchingDivEnd returns the offset just past the </div> that closes the div
// opening at start, so callers can assert containment by position.
func matchingDivEnd(t *testing.T, html string, start int) int {
	t.Helper()
	depth := 0
	for i := start; i < len(html); {
		open := strings.Index(html[i:], "<div")
		close := strings.Index(html[i:], "</div>")
		if close == -1 {
			break
		}
		if open != -1 && open < close {
			depth++
			i += open + len("<div")
			continue
		}
		depth--
		i += close + len("</div>")
		if depth == 0 {
			return i
		}
	}
	t.Fatal("no matching </div> for the row container")
	return -1
}
