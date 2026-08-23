package i18n

import (
	"sort"
	"testing"
)

// TestCatalogsCoverTheSameKeys is what lets Catalog.T treat a missing key as an
// error: every language answers for every message, so a page can never render
// with a hole in it. A message added to en without a translation fails here,
// before it can fail a user's build.
func TestCatalogsCoverTheSameKeys(t *testing.T) {
	for lang, c := range catalogs {
		if lang == Default {
			continue
		}
		for _, key := range keys(en) {
			if _, ok := c[key]; !ok {
				t.Errorf("%s is missing a translation for %q", lang, key)
			}
		}
		for _, key := range keys(c) {
			if _, ok := en[key]; !ok {
				t.Errorf("%s translates %q, which no longer exists in en", lang, key)
			}
		}
	}
}

// Langs drives both the "unknown lang" error message and the per-language
// template sets, so a language with a catalog but no entry there would be
// unreachable — and one listed without a catalog would fall back to English.
func TestLangsMatchTheCatalogs(t *testing.T) {
	if len(Langs) != len(catalogs) {
		t.Fatalf("Langs has %d entries, catalogs has %d", len(Langs), len(catalogs))
	}
	for _, l := range Langs {
		if !l.Valid() {
			t.Errorf("Langs lists %q, which has no catalog", l)
		}
	}
}

func TestOfFallsBackToDefault(t *testing.T) {
	if got := Of(""); got == nil {
		t.Fatal("Of(\"\") returned no catalog")
	}
	msg, err := Of("").T("nav.stories")
	if err != nil {
		t.Fatal(err)
	}
	if want := en["nav.stories"]; msg != want {
		t.Fatalf("Of(\"\") = %q, want the default catalog's %q", msg, want)
	}
}

func TestTReportsAnUnknownKey(t *testing.T) {
	if _, err := Of(Ja).T("nav.nope"); err == nil {
		t.Fatal("expected an error for a key no catalog holds")
	}
}

func keys(c Catalog) []string {
	out := make([]string, 0, len(c))
	for k := range c {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
