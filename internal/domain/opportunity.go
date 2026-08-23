package domain

// Opportunity is one thing the product could take on: a user problem or need
// together with the business benefit of solving it, held as a single unit of
// consideration. Name is the short label carried on cards and filter chips,
// Body the statement of the opportunity itself — the same split Story makes
// between its name and its narrative, and for the same reason.
type Opportunity struct {
	Key  OpportunityKey
	Name string
	Body string
	Meta []MetaField
}
