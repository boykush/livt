package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/i18n"
)

// livt:automates livt://mapping/filter-lists-by-opportunity/rule/R-02
// The Related section links to each opportunity by name (one link per map),
// replacing the single generic "Story Map" link. Paths are relative to the
// story/ directory.
func TestRenderStoryLinksEachOpportunityByName(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "record-rule-automation"}, Name: "Record rule automation"}
	opportunities := []opportunityRef{
		{Name: "協働ディスカバリー", Path: "../story-map/collaborative-discovery.html"},
		{Name: "ディスカバリーと開発のギャップ", Path: "../story-map/discovery-development-gap.html"},
	}

	var buf bytes.Buffer
	if err := renderStory(&buf, i18n.En, story, "", opportunities, nil, nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	for _, o := range opportunities {
		if !strings.Contains(html, `href="`+o.Path+`"`) {
			t.Fatalf("expected a Related link to the %q board", o.Name)
		}
		if !strings.Contains(html, o.Name) {
			t.Fatalf("expected the opportunity name %q on the story page", o.Name)
		}
	}
	if strings.Contains(html, "Story Map &rarr;") {
		t.Fatal(`expected the generic "Story Map" link to be replaced by named opportunity links`)
	}
}

// An opportunity link says what it leads to, in the opportunity's colour: purple
// is the story map's, and would read as a way to the board. The kind alone does
// for a story on one map; on several, each link is named too, since the kind
// would not tell them apart.
func TestRenderStoryOpportunityLinkSaysItsKind(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "s"}, Name: "S"}
	label := i18n.Of(i18n.En).Msg("label.opportunity")
	demo := opportunityRef{Name: "デモ", Path: "../opportunity/demo.html"}
	other := opportunityRef{Name: "別件", Path: "../opportunity/other.html"}

	for _, opportunities := range [][]opportunityRef{{demo}, {demo, other}} {
		var buf bytes.Buffer
		if err := renderStory(&buf, i18n.En, story, "", opportunities, nil, nil); err != nil {
			t.Fatal(err)
		}
		html := buf.String()
		named := len(opportunities) > 1

		for _, o := range opportunities {
			at := strings.Index(html, `href="`+o.Path+`"`)
			if at < 0 {
				t.Fatalf("expected a Related link to %s", o.Path)
			}
			if tag := html[at : at+strings.Index(html[at:], ">")]; strings.Contains(tag, "purple") || !strings.Contains(tag, "orange") {
				t.Errorf("the link to %s wears %q, want the opportunity's orange", o.Path, tag)
			}
			text, _ := linkText(html, o.Path)
			if !strings.Contains(text, label) || strings.Contains(text, o.Name) != named {
				t.Errorf("with %d opportunities the link to %s reads %q, want %q, named: %v", len(opportunities), o.Path, text, label, named)
			}
		}
	}
}

// A story on no map shows no opportunity link in the Related section.
func TestRenderStoryWithoutOpportunitiesShowsNoMapLink(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "orphan"}, Name: "Orphan"}

	var buf bytes.Buffer
	if err := renderStory(&buf, i18n.En, story, "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "../story-map/") {
		t.Fatal("expected no opportunity link for a story on no map")
	}
}

// The story page shows the story's name as its title and its heading, so a
// story with no name: frontmatter would render both blank. The key stands in.
func TestRenderStoryWithoutNameShowsTheKey(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "unnamed-story"}}

	var buf bytes.Buffer
	if err := renderStory(&buf, i18n.En, story, "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if strings.Contains(html, "<title> - livt</title>") {
		t.Error("expected the page title to name the story rather than be blank")
	}
	if !strings.Contains(html, ">unnamed-story</h1>") {
		t.Error("expected the heading to fall back to the story key")
	}
}

// livt:automates livt://mapping/preview-story-map-in-browser/rule/R-06
// The story page leads back to its card in the story map's purple. The kind
// alone does for a story on one map; on several, each link names its map, and a
// story on none has no link at all.
func TestRenderStoryLinksItsCardOnEachStoryMap(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "s"}, Name: "S"}
	label := i18n.Of(i18n.En).Msg("label.story-map")
	demo := storyMapRef{Name: "デモ", Path: "../story-map/demo.html#story-s"}
	other := storyMapRef{Name: "別件", Path: "../story-map/other.html#story-s"}

	for _, cards := range [][]storyMapRef{nil, {demo}, {demo, other}} {
		var buf bytes.Buffer
		if err := renderStory(&buf, i18n.En, story, "", nil, cards, nil); err != nil {
			t.Fatal(err)
		}
		html := buf.String()
		named := len(cards) > 1

		if len(cards) == 0 && strings.Contains(html, "../story-map/") {
			t.Error("expected no story map link for a story on no map")
		}
		for _, c := range cards {
			at := strings.Index(html, `href="`+c.Path+`"`)
			if at < 0 {
				t.Fatalf("expected a Related link to %s", c.Path)
			}
			if tag := html[at : at+strings.Index(html[at:], ">")]; !strings.Contains(tag, "purple") {
				t.Errorf("the link to %s wears %q, want the story map's purple", c.Path, tag)
			}
			text, _ := linkText(html, c.Path)
			if !strings.Contains(text, label) || strings.Contains(text, c.Name) != named {
				t.Errorf("with %d maps the link to %s reads %q, want %q, named: %v", len(cards), c.Path, text, label, named)
			}
		}
	}
}
