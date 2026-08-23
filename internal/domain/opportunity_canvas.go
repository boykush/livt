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
// filling the boxes in, printed on the sheet. Key addresses the box's heading
// and its prompt — the question it asks, which the heading alone does not say —
// in the message catalogs, since both are read by a person and so follow the
// language the site is rendered in.
type CanvasBox struct {
	Key    string
	Number int
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
		{Number: 1, Key: "solution-ideas", Zone: ZoneSolution, Items: c.SolutionIdeas},
		{Number: 2, Key: "problems", Zone: ZoneFacts, Items: c.Problems},
		{Number: 3, Key: "users-and-customers", Zone: ZoneFacts, Items: c.UsersAndCustomers},
		{Number: 4, Key: "solutions-today", Zone: ZoneFacts, Items: c.SolutionsToday},
		{Number: 5, Key: "business-challenges", Zone: ZoneFacts, Items: c.BusinessChallenges},
		{Number: 6, Key: "user-value", Zone: ZoneValue, Items: c.UserValue},
		{Number: 7, Key: "user-metrics", Zone: ZoneValue, Items: c.UserMetrics},
		{Number: 8, Key: "adoption-strategy", Zone: ZoneValue, Items: c.AdoptionStrategy},
		{Number: 9, Key: "business-impact", Zone: ZoneValue, Items: c.BusinessImpact},
		{Number: 10, Key: "budget", Zone: ZoneSolution, Items: c.Budget},
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
