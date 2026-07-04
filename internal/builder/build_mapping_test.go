package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

func TestRenderMappingRuleCarriesIDAnchorAndBadge(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{{ID: "R-01", Name: "Activities and steps can be overviewed"}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if !strings.Contains(html, `id="rule-R-01"`) {
		t.Fatal("expected rule sticky to carry an id anchor")
	}
	if !strings.Contains(html, `href="#rule-R-01" data-copy-link`) {
		t.Fatal("expected a one-click copy-link trigger on the rule sticky")
	}
	if !strings.Contains(html, "🔗") {
		t.Fatal("expected a link emoji to signal the badge is copyable")
	}
	if !strings.Contains(html, "R-01</a>") {
		t.Fatal("expected the rule ID to be displayed as a badge")
	}
}

func TestRenderMappingMarksAutomatedRules(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "An automated rule", Automated: true},
			{ID: "R-02", Name: "A rule not yet automated"},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if got := strings.Count(html, "✓ automated"); got != 1 {
		t.Fatalf("automated badge rendered %d times, want once (only on the automated rule)", got)
	}
	if !strings.Contains(html, "Automated rule") {
		t.Fatal("expected the legend to explain the automated mark")
	}
}

func TestRenderMappingRuleWithoutIDOmitsAnchor(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{{Name: "A rule without an ID"}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if strings.Contains(html, `href="#rule-`) {
		t.Fatal("expected no copy-link anchor when a rule has no ID")
	}
	if strings.Contains(html, `id="rule-"`) {
		t.Fatal("expected no empty rule anchor when a rule has no ID")
	}
}
