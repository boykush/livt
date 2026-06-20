package builder

import (
	"testing"

	"github.com/boykush/livt/internal/domain"
)

func TestMetaFieldViewsAutoLinksURLs(t *testing.T) {
	views := metaFieldViews([]domain.MetaField{
		{Key: "status", Value: "draft"},
		{Key: "ticketUrl", Value: "https://example.com/1"},
	})
	if len(views) != 2 {
		t.Fatalf("got %d views, want 2", len(views))
	}
	if views[0].Href != "" {
		t.Fatalf("plain value should not link, got %q", views[0].Href)
	}
	if views[1].Href != "https://example.com/1" {
		t.Fatalf("URL value should link, got %q", views[1].Href)
	}
}

func TestURLMetaFieldViewsKeepsOnlyURLs(t *testing.T) {
	views := urlMetaFieldViews([]domain.MetaField{
		{Key: "status", Value: "draft"},
		{Key: "ticketUrl", Value: "http://example.com"},
	})
	if len(views) != 1 || views[0].Key != "ticketUrl" {
		t.Fatalf("got %+v, want only ticketUrl", views)
	}
}
