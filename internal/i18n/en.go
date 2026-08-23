package i18n

// en is the reference catalog: every other language is checked against its key
// set, so a message added here is a translation owed everywhere else.
var en = Catalog{
	// The <html lang> attribute of every page.
	"lang.code": "en",

	"nav.example-mappings": "Example Mappings",
	"nav.story-maps":       "Story Maps",
	"nav.stories":          "Stories",
	"nav.ubiquitous":       "Ubiquitous Language",
	"nav.tasks":            "Tasks",
	"nav.tagline-1":        "Collaborate on board.",
	"nav.tagline-2":        "Make it living in text.",

	// Sticky kinds, shared by the board legends, the filter bars and the
	// chips — one word each, so the same key serves wherever it appears.
	"label.story":           "Story",
	"label.rule":            "Rule",
	"label.example":         "Example",
	"label.question":        "Question",
	"label.questions":       "Questions",
	"label.opportunity":     "Opportunity",
	"label.activity":        "Activity",
	"label.step":            "Step",
	"label.context":         "Context",
	"label.example-mapping": "Example Mapping",

	"filter.all": "All",

	"badge.copy-link": "Copy link",
	"badge.copied":    "✓ Copied",

	"empty.example-mappings": "No example mappings yet.",
	"empty.story-maps":       "No story maps yet.",
	"empty.stories":          "No stories found.",
	"empty.terms":            "No terms found.",

	"tasks.questions":                "Open Questions",
	"tasks.questions-hint":           "closed by a conversation",
	"tasks.questions-empty":          "No open questions.",
	"tasks.questions-empty-filtered": "No open questions for this opportunity.",
	"tasks.rules":                    "Un-automated Rules",
	"tasks.rules-hint":               "closed by a test",
	"tasks.rules-empty":              "Every rule is automated.",
	"tasks.rules-empty-filtered":     "Every rule for this opportunity is automated.",

	"glossary.term":       "Term",
	"glossary.key":        "Key",
	"glossary.definition": "Definition",

	"story.metadata":    "Metadata",
	"story.related":     "Related",
	"story.description": "Description",

	"mapping.automated-legend": "Automated rule",
	"mapping.automated-badge":  "✓ automated",
	"mapping.automated-title":  "Automated by tests",
}
