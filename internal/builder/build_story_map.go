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
// and the maps mapped for each opportunity, by opportunity key.
type storyMapBuild struct {
	StoryOpportunities map[string][]opportunityRef
	Tiles              []storyMapTile
	MapsByOpportunity  map[string][]storyMapRef
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
		StoryOpportunities: make(map[string][]opportunityRef),
		MapsByOpportunity:  make(map[string][]storyMapRef),
	}
	for _, sm := range maps {
		ref := mapOpportunity(sm, opportunities)
		// The map page names its opportunity only when a file backs it; a map
		// standing in as its own opportunity would just link to itself.
		var own *opportunityRef
		if _, ok := opportunities[sm.Key]; ok {
			own = &ref
			out.MapsByOpportunity[sm.Key] = append(out.MapsByOpportunity[sm.Key],
				storyMapRef{Name: sm.Name, Path: "../" + uri.StoryMapPage(sm.Name)})
		}

		view := b.toStoryMapView(sm, own)
		outPath := filepath.Join(b.OutDir, uri.StoryMapPage(sm.Name))
		if err := b.buildStoryMap(outPath, view); err != nil {
			return storyMapBuild{}, err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))

		out.Tiles = append(out.Tiles, storyMapTile{Name: sm.Name, Opportunity: own})

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
				}
			}
		}
	}

	return out, nil
}

type storyMapViewStory struct {
	Key    string
	Name   string
	Opened bool
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
			Name:            sm.Name,
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
					Key:    sc.Key.Value,
					Name:   sc.Name,
					Opened: sc.HasKey() && b.hasStoryPage(sc.Key),
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

func (b *Builder) buildStoryMap(path string, view storyMapView) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderStoryMap(f, view)
}
