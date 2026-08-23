package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"gopkg.in/yaml.v3"
)

type storyMapYAML struct {
	Name       string         `yaml:"name"`
	Activities []activityYAML `yaml:"activities"`
	Releases   []releaseYAML  `yaml:"releases"`
	Ubiquitous []string       `yaml:"ubiquitous"`
}

type releaseYAML struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type activityYAML struct {
	Key   string     `yaml:"key"`
	Name  string     `yaml:"name"`
	Steps []stepYAML `yaml:"steps"`
}

type storyRefYAML struct {
	Key     string `yaml:"key"`
	Name    string `yaml:"name"`
	Release string `yaml:"release"`
}

type stepYAML struct {
	Key     string         `yaml:"key"`
	Name    string         `yaml:"name"`
	Stories []storyRefYAML `yaml:"stories"`
}

func ParseStoryMap(path string) (*domain.StoryMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw storyMapYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var activities []domain.Activity
	storyNames := make(map[string]bool)
	for _, a := range raw.Activities {
		var steps []domain.Step
		for _, s := range a.Steps {
			var storyCards []domain.StoryCard
			for _, sk := range s.Stories {
				if sk.Name == "" {
					return nil, fmt.Errorf("story card in step %q is missing name", s.Key)
				}
				if storyNames[sk.Name] {
					return nil, fmt.Errorf("story card name %q is duplicated", sk.Name)
				}
				storyNames[sk.Name] = true
				storyCards = append(storyCards, domain.StoryCard{
					Key:     domain.StoryKey{Value: sk.Key},
					Name:    sk.Name,
					Release: sk.Release,
				})
			}
			steps = append(steps, domain.Step{
				Key:     s.Key,
				Name:    s.Name,
				Stories: storyCards,
			})
		}
		activities = append(activities, domain.Activity{
			Key:   a.Key,
			Name:  a.Name,
			Steps: steps,
		})
	}

	var releases []domain.Release
	releaseIDs := make(map[string]bool)
	for _, r := range raw.Releases {
		if r.ID == "" {
			return nil, fmt.Errorf("release %q is missing id", r.Name)
		}
		if releaseIDs[r.ID] {
			return nil, fmt.Errorf("release id %q is duplicated", r.ID)
		}
		releaseIDs[r.ID] = true
		releases = append(releases, domain.Release{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	for _, a := range activities {
		for _, s := range a.Steps {
			for _, sc := range s.Stories {
				if sc.Release != "" && !releaseIDs[sc.Release] {
					return nil, fmt.Errorf("story %q references unknown release %q", sc.Name, sc.Release)
				}
			}
		}
	}

	return &domain.StoryMap{
		Key:        strings.TrimSuffix(filepath.Base(path), ".yaml"),
		Name:       raw.Name,
		Activities: activities,
		Releases:   releases,
		Ubiquitous: raw.Ubiquitous,
	}, nil
}

func ParseAllStoryMaps(usmDir string) ([]*domain.StoryMap, error) {
	files, err := filepath.Glob(filepath.Join(usmDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	var maps []*domain.StoryMap
	for _, f := range files {
		sm, err := ParseStoryMap(f)
		if err != nil {
			return nil, err
		}
		maps = append(maps, sm)
	}

	return maps, nil
}
