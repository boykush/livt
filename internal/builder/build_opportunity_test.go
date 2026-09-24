package builder

import (
	gohtml "html"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/i18n"
)

const demoCanvasYAML = "canvas:\n" +
	"  solution-ideas:\n    - テキストで残す\n" +
	"  problems:\n    - 共通理解が失われる\n" +
	"  budget:\n    - 余暇で開発\n"

// The canvas is one sheet in three zones — the facts, the solution, the value —
// so every box has to land in the panel the method puts it in.
func TestRenderOpportunityCanvasLaysOutTheThreeZones(t *testing.T) {
	b := outDirsBuilder(t)
	writeFile(t, filepath.Join(b.OpportunitiesDir, "demo.md"), "---\nname: デモ\n---\n\n一文\n")
	writeFile(t, filepath.Join(b.CanvasesDir, "demo.yaml"), demoCanvasYAML)

	if err := b.buildOpportunityCanvas(&domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ"}); err != nil {
		t.Fatal(err)
	}
	html := readFile(t, filepath.Join(b.OutDir, "opportunity-canvas", "demo.html"))

	// All ten boxes are on the sheet, unanswered ones included: a blank box is
	// the record of a question the opportunity has not answered.
	// Headings are compared escaped: one of them is Patton's "Customers &
	// Users", and the template is right to render the ampersand as &amp;.
	en := i18n.Of(i18n.En)
	for _, box := range (&domain.OpportunityCanvas{}).Boxes() {
		name := en.Msg("canvas." + box.Key)
		if !strings.Contains(html, gohtml.EscapeString(name)) {
			t.Errorf("box %q is missing from the canvas", name)
		}
	}
	if got := strings.Count(html, en.Msg("canvas.empty-box")); got != 7 {
		t.Errorf("got %d unanswered boxes, want 7 (three of the ten are filled in)", got)
	}
	for _, panel := range []string{"bg-rose-50", "bg-sky-50", "bg-emerald-50"} {
		if !strings.Contains(html, panel) {
			t.Errorf("the %s zone panel is missing; the sheet is not laid out in three zones", panel)
		}
	}
	// Facts sit left of the solution, which sits left of the value: the sheet
	// reads back from the idea to the problem, then forward to the value.
	facts, solution, value := strings.Index(html, "bg-rose-50"), strings.Index(html, "bg-sky-50"), strings.Index(html, "bg-emerald-50")
	if facts >= solution || solution >= value {
		t.Errorf("zones are out of order: facts=%d solution=%d value=%d", facts, solution, value)
	}
}

// livt:automates livt://opportunity/collaborative-discovery
// The Related section is what says how far an opportunity has been taken. A
// canvas means it was thought through, a story map means it was taken on.
func TestRenderOpportunityLinksItsCanvasAndStoryMap(t *testing.T) {
	b := outDirsBuilder(t)
	o := &domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ", Body: "一文で言う"}
	maps := []storyMapRef{{Name: "デモマップ", Path: "../story-map/デモマップ.html"}}

	out := filepath.Join(b.OutDir, "opportunity", "demo.html")
	if err := b.renderOpportunityPage(out, o, "../opportunity-canvas/demo.html", "", maps, opportunityProgress{}); err != nil {
		t.Fatal(err)
	}
	html := readFile(t, out)

	for _, want := range []string{"Opportunity Canvas", "デモマップ", "一文で言う"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %q on the opportunity page", want)
		}
	}
}

// An opportunity nobody has held a canvas session for links to none, so the
// absence stays visible instead of resolving to an empty board.
func TestRenderOpportunityWithoutCanvasLinksNone(t *testing.T) {
	b := outDirsBuilder(t)
	o := &domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ"}

	out := filepath.Join(b.OutDir, "opportunity", "demo.html")
	if err := b.renderOpportunityPage(out, o, "", "", nil, opportunityProgress{}); err != nil {
		t.Fatal(err)
	}
	if html := readFile(t, out); strings.Contains(html, "opportunity-canvas/") {
		t.Fatal("expected no canvas link for an opportunity with no canvas")
	}
}

// The filename is the join: a map whose key names a committed opportunity earns
// its stories a chip pointing at that opportunity's page, not at the map.
func TestStoryChipsPointAtTheOpportunityTheMapServes(t *testing.T) {
	b := outDirsBuilder(t)
	writeFile(t, filepath.Join(b.OpportunitiesDir, "demo.md"), "---\nname: デモ機会\n---\n\n一文\n")
	writeFile(t, filepath.Join(b.USMDir, "demo.yaml"),
		"name: デモマップ\nactivities:\n  - key: a\n    name: A\n    steps:\n      - key: s\n        name: S\n        stories:\n          - key: card\n            name: カード\n")

	opportunities, err := b.opportunityIndex()
	if err != nil {
		t.Fatal(err)
	}
	built, err := b.buildStoryMaps(opportunities)
	if err != nil {
		t.Fatal(err)
	}

	refs := built.StoryOpportunities["card"]
	if len(refs) != 1 {
		t.Fatalf("got %d opportunity refs, want 1", len(refs))
	}
	if refs[0].Name != "デモ機会" {
		t.Errorf("chip reads %q, want the opportunity's name rather than the map's", refs[0].Name)
	}
	if refs[0].Path != "../opportunity/demo.html" {
		t.Errorf("chip links to %q, want the opportunity's page", refs[0].Path)
	}
	if got := built.MapsByOpportunity["demo"]; len(got) != 1 || got[0].Name != "デモマップ" {
		t.Errorf("MapsByOpportunity[demo] = %+v, want the map mapped for it", got)
	}
}

// A livt repository with no opportunities/ at all keeps working exactly as it
// did before opportunities became files: the map stands in as its own
// opportunity, named by the map. Nothing has to be migrated to keep building.
func TestMapWithNoOpportunityFileStandsInAsItsOwn(t *testing.T) {
	b := outDirsBuilder(t)
	writeFile(t, filepath.Join(b.USMDir, "demo.yaml"),
		"name: デモマップ\nactivities:\n  - key: a\n    name: A\n    steps:\n      - key: s\n        name: S\n        stories:\n          - key: card\n            name: カード\n")

	opportunities, err := b.opportunityIndex()
	if err != nil {
		t.Fatal(err)
	}
	built, err := b.buildStoryMaps(opportunities)
	if err != nil {
		t.Fatal(err)
	}

	refs := built.StoryOpportunities["card"]
	if len(refs) != 1 || refs[0].Name != "デモマップ" || refs[0].Path != "../story-map/デモマップ.html" {
		t.Fatalf("got %+v, want the map standing in as its own opportunity", refs)
	}
	if len(built.MapsByOpportunity) != 0 {
		t.Errorf("MapsByOpportunity = %+v, want empty: no opportunity file claims this map", built.MapsByOpportunity)
	}
	// With nothing to link to, the board does not offer an opportunity link.
	if len(built.Tiles) != 1 || built.Tiles[0].Opportunity != nil {
		t.Errorf("tile = %+v, want no opportunity chip", built.Tiles)
	}
}

// A canvas whose opportunity was renamed or deleted must not outlive it, the
// same guarantee TestBuildDropsPagesForRemovedResources makes for the rest.
func TestBuildDropsPagesForRemovedOpportunities(t *testing.T) {
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.OpportunitiesDir, "kept.md"), "---\nname: 残る\n---\n\n一文\n")
	writeFile(t, filepath.Join(b.CanvasesDir, "kept.yaml"), demoCanvasYAML)
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	orphans := []string{
		filepath.Join(b.OutDir, "opportunity", "removed.html"),
		filepath.Join(b.OutDir, "opportunity-canvas", "removed.html"),
	}
	for _, p := range orphans {
		writeFile(t, p, "<html>stale</html>")
	}
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	for _, p := range orphans {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("%s survived the rebuild (stat error: %v)", p, err)
		}
	}
	for _, p := range []string{
		filepath.Join(b.OutDir, "opportunity", "kept.html"),
		filepath.Join(b.OutDir, "opportunity-canvas", "kept.html"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("rebuild dropped a live page: %v", err)
		}
	}
}

// An opportunity file need not carry a name, and every surface naming one is a
// heading, a filter chip, or a link label. It falls back to the key rather than
// rendering blank — the same call Story.DisplayName makes.
func TestNamelessOpportunityFallsBackToItsKey(t *testing.T) {
	b := outDirsBuilder(t)
	writeFile(t, filepath.Join(b.OpportunitiesDir, "demo.md"), "一文だけの本文\n")
	writeFile(t, filepath.Join(b.USMDir, "demo.yaml"),
		"name: デモマップ\nactivities:\n  - key: a\n    name: A\n    steps:\n      - key: s\n        name: S\n        stories:\n          - key: card\n            name: カード\n")

	opportunities, err := b.opportunityIndex()
	if err != nil {
		t.Fatal(err)
	}
	built, err := b.buildStoryMaps(opportunities)
	if err != nil {
		t.Fatal(err)
	}
	// A blank chip would be an unpickable filter axis, not just an ugly one.
	if refs := built.StoryOpportunities["card"]; len(refs) != 1 || refs[0].Name != "demo" {
		t.Fatalf("chip = %+v, want it to fall back to the key", refs)
	}

	if _, err := b.buildOpportunities(built.MapsByOpportunity, built.StoriesByOpportunity, nil); err != nil {
		t.Fatal(err)
	}
	if html := readFile(t, filepath.Join(b.OutDir, "opportunity", "demo.html")); !strings.Contains(html, ">demo</h1>") {
		t.Error("the opportunity page heading is blank")
	}
}

// outDirsBuilder is emptyDirsBuilder with the per-resource output directories
// laid out, for a test that drives one build step rather than a whole Build.
func outDirsBuilder(t *testing.T) Builder {
	t.Helper()
	b := emptyDirsBuilder(t)
	if err := b.resetGeneratedDirs(); err != nil {
		t.Fatal(err)
	}
	return b
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// progressBuilder lays out one opportunity whose map hangs three keyed stories,
// two of which have been through an example mapping. The mappings between them
// carry every standing the reading reports: automated, un-automated, proposed,
// retired, and an open question.
func progressBuilder(t *testing.T) Builder {
	t.Helper()
	b := outDirsBuilder(t)
	writeFile(t, filepath.Join(b.OpportunitiesDir, "demo.md"), "---\nname: デモ機会\n---\n\n一文\n")
	writeFile(t, filepath.Join(b.USMDir, "demo.yaml"),
		"name: デモマップ\n"+
			"activities:\n"+
			"  - key: a1\n    name: アクティビティ\n"+
			"    steps:\n"+
			"      - key: s1\n        name: ステップ\n"+
			"        stories:\n"+
			"          - name: 押さえたストーリー\n            key: held-story\n"+
			"          - name: 途中のストーリー\n            key: half-story\n"+
			"          - name: まだのストーリー\n            key: unmapped-story\n")
	for key, name := range map[string]string{
		"held-story": "押さえたストーリー", "half-story": "途中のストーリー", "unmapped-story": "まだのストーリー",
	} {
		writeFile(t, filepath.Join(b.StoriesDir, key+".md"), "---\nname: "+name+"\n---\n")
	}
	writeFile(t, filepath.Join(b.MappingsDir, "held-story.yaml"),
		"rules:\n"+
			"  - id: R-01\n    name: 押さえたルール\n"+
			"  - id: R-02\n    name: 退役したルール\n    status: retired\n"+
			"questions: []\n")
	writeAutomations(t, b, "livt://mapping/held-story/rule/R-01")
	writeFile(t, filepath.Join(b.MappingsDir, "half-story.yaml"),
		"rules:\n"+
			"  - id: R-01\n    name: まだのルール\n"+
			"  - id: R-02\n    name: 提案中のルール\n    status: proposed\n"+
			"questions:\n"+
			"  - id: Q-01\n    text: 開いている疑問\n"+
			"  - id: Q-02\n    text: 片付いた疑問\n    retired: true\n")
	return b
}

// livt:automates livt://mapping/show-opportunity-progress/rule/R-02
// The two gauges are the reading, and neither may count what the board has
// closed: a retired rule is not spec anyone is waiting on, so counting it would
// make the opportunity read as less finished than it is.
func TestOpportunityProgressCountsStoriesAndLiveRules(t *testing.T) {
	b := progressBuilder(t)
	_, _, tallies, err := b.buildMappings()
	if err != nil {
		t.Fatal(err)
	}

	o := &domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ機会"}
	p := b.progressOf(o, oneSlice("held-story", "half-story", "unmapped-story"), tallies)

	if p.TotalStories != 3 || p.MappedStories != 2 {
		t.Errorf("mapped stories = %d/%d, want 2/3", p.MappedStories, p.TotalStories)
	}
	// held-story contributes its one live rule, half-story both of its; the
	// retired rule is in neither total.
	if p.Rules != 3 || p.Automated != 1 {
		t.Errorf("automated rules = %d/%d, want 1/3", p.Automated, p.Rules)
	}
	if p.Proposed != 1 || p.Questions != 1 {
		t.Errorf("proposed = %d, questions = %d, want 1 and 1", p.Proposed, p.Questions)
	}
}

// livt:automates livt://mapping/show-opportunity-progress/rule/R-03
// A row leads to the mapping once the story has one and to its card until then,
// because that is where a reader can act on it in either state.
func TestOpportunityProgressRowsLinkToMappingOrStory(t *testing.T) {
	b := progressBuilder(t)
	_, _, tallies, err := b.buildMappings()
	if err != nil {
		t.Fatal(err)
	}

	o := &domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ機会"}
	p := b.progressOf(o, oneSlice("held-story", "unmapped-story"), tallies)

	rows := p.Releases[0].Stories
	if got, want := rows[0].Path, "../mapping/held-story.html"; got != want {
		t.Errorf("mapped row links to %q, want %q", got, want)
	}
	if got, want := rows[1].Path, "../story/unmapped-story.html"; got != want {
		t.Errorf("un-mapped row links to %q, want %q", got, want)
	}
	// The figures lead into the lists that already render these items, narrowed
	// through the filter bar's own parameter — a link on any other name lands on
	// the unfiltered list and shows every opportunity's items as this one's.
	for _, path := range []string{p.StoriesPath, p.QuestionsPath, p.ProposedPath, p.UnautomatedPath} {
		if !strings.Contains(path, "?"+opportunityFilterParam+"=") {
			t.Errorf("%q does not narrow on the filter bar's parameter", path)
		}
		if !strings.Contains(path, url.QueryEscape("デモ機会")) {
			t.Errorf("%q does not narrow to this opportunity", path)
		}
	}
	// The Tasks page keeps three lists; a chip lands on the one it counted
	// rather than at the top of all three.
	for _, want := range []string{tasksQuestionsAnchor, tasksProposedAnchor, tasksRulesAnchor} {
		if !strings.Contains(p.QuestionsPath+p.ProposedPath+p.UnautomatedPath, "#"+want) {
			t.Errorf("no chip lands on #%s", want)
		}
	}
}

// livt:automates livt://mapping/show-opportunity-progress/rule/R-01
// An opportunity nobody has mapped a journey for has no progress to read, so it
// gets no page and its own page links to none — the absence is the record, the
// same way it is for a canvas that was never filled in.
func TestOpportunityWithNoStoryMapGetsNoProgressPage(t *testing.T) {
	b := progressBuilder(t)
	if err := os.Remove(filepath.Join(b.USMDir, "demo.yaml")); err != nil {
		t.Fatal(err)
	}
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(b.OutDir, "opportunity-progress", "demo.html")); !os.IsNotExist(err) {
		t.Error("an opportunity with no map should get no progress page")
	}
	if html := readFile(t, filepath.Join(b.OutDir, "opportunity", "demo.html")); strings.Contains(html, "opportunity-progress/") {
		t.Error("the opportunity page should link to no progress page")
	}
}

// livt:automates livt://mapping/show-opportunity-progress/rule/R-01
// The opportunity's page is the way in and nothing more: what an opportunity is
// reads the same on every visit, and the reading that changes lives on its own
// page rather than being previewed beside the statement.
func TestOpportunityPageOnlyLeadsToTheProgress(t *testing.T) {
	b := progressBuilder(t)
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	html := readFile(t, filepath.Join(b.OutDir, "opportunity", "demo.html"))
	if !strings.Contains(html, "opportunity-progress/demo.html") {
		t.Error("the opportunity page does not lead to the progress page")
	}
	en := i18n.Of(i18n.En)
	for _, moved := range []string{"opportunity.mapped-stories", "opportunity.automated-rules", "opportunity.unmapped"} {
		if strings.Contains(html, en.Msg(moved)) {
			t.Errorf("%s belongs to the progress page, not the way in", moved)
		}
	}
}

// livt:automates livt://mapping/show-opportunity-progress/rule/R-02
// The breakdown is one row per story the opportunity took on, and every figure
// on it carries the whole it is part of.
func TestOpportunityProgressPageBreaksDownByStory(t *testing.T) {
	b := progressBuilder(t)
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	html := readFile(t, filepath.Join(b.OutDir, "opportunity-progress", "demo.html"))
	en := i18n.Of(i18n.En)
	for _, name := range []string{"押さえたストーリー", "途中のストーリー", "まだのストーリー"} {
		if !strings.Contains(html, name) {
			t.Errorf("story %q is missing from the breakdown", name)
		}
	}
	// The column is headed once rather than labelling every row, and that
	// heading is what says the fraction is coverage.
	if !strings.Contains(html, en.Msg("opportunity.coverage")) {
		t.Error("nothing on the page says the per-story fraction is rule coverage")
	}
	// A story with no mapping says so in words; 0/0 would read as a mapping
	// that found no rules.
	if got := strings.Count(html, en.Msg("opportunity.unmapped")); got != 1 {
		t.Errorf("got %d un-mapped rows, want 1", got)
	}
}

// livt:automates livt://mapping/show-opportunity-progress/rule/R-02
// A bare count does not say whether a board is nearly agreed, so every figure on
// the dashboard is drawn as a share of the whole it belongs to — and every
// figure gets a meter of its own, because what closes each of them differs and a
// segment inside another's bar is not a thing anyone can go and do.
func TestOpportunityDashboardMetersCarryTheirWhole(t *testing.T) {
	b := progressBuilder(t)
	_, _, tallies, err := b.buildMappings()
	if err != nil {
		t.Fatal(err)
	}
	o := &domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ機会"}
	if p := b.progressOf(o, oneSlice("held-story", "half-story"), tallies); p.QuestionsAsked != 2 {
		t.Errorf("questions asked = %d, want 2 (one open, one settled)", p.QuestionsAsked)
	}

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}
	html := readFile(t, filepath.Join(b.OutDir, "opportunity-progress", "demo.html"))

	en := i18n.Of(i18n.En)
	// Two of three stories mapped; of three live rules one is automated, one is
	// proposed and one is left waiting for a test — the proposal is not counted
	// as waiting for one; one of the two questions the boards asked is open.
	for label, want := range map[string]string{
		"opportunity.mapped-stories":    `>2<span class="font-normal text-gray-400">/3<`,
		"opportunity.unautomated-rules": `>1<span class="font-normal text-gray-400">/3<`,
		"opportunity.proposed-rules":    `>1<span class="font-normal text-gray-400">/3<`,
		"opportunity.open-questions":    `>1<span class="font-normal text-gray-400">/2<`,
	} {
		msg := en.Msg(label)
		if !strings.Contains(html, msg) {
			t.Errorf("%s has no meter of its own", msg)
			continue
		}
		if !strings.Contains(html[strings.Index(html, msg):], want) {
			t.Errorf("the %s meter is missing its whole: expected %q", msg, want)
		}
	}
}

// livt:automates livt://mapping/show-opportunity-progress/rule/R-02
// A proposal can carry a test ahead of its agreement, and then it is both
// proposed and automated. The un-automated meter counts what the Tasks page
// lists as waiting for a test — neither — rather than subtracting both counts
// from the whole, which took such a rule out twice and could go below zero.
func TestOpportunityUnautomatedMeterDoesNotSubtractAProposalTwice(t *testing.T) {
	b := progressBuilder(t)
	// half-story: one plain rule, one bare proposal, and one proposal whose
	// test was written first.
	writeFile(t, filepath.Join(b.MappingsDir, "half-story.yaml"),
		"rules:\n"+
			"  - id: R-01\n    name: まだのルール\n"+
			"  - id: R-02\n    name: 提案中のルール\n    status: proposed\n"+
			"  - id: R-03\n    name: 先にテストのある提案\n    status: proposed\n"+
			"questions: []\n")
	writeAutomations(t, b,
		"livt://mapping/held-story/rule/R-01",
		"livt://mapping/half-story/rule/R-03")
	_, _, tallies, err := b.buildMappings()
	if err != nil {
		t.Fatal(err)
	}
	o := &domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ機会"}
	p := b.progressOf(o, oneSlice("held-story", "half-story"), tallies)
	// Four live rules across the two boards: one automated, one plain, two
	// proposed of which one is also automated.
	if p.Rules != 4 || p.Automated != 2 || p.Proposed != 2 {
		t.Fatalf("rules/automated/proposed = %d/%d/%d, want 4/2/2", p.Rules, p.Automated, p.Proposed)
	}
	if p.Unautomated != 1 {
		t.Errorf("un-automated = %d, want 1: the plain rule alone is waiting for a test", p.Unautomated)
	}

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}
	html := readFile(t, filepath.Join(b.OutDir, "opportunity-progress", "demo.html"))
	msg := i18n.Of(i18n.En).Msg("opportunity.unautomated-rules")
	if at := strings.Index(html, msg); at < 0 {
		t.Fatalf("%s has no meter", msg)
	} else if want := `>1<span class="font-normal text-gray-400">/4<`; !strings.Contains(html[at:], want) {
		t.Errorf("the %s meter should read 1/4: expected %q", msg, want)
	}
}

// livt:automates livt://mapping/show-opportunity-progress/rule/R-01
// A fraction over a whole does not say which way it is meant to move, so the
// meters stand in two groups: stories under a burn-up heading, done at the
// whole, and what is still open under a burn-down one, done at zero. Each
// carries its own group's arrow, so a card read alone still says which it is.
func TestOpportunityDashboardSplitsBurnUpFromBurnDown(t *testing.T) {
	b := progressBuilder(t)
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}
	html := readFile(t, filepath.Join(b.OutDir, "opportunity-progress", "demo.html"))

	en := i18n.Of(i18n.En)
	up := strings.Index(html, en.Msg("opportunity.burn-up"))
	down := strings.Index(html, en.Msg("opportunity.burn-down"))
	if up < 0 || down < 0 {
		t.Fatalf("burn-up heading at %d, burn-down heading at %d: want both", up, down)
	}
	if up > down {
		t.Error("the burn-up group should lead: it is the figure that says how far discovery has got")
	}
	// The story count climbs to the whole; the rest fall to zero.
	for label, wantDown := range map[string]bool{
		"opportunity.mapped-stories":    false,
		"opportunity.unautomated-rules": true,
		"opportunity.proposed-rules":    true,
		"opportunity.open-questions":    true,
	} {
		msg := en.Msg(label)
		at := strings.Index(html, msg)
		if at < 0 {
			t.Errorf("%s has no meter", msg)
			continue
		}
		if gotDown := at > down; gotDown != wantDown {
			t.Errorf("the %s meter sits under the wrong heading: burn-down = %v, want %v", msg, gotDown, wantDown)
		}
		card := html[at:]
		if end := strings.Index(card, "</div>"); end > 0 {
			card = card[:end]
		}
		if gotDown := strings.Contains(card, "&darr;"); gotDown != wantDown {
			t.Errorf("the %s meter carries the wrong arrow: down = %v, want %v", msg, gotDown, wantDown)
		}
	}
}

// livt:automates livt://mapping/show-opportunity-progress/rule/R-04
// The hub says which opportunity is moving without being opened, carrying the
// same two figures its page leads with.
func TestOpportunitiesHubCarriesTheSameFigures(t *testing.T) {
	b := progressBuilder(t)
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	html := readFile(t, filepath.Join(b.OutDir, "opportunities.html"))
	en := i18n.Of(i18n.En)
	for _, label := range []string{"opportunity.mapped-stories", "opportunity.automated-rules"} {
		if !strings.Contains(html, en.Msg(label)) {
			t.Errorf("%s figure is missing from the tile", label)
		}
	}
	for _, fraction := range []string{">2<", ">1<"} {
		if !strings.Contains(html, fraction) {
			t.Errorf("tile does not carry %q", fraction)
		}
	}
}

// oneSlice is an opportunity whose map declares no release: every story lands
// in the single unscoped slice, which is the common case.
func oneSlice(keys ...string) []opportunityReleaseStories {
	return []opportunityReleaseStories{{Keys: keys}}
}
