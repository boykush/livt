// Package automation carries what an implementation repository's tests say
// they automate. A citation is a claim its author made while writing the test;
// livt records it and shows it, and never decides whether it is true.
package automation

import "time"

// Marker separates a claim from a mention. A livt URI alone in a comment is a
// reference — production code cites rules for context too — so only a marked
// line counts as automation.
const Marker = "livt:automates"

// Citation is one line of one test claiming it automates a point of the spec.
// File is relative to the scanned repository's root, so the same report reads
// the same wherever it is checked out.
type Citation struct {
	URI  string `json:"uri"`
	File string `json:"file"`
	Line int    `json:"line"`
	// URL is where a reader can see that line on the forge. Empty when the
	// forge could not be identified: the link is decoration, and its absence
	// never changes whether the citation counts.
	URL string `json:"url,omitempty"`
}

// Report is one implementation repository's answer, covering the whole of it.
// That coverage is what lets a reader tell "scanned and not automated" from
// "never looked" — the absent report is the second, and needs no declaration
// anywhere to say so.
type Report struct {
	Repo        string     `json:"repo"`
	Rev         string     `json:"rev,omitempty"`
	GeneratedAt time.Time  `json:"generated_at"`
	Citations   []Citation `json:"citations"`
}
