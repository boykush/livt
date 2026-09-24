package automation

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The marker is spelled through the constant so these fixtures are not read as
// citations when livt scans itself — the same reason the scan skips a marker
// written inside code.
var (
	marked = "package p\n\n// " + Marker + " livt://mapping/checkout/rule/R-02/example/EX-01\nfunc TestCard(t *testing.T) {}\n"
	plain  = "package p\n\nfunc TestCard(t *testing.T) {}\n"
)

// gitRepo makes a scratch repository, so each test describes the whole history
// it then asks about.
func gitRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	git(t, root, "init", "-q")
	git(t, root, "config", "user.email", "test@example.com")
	git(t, root, "config", "user.name", "test")
	git(t, root, "config", "commit.gpgsign", "false")
	return root
}

func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func write(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

func commit(t *testing.T, root string) string {
	t.Helper()
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "commit")
	return git(t, root, "rev-parse", "HEAD")
}

func changed(t *testing.T, root, base, head string) bool {
	t.Helper()
	got, err := Changed(root, base, head)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

// livt:automates livt://mapping/collect-automations/rule/R-07/example/EX-01
func TestChangedWhenAMarkerIsAdded(t *testing.T) {
	root := gitRepo(t)
	write(t, root, map[string]string{"card_test.go": plain})
	base := commit(t, root)
	write(t, root, map[string]string{"card_test.go": marked})

	if !changed(t, root, base, commit(t, root)) {
		t.Error("an added marker did not ask for a scan")
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-07/example/EX-02
func TestChangedWhenAMarkerIsRemoved(t *testing.T) {
	for _, tc := range []struct {
		name string
		drop func(root string)
	}{
		{"the line goes", func(root string) {
			write(t, root, map[string]string{"card_test.go": plain})
		}},
		{"the file goes", func(root string) {
			if err := os.Remove(filepath.Join(root, "card_test.go")); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := gitRepo(t)
			write(t, root, map[string]string{"card_test.go": marked})
			base := commit(t, root)
			tc.drop(root)

			if !changed(t, root, base, commit(t, root)) {
				t.Error("a claim outlived its test without asking for a scan")
			}
		})
	}
}

// livt:automates livt://mapping/collect-automations/rule/R-07/example/EX-03
func TestUnchangedWhenTheClaimsStandStill(t *testing.T) {
	t.Run("a marker only moves", func(t *testing.T) {
		root := gitRepo(t)
		write(t, root, map[string]string{"card_test.go": marked})
		base := commit(t, root)
		// Renamed and pushed down a line; the marker line itself is untouched.
		git(t, root, "mv", "card_test.go", "expiry_test.go")
		write(t, root, map[string]string{"expiry_test.go": "//go:build unit\n\n" + marked})

		if changed(t, root, base, commit(t, root)) {
			t.Error("a moved line asked for a scan it cannot change the answer to")
		}
	})

	t.Run("only the report changes", func(t *testing.T) {
		root := gitRepo(t)
		write(t, root, map[string]string{
			"card_test.go":               marked,
			"automations/acme/impl.json": `{"rev":"aaa","citations":[{"livt_uri":"livt://mapping/checkout/rule/R-02"}]}`,
		})
		base := commit(t, root)
		write(t, root, map[string]string{
			"automations/acme/impl.json": `{"rev":"bbb","citations":[{"livt_uri":"livt://mapping/checkout/rule/R-02"}]}`,
		})

		if changed(t, root, base, commit(t, root)) {
			t.Error("collecting the report re-triggered the collection that wrote it")
		}
	})
}
