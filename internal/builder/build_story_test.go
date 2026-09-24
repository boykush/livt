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
		{Name: "協働ディスカバリー", Path: "../story-map/協働ディスカバリー.html"},
		{Name: "ディスカバリーと開発のギャップ", Path: "../story-map/ディスカバリーと開発のギャップ.html"},
	}

	var buf bytes.Buffer
	if err := renderStory(&buf, i18n.En, story, "", opportunities, nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	for _, o := range opportunities {
		if !strings.Contains(html, `href="`+mapHref("../story-map/", o.Name)+`"`) {
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

// Each opportunity link says what it leads to beside the name R-02 keeps on it,
// and wears the opportunity's colour: purple is the story map's, and would read
// as a way to the board rather than to the opportunity.
func TestRenderStoryOpportunityLinkSaysItsKind(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "s"}, Name: "S"}
	opportunities := []opportunityRef{{Name: "デモ", Path: "../opportunity/demo.html"}}

	var buf bytes.Buffer
	if err := renderStory(&buf, i18n.En, story, "", opportunities, nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	at := strings.Index(html, `href="../opportunity/demo.html"`)
	if at < 0 {
		t.Fatal("expected a Related link to the opportunity")
	}
	if tag := html[at : at+strings.Index(html[at:], ">")]; strings.Contains(tag, "purple") || !strings.Contains(tag, "orange") {
		t.Errorf("the opportunity link wears %q, want the opportunity's orange", tag)
	}
	text, _ := linkText(html, "../opportunity/demo.html")
	if want := i18n.Of(i18n.En).Msg("label.opportunity"); !strings.Contains(text, want) || !strings.Contains(text, "デモ") {
		t.Errorf("the opportunity link reads %q, want %q beside the name", text, want)
	}
}

// A story on no map shows no opportunity link in the Related section.
func TestRenderStoryWithoutOpportunitiesShowsNoMapLink(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "orphan"}, Name: "Orphan"}

	var buf bytes.Buffer
	if err := renderStory(&buf, i18n.En, story, "", nil, nil); err != nil {
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
	if err := renderStory(&buf, i18n.En, story, "", nil, nil); err != nil {
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
