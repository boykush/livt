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
