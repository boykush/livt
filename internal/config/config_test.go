package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/i18n"
)

// The file is optional: livt predates it, and most repositories will never add
// one, so its absence has to mean "defaults" rather than "error".
func TestLoadWithoutTheFileTakesTheDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), Path))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Lang != i18n.Default {
		t.Fatalf("lang = %q, want the default %q", cfg.Lang, i18n.Default)
	}
}

func TestLoadReadsTheLanguage(t *testing.T) {
	path := writeConfig(t, "lang: ja\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Lang != i18n.Ja {
		t.Fatalf("lang = %q, want %q", cfg.Lang, i18n.Ja)
	}
}

func TestLoadDefaultsAnEmptyFile(t *testing.T) {
	cfg, err := Load(writeConfig(t, ""))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Lang != i18n.Default {
		t.Fatalf("lang = %q, want the default %q", cfg.Lang, i18n.Default)
	}
}

// An unsupported language cannot be silently downgraded to English: the site
// would be built in a language nobody asked for.
func TestLoadRejectsAnUnknownLanguage(t *testing.T) {
	_, err := Load(writeConfig(t, "lang: de\n"))
	if err == nil {
		t.Fatal("expected an error for an unsupported language")
	}
	for _, want := range []string{`"de"`, "en", "ja"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not mention %q", err, want)
		}
	}
}

// A misspelled key would leave the default in place while looking configured,
// so the error has to name the offending field and what livt.yaml does define.
func TestLoadRejectsAnUnknownField(t *testing.T) {
	_, err := Load(writeConfig(t, "lang: ja\nlanguage: ja\n"))
	if err == nil {
		t.Fatal("expected an error for a field livt.yaml does not define")
	}
	for _, want := range []string{`"language"`, ":2:", "lang"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not mention %q", err, want)
		}
	}
}

func TestLoadRejectsMalformedYAML(t *testing.T) {
	if _, err := Load(writeConfig(t, "lang: [ja\n")); err == nil {
		t.Fatal("expected an error for malformed YAML")
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), Path)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
