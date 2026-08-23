// Package config reads livt.yaml, the site config at the root of the livt
// repository. It is optional and every field falls back to a default, so a
// repository without the file builds exactly as it did before there was one.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strings"

	"github.com/boykush/livt/internal/i18n"
	"gopkg.in/yaml.v3"
)

// Path is where livt looks for the site config, relative to the directory the
// command runs in — the same place the input directories are resolved from.
const Path = "livt.yaml"

// Config is the site config as the builder consumes it: values already
// defaulted and validated, so nothing downstream has to re-check them.
type Config struct {
	// Lang is the language of the site chrome — nav labels, headings, legends,
	// empty states. It does not touch the livt repository's own prose, which is
	// rendered as written.
	Lang i18n.Lang
}

// Default is the config of a repository that ships no livt.yaml.
func Default() Config {
	return Config{Lang: i18n.Default}
}

// configYAML mirrors livt.yaml.
type configYAML struct {
	Lang string `yaml:"lang"`
}

// knownFields is what livt.yaml defines, in the order an error offers them.
// Field names are checked against this rather than by yaml's KnownFields so the
// message names the file's own vocabulary instead of livt's Go types.
var knownFields = []string{"lang"}

// Load reads the config at path. A missing file is not an error — it is the
// common case — but a malformed or unknown value is, since it means the site
// would be built as something other than what was asked for.
func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}

	var raw configYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return cfg, fmt.Errorf("%s: %w", path, err)
	}
	if err := checkFields(path, data); err != nil {
		return cfg, err
	}

	if raw.Lang != "" {
		lang := i18n.Lang(raw.Lang)
		if !lang.Valid() {
			return cfg, fmt.Errorf("%s: unknown lang %q (supported: %s)", path, raw.Lang, i18n.List())
		}
		cfg.Lang = lang
	}
	return cfg, nil
}

// checkFields rejects a field livt.yaml does not define. A typo like
// `language:` would otherwise leave the default in place while looking
// configured, and a config that is quietly ignored is worse than one that
// fails. The document is walked as nodes rather than a map so the field
// reported is the first one in the file, and so it can carry its line.
func checkFields(path string, data []byte) error {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	// An empty file, or a top level that is not a mapping — the struct decode
	// above has already had its say on the latter.
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil
	}
	// A mapping node holds its keys and values as alternating children.
	fields := doc.Content[0].Content
	for i := 0; i < len(fields); i += 2 {
		key := fields[i]
		if !slices.Contains(knownFields, key.Value) {
			return fmt.Errorf("%s:%d: unknown field %q (livt.yaml defines: %s)",
				path, key.Line, key.Value, strings.Join(knownFields, ", "))
		}
	}
	return nil
}
