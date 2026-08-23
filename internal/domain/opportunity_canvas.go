package domain

// OpportunityCanvas is one filled-in Opportunity Canvas: the ten boxes Jeff
// Patton's canvas holds an opportunity in, each carrying the stickies that sat
// in it on the board. OpportunityKey comes from the filename, which is what
// joins a canvas to its opportunity — the same filename join that ties an
// example mapping to its story.
type OpportunityCanvas struct {
	OpportunityKey     OpportunityKey
	SolutionIdeas      []string
	Problems           []string
	UsersAndCustomers  []string
	SolutionsToday     []string
	BusinessChallenges []string
	UserValue          []string
	UserMetrics        []string
	AdoptionStrategy   []string
	BusinessImpact     []string
	Budget             []string
	Ubiquitous         []string
}

// Canvas zones. The sheet is divided three ways, and the division is the
// method's own: what can be checked, what is being proposed, and what is only
// assumed until the thing ships. Reading left to right walks back from the idea
// to the problem it solves, then forward to the value it would create.
const (
	ZoneFacts    = "facts"
	ZoneSolution = "solution"
	ZoneValue    = "value"
)

// CanvasBox is one box of the canvas. Number is the order Patton recommends
// filling the boxes in, printed on the sheet. Prompt is the question the box
// asks, kept beside the heading because the heading alone does not say what
// belongs in it.
type CanvasBox struct {
	Key    string
	Number int
	Name   string
	Prompt string
	Zone   string
	Items  []string
}

// CanvasPanels is the sheet's three columns, which is how the canvas is laid
// out rather than how it is filled in: the facts on the left, the solution
// down the middle, the value on the right.
type CanvasPanels struct {
	Facts    []CanvasBox
	Solution []CanvasBox
	Value    []CanvasBox
}

// Boxes returns the ten boxes in the order Patton recommends filling them.
// Empty boxes are returned too: a blank box is the visible record of a question
// the opportunity has not answered yet, and dropping it would quietly turn a
// gap into a tidy canvas.
func (c *OpportunityCanvas) Boxes() []CanvasBox {
	return []CanvasBox{
		{Number: 1, Key: "solution-ideas", Name: "Solution Ideas", Zone: ZoneSolution, Items: c.SolutionIdeas,
			Prompt: "List a specific product, feature, or enhancement idea."},
		{Number: 2, Key: "problems", Name: "Problems", Zone: ZoneFacts, Items: c.Problems,
			Prompt: "What problems do prospective users and customers have today that your solution addresses?"},
		{Number: 3, Key: "users-and-customers", Name: "Users and Customers", Zone: ZoneFacts, Items: c.UsersAndCustomers,
			Prompt: "What types of users and customers have the challenges your solution addresses?"},
		{Number: 4, Key: "solutions-today", Name: "Solutions Today", Zone: ZoneFacts, Items: c.SolutionsToday,
			Prompt: "How do users address their problems today?"},
		{Number: 5, Key: "business-challenges", Name: "Business Challenges", Zone: ZoneFacts, Items: c.BusinessChallenges,
			Prompt: "How do these users' and customers' challenges impact your business?"},
		{Number: 6, Key: "user-value", Name: "What Will Users Do To Get Value?", Zone: ZoneValue, Items: c.UserValue,
			Prompt: "If your target audience has your solution, what will they do?"},
		{Number: 7, Key: "user-metrics", Name: "User Metrics", Zone: ZoneValue, Items: c.UserMetrics,
			Prompt: "Given the story of what they do to get value, what could you measure that would show they actually did that?"},
		{Number: 8, Key: "adoption-strategy", Name: "Adoption Strategy", Zone: ZoneValue, Items: c.AdoptionStrategy,
			Prompt: "How will customers and users discover, learn to use, and adopt your solution?"},
		{Number: 9, Key: "business-impact", Name: "Business Impact", Zone: ZoneValue, Items: c.BusinessImpact,
			Prompt: "What business performance metrics will be affected by the success of this solution?"},
		{Number: 10, Key: "budget", Name: "Budget", Zone: ZoneSolution, Items: c.Budget,
			Prompt: "How much money or team time would you budget to solve this problem and achieve this outcome?"},
	}
}

// Panels groups the boxes into the sheet's three columns, each in the order it
// is printed down the page.
func (c *OpportunityCanvas) Panels() CanvasPanels {
	var p CanvasPanels
	for _, b := range c.Boxes() {
		switch b.Zone {
		case ZoneFacts:
			p.Facts = append(p.Facts, b)
		case ZoneSolution:
			p.Solution = append(p.Solution, b)
		case ZoneValue:
			p.Value = append(p.Value, b)
		}
	}
	return p
}
