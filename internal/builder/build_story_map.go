package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
)

// storyMapBuild is what building the story maps leaves for the rest of the
// build: the opportunity refs a story's cards earn it, a preview tile per map,
// and the map drawn for each opportunity, by opportunity key.
type storyMapBuild struct {
	StoryOpportunities map[string][]opportunityRef
	Tiles              []storyMapTile
	MapByOpportunity   map[string]storyMapRef
	// StoriesByOpportunity is the reverse of StoryOpportunities: the keyed
	// stories an opportunity took on, in its map's own release slices. Only
	// keyed cards are in it, because a candidate with no card cannot have been
	// through an example mapping — counting it would make every opportunity look
	// less discovered than it is by the only measure the livt repository can take.
	StoriesByOpportunity map[string][]opportunityReleaseStories
}

// buildStoryMaps builds story map HTML pages. Paths in the returned refs are
// relative to the story/ directory, the deepest page that renders one. A story
// can sit on several maps, so each key maps to a slice in map order.
func (b *Builder) buildStoryMaps(opportunities map[string]*domain.Opportunity) (storyMapBuild, error) {
	maps, err := parser.ParseAllStoryMaps(b.USMDir)
	if err != nil {
		return storyMapBuild{}, err
	}

	out := storyMapBuild{
		StoryOpportunities:   make(map[string][]opportunityRef),
		MapByOpportunity:     make(map[string]storyMapRef),
		StoriesByOpportunity: make(map[string][]opportunityReleaseStories),
	}
	for _, sm := range maps {
		key := sm.OpportunityKey.Value
		ref := mapOpportunity(sm, opportunities)
		// The map page names its opportunity only when a file backs it; a map
		// standing in as its own opportunity would just link to itself.
		var own *opportunityRef
		var slices *releaseSlices
		if _, ok := opportunities[key]; ok {
			own = &ref
			out.MapByOpportunity[key] = storyMapRef{Name: sm.DisplayName(), Path: "../" + uri.StoryMapPage(key)}
			slices = newReleaseSlices()
			// Declared first and in the map's own order, so the slices read the
			// way the map draws them rather than the order stories happen to be
			// hung under the backbone.
			for i, r := range sm.Releases {
				slices.declare(r.ID, r.DisplayName(i))
			}
		}

		view := b.toStoryMapView(sm, own)
		view.Diff = b.diffMark("../", uri.StoryMap(key))
		outPath := filepath.Join(b.OutDir, uri.StoryMapPage(key))
		if err := b.buildStoryMap(outPath, view); err != nil {
			return storyMapBuild{}, err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))

		out.Tiles = append(out.Tiles, storyMapTile{Key: key, Name: sm.DisplayName(), Opportunity: own})

		// A key can recur across steps/releases within one map; add its chip once.
		seen := make(map[string]bool)
		for _, a := range sm.Activities {
			for _, s := range a.Steps {
				for _, sc := range s.Stories {
					if !sc.HasKey() || seen[sc.Key.Value] {
						continue
					}
					seen[sc.Key.Value] = true
					out.StoryOpportunities[sc.Key.Value] = append(out.StoryOpportunities[sc.Key.Value], ref)
					if slices != nil {
						slices.add(sc.Release, sc.Key.Value)
					}
				}
			}
		}
		if slices != nil {
			out.StoriesByOpportunity[key] = slices.rows()
		}
	}
	return out, nil
}

// opportunityReleaseStories is one release slice of an opportunity's stories.
// ID is empty for the stories its map left unscoped — which is every story on
// a map that declares no release at all, the common case.
type opportunityReleaseStories struct {
	ID   string
	Name string
	Keys []string
}

// releaseSlices buckets an opportunity's stories by the release its map puts
// them in, keeping the map's declared order. A slice nobody has hung a story
// under is dropped rather than drawn empty: the map is where the plan is read,
// and an empty row here would say a slice exists without saying anything about
// it.
type releaseSlices struct {
	order []string
	byID  map[string]*opportunityReleaseStories
}

func newReleaseSlices() *releaseSlices {
	return &releaseSlices{byID: make(map[string]*opportunityReleaseStories)}
}

func (s *releaseSlices) declare(id, name string) {
	if id == "" || s.byID[id] != nil {
		return
	}
	s.byID[id] = &opportunityReleaseStories{ID: id, Name: name}
	s.order = append(s.order, id)
}

// add files a story under its release. A story naming a release the map never
// declared is filed under it anyway, named by its own id — the story map is
// what reports an unknown release, and dropping the story here would lose it
// from a page that is counting.
func (s *releaseSlices) add(releaseID, storyKey string) {
	if s.byID[releaseID] == nil {
		s.byID[releaseID] = &opportunityReleaseStories{ID: releaseID, Name: releaseID}
		s.order = append(s.order, releaseID)
	}
	s.byID[releaseID].Keys = append(s.byID[releaseID].Keys, storyKey)
}

// rows returns the non-empty slices in declared order, with the unscoped one
// last: it is the remainder, not a stage of the plan.
func (s *releaseSlices) rows() []opportunityReleaseStories {
	var out []opportunityReleaseStories
	for _, id := range s.order {
		if id == "" {
			continue // the remainder is emitted last, below
		}
		if slice := s.byID[id]; len(slice.Keys) > 0 {
			out = append(out, *slice)
		}
	}
	if unscoped := s.byID[""]; unscoped != nil && len(unscoped.Keys) > 0 {
		out = append(out, *unscoped)
	}
	return out
}

type storyMapViewStory struct {
	Key  string
	Name string
	// Path is where the card leads, empty for a plain card.
	Path string
}

type storyMapViewStepHeader struct {
	Key  string
	Name string
}

type storyMapViewActivity struct {
	Key   string
	Name  string
	Steps []storyMapViewStepHeader
}

type storyMapViewActivityStepStories struct {
	StepStories [][]storyMapViewStory
}

type storyMapViewReleaseRow struct {
	Name       string
	Activities []storyMapViewActivityStepStories
}

type storyMapViewData struct {
	Name            string
	Opportunity     *opportunityRef
	Activities      []storyMapViewActivity
	ReleaseRows     []storyMapViewReleaseRow
	UnscopedStories *storyMapViewReleaseRow
	Ubiquitous      []termCard
}

type storyMapView struct {
	StoryMap storyMapViewData
	Diff     *diffMarkView
}

func (b *Builder) toStoryMapView(sm *domain.StoryMap, opportunity *opportunityRef) storyMapView {
	releaseIndexByID := make(map[string]int)
	for i, r := range sm.Releases {
		releaseIndexByID[r.ID] = i
	}

	var activities []storyMapViewActivity
	for _, a := range sm.Activities {
		var stepHeaders []storyMapViewStepHeader
		for _, s := range a.Steps {
			stepHeaders = append(stepHeaders, storyMapViewStepHeader{
				Key:  s.Key,
				Name: s.Name,
			})
		}
		activities = append(activities, storyMapViewActivity{
			Key:   a.Key,
			Name:  a.Name,
			Steps: stepHeaders,
		})
	}

	releaseRows, unscopedStories := b.buildReleaseRows(sm.Activities, sm.Releases, releaseIndexByID)

	return storyMapView{
		StoryMap: storyMapViewData{
			Name:            sm.DisplayName(),
			Opportunity:     opportunity,
			Activities:      activities,
			ReleaseRows:     releaseRows,
			UnscopedStories: unscopedStories,
			Ubiquitous:      b.resolveTermCards(sm.Ubiquitous),
		},
	}
}

func (b *Builder) buildReleaseRows(allActivities []domain.Activity, releases []domain.Release, releaseIndexByID map[string]int) ([]storyMapViewReleaseRow, *storyMapViewReleaseRow) {
	// Collect per-activity, per-step story grouping
	type stepGroup struct {
		perRelease map[int][]storyMapViewStory
		unscoped   []storyMapViewStory
	}
	type activityGroup struct {
		stepGroups []stepGroup
	}

	var actGroups []activityGroup
	for _, a := range allActivities {
		ag := activityGroup{}
		for _, s := range a.Steps {
			sg := stepGroup{perRelease: make(map[int][]storyMapViewStory)}
			for _, sc := range s.Stories {
				vs := storyMapViewStory{
					Key:  sc.Key.Value,
					Name: sc.Name,
					Path: b.storyCardPath(sc),
				}
				if sc.Release != "" {
					if idx, ok := releaseIndexByID[sc.Release]; ok {
						sg.perRelease[idx] = append(sg.perRelease[idx], vs)
						continue
					}
				}
				sg.unscoped = append(sg.unscoped, vs)
			}
			ag.stepGroups = append(ag.stepGroups, sg)
		}
		actGroups = append(actGroups, ag)
	}

	if len(releases) == 0 {
		row := storyMapViewReleaseRow{}
		for _, ag := range actGroups {
			actStories := storyMapViewActivityStepStories{}
			for _, sg := range ag.stepGroups {
				actStories.StepStories = append(actStories.StepStories, sg.unscoped)
			}
			row.Activities = append(row.Activities, actStories)
		}
		return nil, &row
	}

	var rows []storyMapViewReleaseRow
	for i, r := range releases {
		row := storyMapViewReleaseRow{Name: r.DisplayName(i)}
		for _, ag := range actGroups {
			actStories := storyMapViewActivityStepStories{}
			for _, sg := range ag.stepGroups {
				actStories.StepStories = append(actStories.StepStories, sg.perRelease[i])
			}
			row.Activities = append(row.Activities, actStories)
		}
		rows = append(rows, row)
	}

	// Unscoped stories
	hasUnscoped := false
	unscopedRow := storyMapViewReleaseRow{}
	for _, ag := range actGroups {
		actStories := storyMapViewActivityStepStories{}
		for _, sg := range ag.stepGroups {
			actStories.StepStories = append(actStories.StepStories, sg.unscoped)
			if len(sg.unscoped) > 0 {
				hasUnscoped = true
			}
		}
		unscopedRow.Activities = append(unscopedRow.Activities, actStories)
	}

	var unscoped *storyMapViewReleaseRow
	if hasUnscoped {
		unscoped = &unscopedRow
	}

	return rows, unscoped
}

// storyCardPath is where a card on the map leads: the story's page, or — for a
// key with an example mapping but no story file — the mapping, since the key
// names the story whether or not a file describes it. A card with neither, or
// with no key, leads nowhere.
func (b *Builder) storyCardPath(sc domain.StoryCard) string {
	switch {
	case !sc.HasKey():
		return ""
	case b.hasStoryPage(sc.Key):
		return "../" + uri.StoryPage(sc.Key.Value)
	case b.hasExampleMapping(sc.Key):
		return "../" + uri.MappingPage(sc.Key.Value)
	}
	return ""
}

func (b *Builder) buildStoryMap(path string, view storyMapView) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderStoryMap(f, b.Lang, view)
}
