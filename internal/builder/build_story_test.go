package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

// R-02 on the story page: the Related section links to each opportunity by name
// (one link per map), replacing the single generic "Story Map" link. Paths are
// relative to the story/ directory.
func TestRenderStoryLinksEachOpportunityByName(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "record-rule-automation"}, Name: "Record rule automation"}
	opportunities := []opportunityRef{
		{Name: "協働ディスカバリー", Path: "../story-map/協働ディスカバリー.html"},
		{Name: "ディスカバリーと開発のギャップ", Path: "../story-map/ディスカバリーと開発のギャップ.html"},
	}

	var buf bytes.Buffer
	if err := renderStory(&buf, story, "", opportunities); err != nil {
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

// A story on no map shows no opportunity link in the Related section.
func TestRenderStoryWithoutOpportunitiesShowsNoMapLink(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "orphan"}, Name: "Orphan"}

	var buf bytes.Buffer
	if err := renderStory(&buf, story, "", nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "../story-map/") {
		t.Fatal("expected no opportunity link for a story on no map")
	}
}
