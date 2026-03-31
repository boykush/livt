package recorder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Recorder struct {
	StoriesDir string
	USMDir     string
}

type storyRef struct {
	Key  string
	Name string
}

type storyMapYAML struct {
	Activities []activityYAML `yaml:"activities"`
}

type activityYAML struct {
	Steps []stepYAML `yaml:"steps"`
}

type stepYAML struct {
	Stories []storyRefYAML `yaml:"stories"`
}

type storyRefYAML struct {
	Key  string `yaml:"key"`
	Name string `yaml:"name"`
}

func (r *Recorder) Record() error {
	refs, err := r.collectStoryRefs()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(r.StoriesDir, 0755); err != nil {
		return err
	}

	for _, ref := range refs {
		path := filepath.Join(r.StoriesDir, ref.Key+".md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := createStory(path, ref); err != nil {
				return err
			}
			fmt.Printf("Created %s\n", path)
		} else {
			updated, err := updateStoryName(path, ref.Name)
			if err != nil {
				return err
			}
			if updated {
				fmt.Printf("Updated %s\n", path)
			}
		}
	}

	return nil
}

func (r *Recorder) collectStoryRefs() ([]storyRef, error) {
	files, err := filepath.Glob(filepath.Join(r.USMDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	var refs []storyRef
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}

		var raw storyMapYAML
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return nil, err
		}

		for _, a := range raw.Activities {
			for _, s := range a.Steps {
				for _, sr := range s.Stories {
					refs = append(refs, storyRef{Key: sr.Key, Name: sr.Name})
				}
			}
		}
	}

	return refs, nil
}

func createStory(path string, ref storyRef) error {
	content := fmt.Sprintf("---\nname: %s\n---\n\nAs a\nI want to\nSo that\n", ref.Name)
	return os.WriteFile(path, []byte(content), 0644)
}

func updateStoryName(path string, name string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	frontmatterDone := false
	updated := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "---" && !frontmatterDone {
			inFrontmatter = !inFrontmatter
			if !inFrontmatter {
				frontmatterDone = true
			}
			lines = append(lines, line)
			continue
		}

		if inFrontmatter && strings.HasPrefix(line, "name:") {
			newLine := "name: " + name
			if line != newLine {
				updated = true
				lines = append(lines, newLine)
			} else {
				lines = append(lines, line)
			}
			continue
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	if updated {
		return true, os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	}

	return false, nil
}
