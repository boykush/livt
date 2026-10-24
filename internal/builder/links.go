package builder

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
)

// bodyLinkURI reads a link destination written in a body as a point of the livt
// repository. ours is false for a destination that points elsewhere — another
// site, an anchor on the same page — which the body keeps as written. A
// destination that is ours but names nothing parses to ok false.
func (b *Builder) bodyLinkURI(dest, sourceDir string) (p uri.Parsed, ok, ours bool) {
	if strings.HasPrefix(dest, "livt://") {
		p, ok = uri.Parse(dest)
		return p, ok, true
	}
	u, err := url.Parse(dest)
	if err != nil || u.Scheme != "" || u.Host != "" || u.Path == "" {
		return uri.Parsed{}, false, false
	}
	// A leading slash is the livt repository's root, the way GitHub reads it,
	// so the same link works on the forge and on the site.
	full := filepath.Join(sourceDir, filepath.FromSlash(u.Path))
	if strings.HasPrefix(u.Path, "/") {
		full = filepath.Join(b.root(), filepath.FromSlash(u.Path))
	}
	p, ok = b.pathURI(full)
	return p, ok, true
}

// pathURI is the livt URI of the file at path, if it is one of the files a livt
// URI can name.
func (b *Builder) pathURI(path string) (uri.Parsed, bool) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return uri.Parsed{}, false
	}
	for _, d := range []struct {
		dir, ext string
		uri      func(ref string) string
	}{
		{b.UbiquitousDir, ".md", func(ref string) string { return uri.Term("", ref) }},
		{b.PoliciesDir, ".md", func(ref string) string { return uri.Policy("", ref) }},
		{b.StoriesDir, ".md", uri.Story},
		{b.OpportunitiesDir, ".md", uri.Opportunity},
		{b.MappingsDir, ".yaml", uri.Mapping},
		{b.USMDir, ".yaml", uri.StoryMap},
		{b.CanvasesDir, ".yaml", uri.OpportunityCanvas},
	} {
		if d.dir == "" {
			continue
		}
		dir, err := filepath.Abs(d.dir)
		if err != nil {
			continue
		}
		rel, err := filepath.Rel(dir, abs)
		if err != nil || !filepath.IsLocal(rel) || !strings.HasSuffix(rel, d.ext) {
			continue
		}
		return uri.Parse(d.uri(filepath.ToSlash(strings.TrimSuffix(rel, d.ext))))
	}
	return uri.Parsed{}, false
}

// linkName is what the point a livt URI names is called, and whether the livt
// repository holds it at all: a link to a point it does not hold leads nowhere.
func (b *Builder) linkName(p uri.Parsed) (string, bool) {
	switch p.Kind {
	case uri.KindTerm:
		t, err := parser.ParseTerm(b.UbiquitousDir, uri.TermRef(p.TermCtx, p.TermKey))
		if err != nil {
			return "", false
		}
		return orKey(t.Name, t.Key), true
	case uri.KindPolicy:
		pol, err := parser.ParsePolicy(b.PoliciesDir, uri.PolicyRef(p.PolicyCtx, p.PolicyKey))
		if err != nil {
			return "", false
		}
		return orKey(pol.Name, pol.Key), true
	case uri.KindStory:
		s, err := parser.ParseStory(filepath.Join(b.StoriesDir, p.StoryKey+".md"))
		if err != nil {
			return "", false
		}
		return s.DisplayName(), true
	case uri.KindOpportunity:
		o, err := parser.ParseOpportunity(filepath.Join(b.OpportunitiesDir, p.OpportunityKey+".md"))
		if err != nil {
			return "", false
		}
		return o.DisplayName(), true
	case uri.KindOpportunityCanvas:
		if _, err := os.Stat(filepath.Join(b.CanvasesDir, p.OpportunityKey+".yaml")); err != nil {
			return "", false
		}
		if o, err := parser.ParseOpportunity(filepath.Join(b.OpportunitiesDir, p.OpportunityKey+".md")); err == nil {
			return o.DisplayName(), true
		}
		return p.OpportunityKey, true
	case uri.KindStoryMap:
		sm, err := parser.ParseStoryMap(filepath.Join(b.USMDir, p.OpportunityKey+".yaml"))
		if err != nil {
			return "", false
		}
		return sm.DisplayName(), true
	case uri.KindMapping, uri.KindRule, uri.KindExample, uri.KindQuestion:
		return b.mappingLinkName(p)
	}
	return "", false
}

// mappingLinkName names a mapping or a sticky on it. A retired sticky still
// resolves: its id stays taken, and the link lands on the board it was on.
func (b *Builder) mappingLinkName(p uri.Parsed) (string, bool) {
	em, err := parser.ParseExampleMapping(filepath.Join(b.MappingsDir, p.StoryKey+".yaml"))
	if err != nil {
		return "", false
	}
	switch p.Kind {
	case uri.KindMapping:
		story, err := parser.ParseStory(filepath.Join(b.StoriesDir, p.StoryKey+".md"))
		if err != nil {
			story = &domain.Story{Key: em.StoryKey}
		}
		return em.DisplayName(story), true
	case uri.KindQuestion:
		for _, q := range em.Questions {
			if q.ID == p.QuestionID {
				return q.Text, true
			}
		}
		return "", false
	}
	for _, r := range em.Rules {
		if r.ID != p.RuleID {
			continue
		}
		if p.Kind == uri.KindRule {
			return r.Name, true
		}
		for _, ex := range r.Examples {
			if ex.ID == p.ExampleID {
				return ex.Name, true
			}
		}
	}
	return "", false
}

func orKey(name, key string) string {
	if name != "" {
		return name
	}
	return key
}
