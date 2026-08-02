package uri

// Kind names which of the livt repository's items a livt URI addresses.
type Kind string

const (
	KindMapping  Kind = "example mapping"
	KindRule     Kind = "rule"
	KindExample  Kind = "example"
	KindQuestion Kind = "question"
	KindStoryMap Kind = "story map"
	KindStory    Kind = "story"
	KindPersona  Kind = "persona"
	KindTerm     Kind = "ubiquitous language term"
)

// Templates lists every shape a livt URI can take, so a caller can show the
// full set without restating it.
var Templates = []string{
	MappingTemplate,
	RuleTemplate,
	ExampleTemplate,
	QuestionTemplate,
	StoryMapTemplate,
	StoryTemplate,
	PersonaTemplate,
	TermTemplate,
	ScopedTermTemplate,
}

// Parsed is a livt URI taken apart. Only the fields its Kind addresses are set.
type Parsed struct {
	Kind       Kind
	StoryKey   string // mapping, rule, example, question, story
	RuleID     string // rule, example
	ExampleID  string // example
	QuestionID string // question
	MapName    string // story map, percent-decoded
	PersonaKey string // persona
	TermKey    string // term
	TermCtx    string // term, empty when the term holds across contexts
}

// Parse recognises any livt URI shape. A caller that already knows which shape
// it expects — an MCP resource template bound to one handler — uses the
// per-kind parsers instead; Parse is for a URI arriving as free text, such as a
// CLI argument. The shapes are mutually exclusive, since ValidSegment rejects
// the extra segments a longer shape would leave behind.
func Parse(s string) (Parsed, bool) {
	if key, ruleID, exampleID, ok := ParseExample(s); ok {
		return Parsed{Kind: KindExample, StoryKey: key, RuleID: ruleID, ExampleID: exampleID}, true
	}
	if key, ruleID, ok := ParseRule(s); ok {
		return Parsed{Kind: KindRule, StoryKey: key, RuleID: ruleID}, true
	}
	if key, questionID, ok := ParseQuestion(s); ok {
		return Parsed{Kind: KindQuestion, StoryKey: key, QuestionID: questionID}, true
	}
	if key, ok := ParseMapping(s); ok {
		return Parsed{Kind: KindMapping, StoryKey: key}, true
	}
	if name, ok := ParseStoryMap(s); ok {
		return Parsed{Kind: KindStoryMap, MapName: name}, true
	}
	if key, ok := ParseStory(s); ok {
		return Parsed{Kind: KindStory, StoryKey: key}, true
	}
	if key, ok := ParsePersona(s); ok {
		return Parsed{Kind: KindPersona, PersonaKey: key}, true
	}
	if ctx, key, ok := ParseTerm(s); ok {
		return Parsed{Kind: KindTerm, TermKey: key, TermCtx: ctx}, true
	}
	return Parsed{}, false
}

// String rebuilds the URI in canonical form, so a story map named in plain text
// comes back percent-encoded.
func (p Parsed) String() string {
	switch p.Kind {
	case KindMapping:
		return Mapping(p.StoryKey)
	case KindRule:
		return Rule(p.StoryKey, p.RuleID)
	case KindExample:
		return Example(p.StoryKey, p.RuleID, p.ExampleID)
	case KindQuestion:
		return Question(p.StoryKey, p.QuestionID)
	case KindStoryMap:
		return StoryMap(p.MapName)
	case KindStory:
		return Story(p.StoryKey)
	case KindPersona:
		return Persona(p.PersonaKey)
	case KindTerm:
		return Term(p.TermCtx, p.TermKey)
	}
	return ""
}

// Page is where the URI lands on the generated site, relative to the output
// root. Deployment stays out of it: a deployed URL is this path with the site
// root prefixed, which is what lets one livt repository be served from several places.
func (p Parsed) Page() string {
	switch p.Kind {
	case KindMapping:
		return MappingPage(p.StoryKey)
	case KindRule:
		return RulePage(p.StoryKey, p.RuleID)
	case KindExample:
		return ExamplePage(p.StoryKey, p.RuleID, p.ExampleID)
	case KindQuestion:
		return QuestionPage(p.StoryKey, p.QuestionID)
	case KindStoryMap:
		return StoryMapPage(p.MapName)
	case KindStory:
		return StoryPage(p.StoryKey)
	case KindPersona:
		return PersonaPage(p.PersonaKey)
	case KindTerm:
		return TermPage(p.TermCtx, p.TermKey)
	}
	return ""
}
