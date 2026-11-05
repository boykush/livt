package domain

// livt:automates livt://mapping/filter-lists-by-opportunity/rule/R-02/example/EX-06
// livt:automates livt://mapping/filter-lists-by-opportunity/rule/R-02/example/EX-07

import (
	"reflect"
	"testing"
)

func oneCardMap(opportunity, story, release string) *StoryMap {
	return &StoryMap{
		OpportunityKey: OpportunityKey{Value: opportunity},
		Activities: []Activity{{Steps: []Step{{Stories: []StoryCard{
			{Key: StoryKey{Value: story}, Release: release},
		}}}}},
	}
}

func TestDeriveOpportunityStoriesJoinsBothWaysOnce(t *testing.T) {
	maps := []*StoryMap{oneCardMap("mapped", "on-map", "r1")}
	stories := []*Story{
		{Key: StoryKey{Value: "on-map"}, Opportunity: OpportunityKey{Value: "mapped"}},
		{Key: StoryKey{Value: "card-only"}, Opportunity: OpportunityKey{Value: "unmapped"}},
		{Key: StoryKey{Value: "both"}, Opportunity: OpportunityKey{Value: "mapped"}},
		{Key: StoryKey{Value: "none"}},
	}
	maps = append(maps, oneCardMap("other", "both", ""))

	got := DeriveOpportunityStories(maps, stories)

	wantByStory := map[StoryKey][]OpportunityKey{
		{Value: "on-map"}:    {{Value: "mapped"}},
		{Value: "card-only"}: {{Value: "unmapped"}},
		{Value: "both"}:      {{Value: "other"}, {Value: "mapped"}},
	}
	if !reflect.DeepEqual(got.ByStory, wantByStory) {
		t.Errorf("ByStory = %v, want %v", got.ByStory, wantByStory)
	}
	wantMapped := []OpportunityStory{
		{Key: StoryKey{Value: "on-map"}, OnMap: true, Release: "r1"},
		{Key: StoryKey{Value: "both"}},
	}
	if m := got.ByOpportunity[OpportunityKey{Value: "mapped"}]; !reflect.DeepEqual(m, wantMapped) {
		t.Errorf("mapped's stories = %v, want %v", m, wantMapped)
	}
	if u := got.ByOpportunity[OpportunityKey{Value: "unmapped"}]; len(u) != 1 || u[0].OnMap {
		t.Errorf("unmapped's stories = %v, want card-only off the map", u)
	}
}
