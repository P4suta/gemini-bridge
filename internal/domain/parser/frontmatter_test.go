package parser

import (
	"strings"
	"testing"
	"time"

	"gemini-bridge/internal/domain/model"
)

func TestExtractFrontMatter_ValidYAML(t *testing.T) {
	input := "```yaml\ntitle: Hello World\ndate: 2024-01-15\ntags: tech, gemini\nlang: ja\n```\n# Hello World\nSome content here.\n"
	nodes, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm, remaining := ExtractFrontMatter(nodes)

	if fm.Title != "Hello World" {
		t.Errorf("expected title 'Hello World', got '%s'", fm.Title)
	}

	expectedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if !fm.Date.Equal(expectedDate) {
		t.Errorf("expected date %v, got %v", expectedDate, fm.Date)
	}

	if len(fm.Tags) != 2 || fm.Tags[0] != "tech" || fm.Tags[1] != "gemini" {
		t.Errorf("expected tags [tech, gemini], got %v", fm.Tags)
	}

	if fm.Language != "ja" {
		t.Errorf("expected language 'ja', got '%s'", fm.Language)
	}

	// Remaining nodes should be the heading and text
	if len(remaining) != 2 {
		t.Fatalf("expected 2 remaining nodes, got %d", len(remaining))
	}

	if _, ok := remaining[0].(model.Heading); !ok {
		t.Errorf("expected remaining[0] to be Heading, got %T", remaining[0])
	}
}

func TestExtractFrontMatter_DefaultLanguage(t *testing.T) {
	input := "```yaml\ntitle: Test\ndate: 2024-01-01\n```\n"
	nodes, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm, _ := ExtractFrontMatter(nodes)

	if fm.Language != "ja" {
		t.Errorf("expected default language 'ja', got '%s'", fm.Language)
	}
}

func TestExtractFrontMatter_NonYAMLBlock(t *testing.T) {
	input := "```go\nfunc main() {}\n```\n"
	nodes, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm, remaining := ExtractFrontMatter(nodes)

	if fm.Title != "" {
		t.Errorf("expected empty title, got '%s'", fm.Title)
	}

	// All nodes should remain unchanged
	if len(remaining) != 1 {
		t.Fatalf("expected 1 remaining node, got %d", len(remaining))
	}
}

func TestExtractFrontMatter_EmptyNodes(t *testing.T) {
	fm, remaining := ExtractFrontMatter(nil)

	if fm.Title != "" {
		t.Errorf("expected empty title, got '%s'", fm.Title)
	}

	if len(remaining) != 0 {
		t.Errorf("expected 0 remaining nodes, got %d", len(remaining))
	}
}

func TestExtractFrontMatter_Draft(t *testing.T) {
	input := "```yaml\ntitle: Draft Post\ndraft: true\n```\n"
	nodes, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm, _ := ExtractFrontMatter(nodes)

	if !fm.Draft {
		t.Error("expected draft to be true")
	}
}

func TestExtractFrontMatter_AllFields(t *testing.T) {
	input := "```yaml\ntitle: Full Post\ndate: 2024-06-15\nslug: custom-slug\ntags: go, testing, gemini\nlang: en\ndescription: A full test post\ndraft: false\n```\n"
	nodes, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fm, _ := ExtractFrontMatter(nodes)

	if fm.Title != "Full Post" {
		t.Errorf("expected title 'Full Post', got '%s'", fm.Title)
	}
	if fm.Slug != "custom-slug" {
		t.Errorf("expected slug 'custom-slug', got '%s'", fm.Slug)
	}
	if len(fm.Tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(fm.Tags))
	}
	if fm.Language != "en" {
		t.Errorf("expected language 'en', got '%s'", fm.Language)
	}
	if fm.Description != "A full test post" {
		t.Errorf("expected description 'A full test post', got '%s'", fm.Description)
	}
	if fm.Draft {
		t.Error("expected draft to be false")
	}
}
