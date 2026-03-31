package parser

import (
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
	"gopkg.in/yaml.v3"
)

type storyMapYAML struct {
	Name       string         `yaml:"name"`
	Activities []activityYAML `yaml:"activities"`
}

type activityYAML struct {
	Key   string     `yaml:"key"`
	Name  string     `yaml:"name"`
	Steps []stepYAML `yaml:"steps"`
}

type stepYAML struct {
	Key     string   `yaml:"key"`
	Name    string   `yaml:"name"`
	Stories []string `yaml:"stories"`
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
			var storyKeys []domain.StoryKey
			for _, sk := range s.Stories {
				storyKeys = append(storyKeys, domain.StoryKey{Value: sk})
			}
			steps = append(steps, domain.Step{
				Key:     s.Key,
				Name:    s.Name,
				Stories: storyKeys,
			})
		}
		activities = append(activities, domain.Activity{
			Key:   a.Key,
			Name:  a.Name,
			Steps: steps,
		})
	}

	return &domain.StoryMap{
		Name:       raw.Name,
		Activities: activities,
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
