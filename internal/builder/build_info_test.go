package builder

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/boykush/livt/internal/automation"
)

// fixedClock pins the build time so a page that prints it can be asserted on,
// and restores the real one when the test ends.
func fixedClock(t *testing.T, at time.Time) {
	t.Helper()
	was := now
	now = func() time.Time { return at }
	t.Cleanup(func() { now = was })
}

// writeReportAt is writeAutomations with a collection time, which is the whole
// subject here: a report that does not say when it was read cannot read as old.
func writeReportAt(t *testing.T, b Builder, repo, rev string, at time.Time, uris ...string) {
	t.Helper()
	citations := make([]automation.Citation, 0, len(uris))
	for i, u := range uris {
		citations = append(citations, automation.Citation{URI: u, File: "x_test.go", Line: i + 1})
	}
	data, err := json.MarshalIndent(automation.Report{Repo: repo, Rev: rev, GeneratedAt: at, Citations: citations}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(b.AutomationsDir, strings.ReplaceAll(repo, "/", "-")+".json"), string(data))
}

// livt:automates livt://mapping/collect-automations/rule/R-05/example/EX-06
// A report names when it was collected, and the build names when it ran, so the
// two are read against each other. Neither is useful alone: a timestamp with
// nothing to compare it to says nothing about whether the board is current.
func TestBuildInfoReadsAReportAgainstTheBuild(t *testing.T) {
	fixedClock(t, time.Date(2026, 9, 25, 8, 45, 0, 0, time.UTC))
	b := emptyDirsBuilder(t)
	b.LivtVersion = "v0.15.1"
	b.SpecVersion = "79aa298"
	writeReportAt(t, b, "acme/impl", "b61f1324d711dbdbb9d4ad66d62fc8657a8d65e8",
		time.Date(2026, 9, 12, 2, 4, 0, 0, time.UTC),
		"livt://mapping/checkout/rule/R-01")

	if err := b.buildInfo(); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "build.html"))

	for _, want := range []string{
		"v0.15.1",
		"79aa298",
		"2026-09-25 08:45", // the build
		"2026-09-12 02:04", // the report, thirteen days behind it
		"acme/impl",
		"b61f132",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("build.html missing %q", want)
		}
	}
	// The whole revision stays reachable, so the short one is an abbreviation
	// rather than the only thing recorded.
	if !strings.Contains(html, "b61f1324d711dbdbb9d4ad66d62fc8657a8d65e8") {
		t.Fatal("build.html abbreviates the revision without carrying it")
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-05/example/EX-06
// The page exists before any report does: with nothing collected it is the one
// place that says so, since every other surface now shows no automation at all.
func TestBuildInfoSaysSoWhenNothingWasCollected(t *testing.T) {
	fixedClock(t, time.Date(2026, 9, 25, 8, 45, 0, 0, time.UTC))
	b := emptyDirsBuilder(t)

	if err := b.buildInfo(); err != nil {
		t.Fatal(err)
	}
	html := readRendered(t, filepath.Join(b.OutDir, "build.html"))

	if !strings.Contains(html, "No report has been collected.") {
		t.Fatal("build.html does not say that nothing was collected")
	}
	if !strings.Contains(html, "2026-09-25 08:45") {
		t.Fatal("build.html drops the build time when there is no report to compare it to")
	}
}

// The foot of every page carries the build, so a board nobody has refreshed
// says so wherever it is read rather than only on the page about it.
func TestEveryPageFootsWithTheBuild(t *testing.T) {
	fixedClock(t, time.Date(2026, 9, 25, 8, 45, 0, 0, time.UTC))
	b := emptyDirsBuilder(t)
	b.LivtVersion = "v0.15.1"
	b.SpecVersion = "79aa298"

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}
	for _, page := range []string{"index.html", "tasks.html", "ubiquitous.html"} {
		html := readRendered(t, filepath.Join(b.OutDir, page))
		for _, want := range []string{"v0.15.1", "79aa298", "2026-09-25 08:45", "build.html"} {
			if !strings.Contains(html, want) {
				t.Fatalf("%s does not foot with %q", page, want)
			}
		}
	}
}
