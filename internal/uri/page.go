package uri

// Where each livt URI lands on the generated site: a path relative to the
// output root, plus the in-page anchor for the points that live on a board.
// The deployment URL is prefixed onto these only at render time, so neither the
// master nor this package knows where the site is deployed. Deriving it here
// keeps the anchors the SSG writes and the URLs the CLI resolves in one piece.

// MappingPage is where Mapping lands.
func MappingPage(storyKey string) string {
	return "mapping/" + storyKey + ".html"
}

// RulePage is where Rule lands: its sticky on the story's board.
func RulePage(storyKey, ruleID string) string {
	return MappingPage(storyKey) + "#" + RuleAnchor(ruleID)
}

// ExamplePage is where Example lands.
func ExamplePage(storyKey, ruleID, exampleID string) string {
	return MappingPage(storyKey) + "#" + ExampleAnchor(ruleID, exampleID)
}

// QuestionPage is where Question lands.
func QuestionPage(storyKey, questionID string) string {
	return MappingPage(storyKey) + "#" + QuestionAnchor(questionID)
}

// StoryPage is where Story lands.
func StoryPage(storyKey string) string {
	return "story/" + storyKey + ".html"
}

// StoryMapPage is where StoryMap lands. The name stays raw here: the URI
// percent-encodes it to hold one path segment, but the build writes the file
// under the display name itself.
func StoryMapPage(name string) string {
	return "story-map/" + name + ".html"
}

// TermPage is where Term lands. The whole ubiquitous language is one table, so
// a term's anchor is its key alone.
func TermPage(termKey string) string {
	return "ubiquitous.html#" + termKey
}

// RuleAnchor is the id of a rule's sticky on its board.
func RuleAnchor(ruleID string) string {
	return "rule-" + ruleID
}

// ExampleAnchor is the id of an example's sticky. The rule is part of it for
// the same reason it is part of the URI: example ids are numbered within their
// rule, so one board holds an EX-01 under every rule.
func ExampleAnchor(ruleID, exampleID string) string {
	return RuleAnchor(ruleID) + "-example-" + exampleID
}

// QuestionAnchor is the id of a question's sticky.
func QuestionAnchor(questionID string) string {
	return "question-" + questionID
}

// StoryCardAnchor is the id of a story's card on a story map. A story's own
// livt URI lands on its page, not here, so this anchor serves only a link
// shared from the board.
func StoryCardAnchor(storyKey string) string {
	return "story-" + storyKey
}
