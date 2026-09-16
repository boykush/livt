package diff

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Reading a revision goes through git the same way spec_version does: the
// binary is already what a livt repository is kept in, and an export is read
// once and thrown away. `git archive` rather than a worktree, because a build
// that only wants to read should leave nothing behind in .git.

// ErrNoRepo marks a root that git does not hold, so the caller can say the one
// thing that fixes it rather than passing git's own wording on.
var ErrNoRepo = errors.New("not a git repository")

// requireRepo fails before anything is exported, so a build asked for a diff
// outside a repository stops with the reason rather than with an empty one.
func requireRepo(root string) error {
	if err := exec.Command("git", "-C", root, "rev-parse", "--git-dir").Run(); err != nil {
		return fmt.Errorf("%w: %s", ErrNoRepo, root)
	}
	return nil
}

// resolveRev turns a revision into the short hash the page prints, and is what
// rejects one git cannot find. ^{commit} is what makes a tag or a branch resolve
// to the commit it points at, and a blob or a tree fail rather than be exported.
func resolveRev(root, rev string) (string, error) {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--short", rev+"^{commit}").Output()
	if err != nil {
		return "", fmt.Errorf("revision %q not found in %s", rev, root)
	}
	return strings.TrimSpace(string(out)), nil
}

// exportRev writes the livt repository's input directories at rev into dir.
// Directories the revision does not hold are skipped rather than failing the
// export: a revision from before a resource type existed is a legitimate side
// of a diff, and showing that it had none is the job.
func exportRev(root, rev, dir string, paths []string) error {
	held := heldPaths(root, rev, paths)
	if len(held) == 0 {
		return nil
	}
	args := append([]string{"-C", root, "archive", "--format=tar", rev, "--"}, held...)
	cmd := exec.Command("git", args...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := extract(out, dir); err != nil {
		_ = cmd.Wait()
		return err
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("git archive %s: %v: %s", rev, err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func heldPaths(root, rev string, paths []string) []string {
	var held []string
	for _, p := range paths {
		if exec.Command("git", "-C", root, "rev-parse", "-q", "--verify", rev+":"+p).Run() == nil {
			held = append(held, p)
		}
	}
	return held
}

// extract unpacks the archive with the standard library rather than piping into
// tar, so the export needs no second binary and every name is checked before it
// is written.
func extract(r io.Reader, dir string) error {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		target, ok := safeJoin(dir, header.Name)
		if !ok {
			continue
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := writeFile(target, tr); err != nil {
				return err
			}
		}
	}
}

// safeJoin refuses any name that would land outside the export directory. git
// writes repository-relative names, so this guards against a crafted archive
// rather than against git.
func safeJoin(dir, name string) (string, bool) {
	clean := filepath.Clean(filepath.FromSlash(name))
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", false
	}
	return filepath.Join(dir, clean), true
}

func writeFile(path string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}
