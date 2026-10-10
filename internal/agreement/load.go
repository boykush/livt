package agreement

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/parser"
)

// MappingsDir is where a livt repository keeps its example mappings.
const MappingsDir = "discoveries/example-mappings"

// LoadTree reads every mapping and its record from the working tree.
func LoadTree(root string) (map[string]Side, error) {
	mappings, err := parser.ParseAllExampleMappings(filepath.Join(root, filepath.FromSlash(MappingsDir)))
	if err != nil {
		return nil, err
	}
	sides := map[string]Side{}
	for _, m := range mappings {
		file, err := readFile(root, m.StoryKey.Value)
		if err != nil {
			return nil, err
		}
		sides[m.StoryKey.Value] = Side{Mapping: m, File: file}
	}
	return sides, nil
}

// LoadRev reads the same from a revision, exported to a directory thrown away
// afterwards so the repository is left as it was found.
func LoadRev(root, rev string) (map[string]Side, error) {
	if err := exec.Command("git", "-C", root, "rev-parse", "--verify", "--quiet", rev+"^{commit}").Run(); err != nil {
		return nil, fmt.Errorf("revision %q not found in %s", rev, root)
	}
	dir, err := os.MkdirTemp("", "livt-agreements-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	for _, sub := range []string{MappingsDir, Dir} {
		if err := export(root, rev, sub, dir); err != nil {
			return nil, err
		}
	}
	return LoadTree(dir)
}

// export copies one directory's files at rev into dir. A directory the
// revision does not hold is simply empty there.
func export(root, rev, sub, dir string) error {
	out, err := exec.Command("git", "-C", root, "ls-tree", "-r", "-z", "--name-only", rev, "--", sub).Output()
	if err != nil {
		return fmt.Errorf("list %s at %s: %w", sub, rev, err)
	}
	for _, name := range strings.Split(strings.TrimRight(string(out), "\x00"), "\x00") {
		if name == "" {
			continue
		}
		data, err := exec.Command("git", "-C", root, "show", rev+":"+name).Output()
		if err != nil {
			var exit *exec.ExitError
			if errors.As(err, &exit) {
				return fmt.Errorf("read %s at %s: %s", name, rev, strings.TrimSpace(string(exit.Stderr)))
			}
			return err
		}
		dest := filepath.Join(dir, filepath.FromSlash(path.Clean(name)))
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
