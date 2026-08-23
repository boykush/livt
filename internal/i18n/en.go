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
	"nav.opportunities":        "Opportunities",

	"label.story-map":          "Story Map",
	"label.opportunity-canvas": "Opportunity Canvas",

	"empty.opportunities":      "No opportunities yet.",
	"empty.opportunity-canvas": "No opportunity canvas yet.",

	"opportunity.opportunity": "Opportunity",

	// The canvas's three zones, and the ten boxes keyed by domain.CanvasBox.Key.
	// A prompt is the question the box asks: the heading alone does not say what
	// belongs in it, which is the whole reason Patton prints both.
	"canvas.zone.facts":    "Verifiable facts",
	"canvas.zone.solution": "Solution",
	"canvas.zone.value":    "Assumptions about value",
	"canvas.empty-box":     "Not filled in yet",

	"canvas.solution-ideas":             "Solution Ideas",
	"canvas.solution-ideas.prompt":      "List a specific product, feature, or enhancement idea.",
	"canvas.problems":                   "Problems",
	"canvas.problems.prompt":            "What problems do prospective users and customers have today that your solution addresses?",
	"canvas.users-and-customers":        "Users and Customers",
	"canvas.users-and-customers.prompt": "What types of users and customers have the challenges your solution addresses?",
	"canvas.solutions-today":            "Solutions Today",
	"canvas.solutions-today.prompt":     "How do users address their problems today?",
	"canvas.business-challenges":        "Business Challenges",
	"canvas.business-challenges.prompt": "How do these users' and customers' challenges impact your business?",
	"canvas.user-value":                 "What Will Users Do To Get Value?",
	"canvas.user-value.prompt":          "If your target audience has your solution, what will they do?",
	"canvas.user-metrics":               "User Metrics",
	"canvas.user-metrics.prompt":        "Given the story of what they do to get value, what could you measure that would show they actually did that?",
	"canvas.adoption-strategy":          "Adoption Strategy",
	"canvas.adoption-strategy.prompt":   "How will customers and users discover, learn to use, and adopt your solution?",
	"canvas.business-impact":            "Business Impact",
	"canvas.business-impact.prompt":     "What business performance metrics will be affected by the success of this solution?",
	"canvas.budget":                     "Budget",
	"canvas.budget.prompt":              "How much money or team time would you budget to solve this problem and achieve this outcome?",
}
