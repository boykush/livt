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

type storyMapViewReleaseSlice struct {
	Name    string
	Stories []storyMapViewStory
}

type storyMapViewStep struct {
	Key           string
	Name          string
	ReleaseSlices []storyMapViewReleaseSlice
}

type storyMapViewActivity struct {
	Key   string
	Name  string
	Steps []storyMapViewStep
}

type storyMapViewData struct {
	Name       string
	Activities []storyMapViewActivity
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
		var steps []storyMapViewStep
		for _, s := range a.Steps {
			releaseSlices := b.groupStoriesByRelease(s.Stories, sm.Releases, storyRelease)
			steps = append(steps, storyMapViewStep{
				Key:           s.Key,
				Name:          s.Name,
				ReleaseSlices: releaseSlices,
			})
		}
		activities = append(activities, storyMapViewActivity{
			Key:   a.Key,
			Name:  a.Name,
			Steps: steps,
		})
	}

	return storyMapView{
		StoryMap: storyMapViewData{
			Name:       sm.Name,
			Activities: activities,
		},
	}
}

func (b *Builder) groupStoriesByRelease(storyKeys []domain.StoryKey, releases []domain.Release, storyRelease map[string]int) []storyMapViewReleaseSlice {
	if len(releases) == 0 {
		// No releases defined: single slice with no name (R-05)
		var stories []storyMapViewStory
		for _, sk := range storyKeys {
			stories = append(stories, storyMapViewStory{
				Key:     sk.Value,
				Name:    b.resolveStoryName(sk),
				HasPage: b.hasStoryPage(sk),
			})
		}
		if len(stories) == 0 {
			return nil
		}
		return []storyMapViewReleaseSlice{{Stories: stories}}
	}

	// Group stories by release index
	perRelease := make(map[int][]storyMapViewStory)
	var unscoped []storyMapViewStory
	for _, sk := range storyKeys {
		vs := storyMapViewStory{
			Key:     sk.Value,
			Name:    b.resolveStoryName(sk),
			HasPage: b.hasStoryPage(sk),
		}
		if idx, ok := storyRelease[sk.Value]; ok {
			perRelease[idx] = append(perRelease[idx], vs)
		} else {
			unscoped = append(unscoped, vs)
		}
	}

	var slices []storyMapViewReleaseSlice
	for i, r := range releases {
		if stories, ok := perRelease[i]; ok {
			slices = append(slices, storyMapViewReleaseSlice{
				Name:    r.Name,
				Stories: stories,
			})
		}
	}
	if len(unscoped) > 0 {
		slices = append(slices, storyMapViewReleaseSlice{Stories: unscoped})
	}

	return slices
}

func (b *Builder) buildStoryMap(path string, view storyMapView) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderStoryMap(f, view)
}
