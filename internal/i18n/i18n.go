// Package i18n holds the message catalogs for the site chrome — nav labels,
// headings, legends, empty states. The livt repository's own prose stays in
// whatever language the team writes it in; only what livt renders around it is
// translated here.
package i18n

import "fmt"

// Lang is a language the site chrome can be rendered in.
type Lang string

const (
	En Lang = "en"
	Ja Lang = "ja"
)

// Default is what livt renders in when the site config names no language.
const Default = En

// Langs are the supported languages, in the order an error message lists them.
var Langs = []Lang{En, Ja}

// Valid reports whether the language has a catalog, so the site config can
// reject an unknown one at load time rather than at render time.
func (l Lang) Valid() bool {
	_, ok := catalogs[l]
	return ok
}

// List names every supported language, for an error telling the user what to
// pick instead.
func List() string {
	s := ""
	for i, l := range Langs {
		if i > 0 {
			s += ", "
		}
		s += string(l)
	}
	return s
}

// Catalog is one language's messages, keyed by the ids the templates use.
type Catalog map[string]string

// T looks up a message. An unknown key is an error rather than an empty label:
// it aborts the build loudly instead of shipping a page with a hole in it.
// Catalogs are kept complete by TestCatalogsCoverTheSameKeys, so the error only
// fires on a key the templates made up.
func (c Catalog) T(key string) (string, error) {
	if msg, ok := c[key]; ok {
		return msg, nil
	}
	return "", fmt.Errorf("i18n: no message for %q", key)
}

// Msg is T for Go callers, where a key that does not resolve can only be a
// mistake in the code rather than in a template. The key stands in for the
// message so a payload built from one stays readable if that ever happens.
func (c Catalog) Msg(key string) string {
	if msg, ok := c[key]; ok {
		return msg
	}
	return key
}

// Of returns the catalog for a language, falling back to the default for the
// zero value — a Builder left unconfigured renders in English.
func Of(l Lang) Catalog {
	if c, ok := catalogs[l]; ok {
		return c
	}
	return catalogs[Default]
}

var catalogs = map[Lang]Catalog{
	En: en,
	Ja: ja,
}
