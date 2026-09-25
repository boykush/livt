package builder

import (
	"os"
	"path/filepath"
	"time"
)

// now is the build's clock. Every page foots with the build time, which is the
// only reference a static page has for calling a report old, and tests replace
// this so that reference does not move under them.
var now = time.Now

// stamp writes both times the build page compares, UTC to the minute, so a
// report's and the build's are read against each other without arithmetic.
// Seconds are noise at the distance the two are apart.
func stamp(t time.Time) string { return t.UTC().Format("2006-01-02 15:04") }

// buildInfo renders build.html: the version, the spec revision and the time
// this site was made from, and the reports it read. Every other surface shows
// what was derived; this one says whose answer it is and how old, which is what
// tells a board that is current from one nobody has refreshed.
func (b *Builder) buildInfo() error {
	idx, err := b.automationIndex()
	if err != nil {
		return err
	}
	sb, err := b.sidebar("build", "")
	if err != nil {
		return err
	}
	reports := idx.Reports()
	rows := make([]buildReportRow, 0, len(reports))
	for _, r := range reports {
		rows = append(rows, buildReportRow{
			Repo:      r.Repo,
			Rev:       r.Rev,
			Short:     shortRev(r.Rev),
			Collected: stamp(r.GeneratedAt),
			Citations: len(r.Citations),
		})
	}
	f, err := os.Create(filepath.Join(b.OutDir, "build.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderBuildInfo(f, b.Lang, buildInfoView{
		Sidebar:     sb,
		LivtVersion: b.LivtVersion,
		SpecVersion: b.SpecVersion,
		Built:       sb.Built,
		Reports:     rows,
	})
}

// shortRev abbreviates a revision for the column, and says so with a dash when
// there is none: a scan of a dirty tree names no revision, and an empty cell
// would read as a column nobody filled in.
func shortRev(rev string) string {
	if rev == "" {
		return "—"
	}
	if len(rev) > 7 {
		return rev[:7]
	}
	return rev
}
