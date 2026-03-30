package builder

import "os"

func (b *Builder) buildIndex(path string, entries []IndexEntry) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderIndex(f, entries)
}
