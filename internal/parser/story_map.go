package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
	"gopkg.in/yaml.v3"
)

type storyMapYAML struct {
	Name       string         `yaml:"name"`
	Activities []activityYAML `yaml:"activities"`
	Releases   []releaseYAML  `yaml:"releases"`
}

type releaseYAML struct {
	Name    string         `yaml:"name"`
	Stories []string       `yaml:"stories"`
}

type activityYAML struct {
	Key   string     `yaml:"key"`
	Name  string     `yaml:"name"`
	Steps []stepYAML `yaml:"steps"`
}

type storyRefYAML struct {
	Key  string `yaml:"key"`
	Name string `yaml:"name"`
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
	for _, a := range raw.Activities {
		var steps []domain.Step
		for _, s := range a.Steps {
			var storyCards []domain.StoryCard
			for _, sk := range s.Stories {
				storyCards = append(storyCards, domain.StoryCard{
					Key:  domain.StoryKey{Value: sk.Key},
					Name: sk.Name,
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

	// Build key-to-name lookup from activities for release resolution
	cardNameByKey := make(map[string]string)
	for _, a := range activities {
		for _, s := range a.Steps {
			for _, sc := range s.Stories {
				cardNameByKey[sc.Key.Value] = sc.Name
			}
		}
	}

	var releases []domain.Release
	seen := make(map[string]bool)
	for _, r := range raw.Releases {
		var storyCards []domain.StoryCard
		for _, sk := range r.Stories {
			if seen[sk] {
				return nil, fmt.Errorf("story %q belongs to multiple releases", sk)
			}
			seen[sk] = true
			storyCards = append(storyCards, domain.StoryCard{
				Key:  domain.StoryKey{Value: sk},
				Name: cardNameByKey[sk],
			})
		}
		releases = append(releases, domain.Release{
			Name:    r.Name,
			Stories: storyCards,
		})
	}

	return &domain.StoryMap{
		Name:       raw.Name,
		Activities: activities,
		Releases:   releases,
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
