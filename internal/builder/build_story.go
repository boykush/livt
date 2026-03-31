package builder

import (
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
)

func (b *Builder) hasExampleMapping(storyKey domain.StoryKey) bool {
	path := filepath.Join(b.MappingsDir, storyKey.Value+".yaml")
	_, err := os.Stat(path)
	return err == nil
}

func (b *Builder) buildStory(path string, story *domain.Story, mappingPath, storyMapPath string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderStory(f, story, mappingPath, storyMapPath)
}
