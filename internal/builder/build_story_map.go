package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
)

// buildStoryMaps builds story map HTML pages and returns a map of story key to story map path
// (relative from story/ directory).
func (b *Builder) buildStoryMaps() (map[string]string, error) {
	maps, err := parser.ParseAllStoryMaps(b.USMDir)
	if err != nil {
		return nil, err
	}

	storyToMap := make(map[string]string)
	for _, sm := range maps {
		view := b.toStoryMapView(sm)
		outPath := filepath.Join(b.OutDir, "story-map", sm.Name+".html")
		if err := b.buildStoryMap(outPath, view); err != nil {
			return nil, err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))

		relativePath := "../story-map/" + sm.Name + ".html"
		for _, a := range sm.Activities {
			for _, s := range a.Steps {
				for _, sk := range s.Stories {
					storyToMap[sk.Value] = relativePath
				}
			}
		}
	}

	return storyToMap, nil
}

type storyMapViewStory struct {
	Key     string
	Name    string
	HasPage bool
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
	Activities      []storyMapViewActivity
	ReleaseRows     []storyMapViewReleaseRow
	UnscopedStories *storyMapViewReleaseRow
}

type storyMapView struct {
	StoryMap storyMapViewData
}

func (b *Builder) toStoryMapView(sm *domain.StoryMap) storyMapView {
	// Build story key → release index map
	storyRelease := make(map[string]int)
	for i, r := range sm.Releases {
		for _, sk := range r.Stories {
			storyRelease[sk.Value] = i
		}
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

	releaseRows, unscopedStories := b.buildReleaseRows(sm.Activities, sm.Releases, storyRelease)

	return storyMapView{
		StoryMap: storyMapViewData{
			Name:            sm.Name,
			Activities:      activities,
			ReleaseRows:     releaseRows,
			UnscopedStories: unscopedStories,
		},
	}
}

func (b *Builder) buildReleaseRows(allActivities []domain.Activity, releases []domain.Release, storyRelease map[string]int) ([]storyMapViewReleaseRow, *storyMapViewReleaseRow) {
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
			for _, sk := range s.Stories {
				vs := storyMapViewStory{
					Key:     sk.Value,
					Name:    b.resolveStoryName(sk),
					HasPage: b.hasStoryPage(sk),
				}
				if idx, ok := storyRelease[sk.Value]; ok {
					sg.perRelease[idx] = append(sg.perRelease[idx], vs)
				} else {
					sg.unscoped = append(sg.unscoped, vs)
				}
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
