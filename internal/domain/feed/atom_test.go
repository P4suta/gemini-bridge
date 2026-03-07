package feed

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"gemini-bridge/internal/domain/model"
)

func TestGenerateAtom_ValidXML(t *testing.T) {
	site := model.Site{
		Title:   "Test Blog",
		BaseURL: "https://example.com",
		Author:  "Test Author",
		Posts: []model.PostMeta{
			{
				Slug:        "hello-world",
				Title:       "Hello World",
				Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Description: "A test post",
			},
		},
	}

	result, err := GenerateAtom(site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify valid XML
	var feed atomFeed
	if err := xml.Unmarshal([]byte(result), &feed); err != nil {
		t.Fatalf("generated XML is invalid: %v", err)
	}

	if feed.Title != "Test Blog" {
		t.Errorf("expected title 'Test Blog', got '%s'", feed.Title)
	}

	if len(feed.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(feed.Entries))
	}

	if feed.Entries[0].Title != "Hello World" {
		t.Errorf("expected entry title 'Hello World', got '%s'", feed.Entries[0].Title)
	}
}

func TestGenerateAtom_EntriesOrder(t *testing.T) {
	site := model.Site{
		Title:   "Test Blog",
		BaseURL: "https://example.com",
		Author:  "Test Author",
		Posts: []model.PostMeta{
			{Slug: "newer", Title: "Newer Post", Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
			{Slug: "older", Title: "Older Post", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	result, err := GenerateAtom(site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var feed atomFeed
	if err := xml.Unmarshal([]byte(result), &feed); err != nil {
		t.Fatalf("generated XML is invalid: %v", err)
	}

	if len(feed.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(feed.Entries))
	}

	// Order should be preserved (caller is responsible for sorting)
	if feed.Entries[0].Title != "Newer Post" {
		t.Errorf("expected first entry 'Newer Post', got '%s'", feed.Entries[0].Title)
	}
	if feed.Entries[1].Title != "Older Post" {
		t.Errorf("expected second entry 'Older Post', got '%s'", feed.Entries[1].Title)
	}
}

func TestGenerateAtom_EmptyFeed(t *testing.T) {
	site := model.Site{
		Title:   "Empty Blog",
		BaseURL: "https://example.com",
		Author:  "Test Author",
		Posts:   nil,
	}

	result, err := GenerateAtom(site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still be valid XML
	var feed atomFeed
	if err := xml.Unmarshal([]byte(result), &feed); err != nil {
		t.Fatalf("generated XML is invalid: %v", err)
	}

	if len(feed.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(feed.Entries))
	}
}

func TestGenerateAtom_XMLHeader(t *testing.T) {
	site := model.Site{
		Title:   "Test Blog",
		BaseURL: "https://example.com",
		Author:  "Test Author",
	}

	result, err := GenerateAtom(site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(result, "<?xml") {
		t.Errorf("expected XML header, got %q", result[:50])
	}
}

func TestGenerateAtom_GemtextAlternateLink(t *testing.T) {
	site := model.Site{
		Title:   "Test Blog",
		BaseURL: "https://example.com",
		Author:  "Test Author",
		Posts: []model.PostMeta{
			{Slug: "test", Title: "Test", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	result, err := GenerateAtom(site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "text/gemini") {
		t.Errorf("expected gemtext alternate link, got %q", result)
	}
}
