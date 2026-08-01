package builder

import (
	"fmt"
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
	if err := b.buildMappingsIndex([]mappingTile{{Key: "checkout", StoryName: "Checkout flow"}}, nil); err != nil {
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
	if err := b.buildMappingsIndex(nil, nil); err != nil {
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

// overview-open-questions R-01 EX-01 / R-02 EX-01 and overview-unautomated-rules
// R-01 EX-01 / R-02 EX-01: the Tasks page gathers both flanks in one place, each
// item carrying the story it came from and a link to its sticky.
func TestBuildTasksListsOpenQuestionsAndUnautomatedRules(t *testing.T) {
	b := emptyDirsBuilder(t)
	questions := []taskItem{{
		Kind: "question", ID: "Q-01", Text: "解決した疑問はどう扱うか",
		StoryKey: "overview-open-questions", StoryName: "未解決の疑問を横断で見渡す",
		StoryPath: "story/overview-open-questions.html", MappingPath: "mapping/overview-open-questions.html#question-Q-01",
	}}
	rules := []taskItem{{
		Kind: "rule", ID: "R-01", Text: "未自動化のルールは全実例マッピングを横断して一覧できる",
		StoryKey: "overview-unautomated-rules", StoryName: "未自動化のルールを横断で見渡す",
		StoryPath: "story/overview-unautomated-rules.html", MappingPath: "mapping/overview-unautomated-rules.html#rule-R-01",
	}}
	if err := b.buildTasks(questions, rules, nil); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "tasks.html"))

	for _, want := range []string{
		"解決した疑問はどう扱うか",                // a question, gathered from its mapping
		"未自動化のルールは全実例マッピングを横断して一覧できる", // an un-automated rule, gathered from another
		"未解決の疑問を横断で見渡す",               // the story a question came from
		"未自動化のルールを横断で見渡す",             // the story a rule came from
		`href="mapping/overview-open-questions.html#question-Q-01"`,
		`href="mapping/overview-unautomated-rules.html#rule-R-01"`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("tasks.html missing %q", want)
		}
	}
}

// overview-open-questions R-01 EX-03 and overview-unautomated-rules R-01 EX-03:
// each list says so when there is nothing left in it — an empty questions list
// and a fully automated master read differently.
func TestBuildTasksEmptyStatesReadPerList(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.buildTasks(nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "tasks.html"))

	for _, want := range []string{"No open questions.", "Every rule is automated."} {
		if !strings.Contains(html, want) {
			t.Fatalf("tasks.html missing empty state %q", want)
		}
	}
}

// overview-open-questions R-03 and overview-unautomated-rules R-03: one filter
// bar, one axis, driving every card on the page whichever list it sits in.
func TestBuildTasksFilterCoversBothLists(t *testing.T) {
	b := emptyDirsBuilder(t)
	opp := []opportunityRef{{Name: "協働ディスカバリー", Path: "story-map/協働ディスカバリー.html"}}
	questions := []taskItem{{Kind: "question", Text: "A question", StoryName: "S", Opportunities: opp}}
	rules := []taskItem{{Kind: "rule", Text: "A rule", StoryName: "S", Opportunities: opp}}
	if err := b.buildTasks(questions, rules, []string{"協働ディスカバリー"}); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "tasks.html"))

	// The attribute also appears in the script's selector, so match the element.
	if got := strings.Count(html, "data-opportunity-filter>"); got != 1 {
		t.Fatalf("filter bar rendered %d times, want one bar driving both lists", got)
	}
	if got := strings.Count(html, `data-opportunities="[&#34;協働ディスカバリー&#34;]"`); got != 2 {
		t.Fatalf("filter hook on %d cards, want both the question and the rule", got)
	}
	if !strings.Contains(html, `data-opportunity="協働ディスカバリー"`) {
		t.Fatal("expected a filter button carrying the opportunity name")
	}
	// R-03 EX-03: the selection rides in the query param, so a filtered Tasks
	// page is shareable and restores on load.
	for _, hook := range []string{"URLSearchParams", "'opportunity'", "history.replaceState"} {
		if !strings.Contains(html, hook) {
			t.Fatalf("expected the filter to sync the URL, missing %q", hook)
		}
	}
}

// Filtering to an opportunity a section has nothing for must not leave its
// heading claiming the unfiltered count over an empty gap. The section needs
// three things for the shared filter to put that right, and the filter has to
// actually use them — hooks nothing reads would render the same defect.
func TestBuildTasksSectionsStayHonestWhenAFilterEmptiesThem(t *testing.T) {
	b := emptyDirsBuilder(t)
	opp := []opportunityRef{{Name: "協働ディスカバリー", Path: "story-map/協働ディスカバリー.html"}}
	questions := []taskItem{{Kind: "question", Text: "A question", StoryName: "S", Opportunities: opp}}
	rules := []taskItem{{Kind: "rule", Text: "A rule", StoryName: "S", Opportunities: opp}}
	if err := b.buildTasks(questions, rules, []string{"協働ディスカバリー"}); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "tasks.html"))

	// Match the elements, not the bare attribute: the filter script mentions
	// each name again in its own selectors, which the last assertion checks for.
	for _, hook := range []string{
		"<section data-filter-scope", // the countable group
		"<span data-filter-count",    // the number to correct
		"<p data-filter-empty",       // the line explaining the gap
	} {
		if got := strings.Count(html, hook); got != 2 {
			t.Errorf("%s appears %d times, want one per section", hook, got)
		}
	}

	// The filtered-empty line ships hidden, so an unfiltered page does not carry
	// two empty states at once.
	if !strings.Contains(html, `data-filter-empty style="display: none"`) {
		t.Error("expected the filtered-empty line to start hidden")
	}
	// It has to say the filter emptied the section, not that the master is done —
	// "Every rule is automated." would be a lie about the master.
	for _, msg := range []string{
		"No open questions for this opportunity.",
		"Every rule for this opportunity is automated.",
	} {
		if !strings.Contains(html, msg) {
			t.Errorf("missing filtered-empty line %q", msg)
		}
	}

	// Without this the hooks are inert decoration and the count never moves.
	for _, selector := range []string{"[data-filter-scope]", "[data-filter-count]", "[data-filter-empty]"} {
		if !strings.Contains(html, selector) {
			t.Errorf("filter script never reads %s, so the section cannot correct itself", selector)
		}
	}
}

// An item whose story has no page still names its story; only the link drops.
func TestBuildTasksItemWithoutStoryPageStillNamesItsStory(t *testing.T) {
	b := emptyDirsBuilder(t)
	questions := []taskItem{{Kind: "question", Text: "A question", StoryName: "orphan-story"}}
	if err := b.buildTasks(questions, nil, nil); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "tasks.html"))

	if !strings.Contains(html, "orphan-story") {
		t.Fatal("expected the story name even with no page to link to")
	}
	if strings.Contains(html, `href="story/`) {
		t.Fatal("expected no story link when the story has no page")
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
	if !strings.Contains(html, `href="story-map/discovery.html"`) {
		t.Fatal("expected tile to link to the story map board")
	}
	if !strings.Contains(html, `src="story-map/discovery.html"`) {
		t.Fatal("expected story map preview iframe source")
	}
	if !strings.Contains(html, "Story Maps") {
		t.Fatal("expected Story Maps label")
	}
}

func TestBuildStoriesIndexRendersList(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.buildStoriesIndex([]storyItem{{Key: "first-story", Name: "First story"}}, nil); err != nil {
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

func readRendered(t *testing.T, path string) string {
	t.Helper()
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// mapHref mirrors how html/template normalizes an href: bytes outside the URL
// unreserved set (so every byte of a non-ASCII map name) become lowercase %xx.
// Chip links to Japanese-named maps therefore ship URL-encoded, which is also
// what lets a filtered URL round-trip (R-04).
func mapHref(prefix, name string) string {
	var b strings.Builder
	b.WriteString(prefix)
	for _, c := range []byte(name) {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '.' || c == '_' || c == '~' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02x", c)
		}
	}
	b.WriteString(".html")
	return b.String()
}

// R-02 EX-01: a story on one map carries a chip named after that opportunity
// (the map), linked to its board, in place of the old generic "Story Map" badge.
func TestBuildStoriesIndexShowsOpportunityChipNamedAfterTheMap(t *testing.T) {
	b := emptyDirsBuilder(t)
	items := []storyItem{{
		Key:  "filter-lists-by-opportunity",
		Name: "Filter lists by opportunity",
		Opportunities: []opportunityRef{
			{Name: "協働ディスカバリー", Path: "story-map/協働ディスカバリー.html"},
		},
	}}
	if err := b.buildStoriesIndex(items, nil); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "stories.html"))

	if !strings.Contains(html, "協働ディスカバリー") {
		t.Fatal("expected the opportunity (map) name on the story card")
	}
	if !strings.Contains(html, `href="`+mapHref("story-map/", "協働ディスカバリー")+`"`) {
		t.Fatal("expected the opportunity chip to link to the map board")
	}
	if strings.Contains(html, `>Story Map</span>`) {
		t.Fatal(`expected the generic "Story Map" badge to be replaced by an opportunity chip`)
	}
}

// R-02 EX-02: a story on several maps carries one chip per map.
func TestBuildStoriesIndexShowsAChipPerMapForMultiMapStory(t *testing.T) {
	b := emptyDirsBuilder(t)
	items := []storyItem{{
		Key:  "record-rule-automation",
		Name: "Record rule automation",
		Opportunities: []opportunityRef{
			{Name: "協働ディスカバリー", Path: "story-map/協働ディスカバリー.html"},
			{Name: "ディスカバリーと開発のギャップ", Path: "story-map/ディスカバリーと開発のギャップ.html"},
		},
	}}
	if err := b.buildStoriesIndex(items, nil); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "stories.html"))

	for _, name := range []string{"協働ディスカバリー", "ディスカバリーと開発のギャップ"} {
		if !strings.Contains(html, name) {
			t.Fatalf("expected a chip for map %q", name)
		}
		if !strings.Contains(html, `href="`+mapHref("story-map/", name)+`"`) {
			t.Fatalf("expected chip for %q to link to its board", name)
		}
	}
}

// R-02 EX-03: a story on no map carries no chip and its filter data is empty, so
// it matches no opportunity filter axis.
func TestBuildStoriesIndexStoryOnNoMapHasNoChipAndMatchesNoFilter(t *testing.T) {
	b := emptyDirsBuilder(t)
	items := []storyItem{{
		Key:  "orphan-story",
		Name: "Orphan story",
	}}
	if err := b.buildStoriesIndex(items, nil); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "stories.html"))

	if !strings.Contains(html, `data-opportunities="[]"`) {
		t.Fatal("expected an empty opportunity set so the card matches no filter")
	}
	if strings.Contains(html, "story-map/") {
		t.Fatal("expected no opportunity chip linking to a map board")
	}
}

// R-01: the Stories list renders opportunity filter controls, including a reset.
func TestBuildStoriesIndexRendersOpportunityFilterControls(t *testing.T) {
	b := emptyDirsBuilder(t)
	items := []storyItem{{Key: "s", Name: "S", Opportunities: []opportunityRef{{Name: "協働ディスカバリー", Path: "story-map/協働ディスカバリー.html"}}}}
	if err := b.buildStoriesIndex(items, []string{"協働ディスカバリー"}); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "stories.html"))

	if !strings.Contains(html, "data-opportunity-filter") {
		t.Fatal("expected an opportunity filter bar")
	}
	if !strings.Contains(html, `data-opportunity="協働ディスカバリー"`) {
		t.Fatal("expected a filter button carrying the opportunity name")
	}
	if !strings.Contains(html, `data-opportunity=""`) || !strings.Contains(html, ">All<") {
		t.Fatal("expected an All button that clears the filter")
	}
}

// R-03 & R-04: the filter runs client-side (no network) and reflects its state
// in the ?opportunity= query param so a shared URL restores the same view.
func TestBuildStoriesIndexFilterIsClientSideAndSyncsURL(t *testing.T) {
	b := emptyDirsBuilder(t)
	items := []storyItem{{Key: "s", Name: "S"}}
	if err := b.buildStoriesIndex(items, []string{"協働ディスカバリー"}); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "stories.html"))

	for _, hook := range []string{"<script>", "URLSearchParams", "'opportunity'", "history.replaceState"} {
		if !strings.Contains(html, hook) {
			t.Fatalf("expected client-side filter hook %q", hook)
		}
	}
	for _, network := range []string{"fetch(", "XMLHttpRequest"} {
		if strings.Contains(html, network) {
			t.Fatalf("filter must stay client-side, found %q", network)
		}
	}
}

// R-01 EX-02 & R-02: the Example Mappings list carries the same opportunity chip
// and filter, resolved through mapping → story → map.
func TestBuildMappingsIndexShowsOpportunityChipsAndFilter(t *testing.T) {
	b := emptyDirsBuilder(t)
	tiles := []mappingTile{{
		Key:       "filter-lists-by-opportunity",
		StoryName: "Filter lists by opportunity",
		Opportunities: []opportunityRef{
			{Name: "協働ディスカバリー", Path: "story-map/協働ディスカバリー.html"},
		},
	}}
	if err := b.buildMappingsIndex(tiles, []string{"協働ディスカバリー"}); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "index.html"))

	if !strings.Contains(html, `data-opportunities="[&#34;協働ディスカバリー&#34;]"`) {
		t.Fatal("expected the tile to carry its opportunity set as the filter hook")
	}
	if !strings.Contains(html, `href="`+mapHref("story-map/", "協働ディスカバリー")+`"`) {
		t.Fatal("expected the mapping tile's opportunity chip to link to the map board")
	}
	if !strings.Contains(html, "data-opportunity-filter") {
		t.Fatal("expected an opportunity filter bar on the Example Mappings list")
	}
}

// A mapping tile on no map has an empty filter set, matching no opportunity axis.
func TestBuildMappingsIndexTileOnNoMapHasEmptyFilterData(t *testing.T) {
	b := emptyDirsBuilder(t)
	tiles := []mappingTile{{Key: "orphan", StoryName: "Orphan"}}
	if err := b.buildMappingsIndex(tiles, nil); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "index.html"))
	if !strings.Contains(html, `data-opportunities="[]"`) {
		t.Fatal("expected an empty opportunity set on a mapping with no map")
	}
}
