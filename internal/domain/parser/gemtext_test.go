package parser

import (
	"strings"
	"testing"

	"gemini-bridge/internal/domain/model"
)

func TestParse_LinkLine(t *testing.T) {
	input := "=> https://example.com Example Site\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}

	link, ok := nodes[0].(model.Link)
	if !ok {
		t.Fatalf("expected Link node, got %T", nodes[0])
	}

	if link.URL != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got '%s'", link.URL)
	}

	if link.Label != "Example Site" {
		t.Errorf("expected label 'Example Site', got '%s'", link.Label)
	}
}

func TestParse_LinkLine_TabSeparated(t *testing.T) {
	input := "=> https://example.com\tExample Site\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}

	link, ok := nodes[0].(model.Link)
	if !ok {
		t.Fatalf("expected Link node, got %T", nodes[0])
	}

	if link.URL != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got '%s'", link.URL)
	}

	if link.Label != "Example Site" {
		t.Errorf("expected label 'Example Site', got '%s'", link.Label)
	}
}

func TestParse_LinkLine_NoLabel(t *testing.T) {
	input := "=> https://example.com\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	link, ok := nodes[0].(model.Link)
	if !ok {
		t.Fatalf("expected Link node, got %T", nodes[0])
	}

	if link.URL != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got '%s'", link.URL)
	}

	if link.Label != "" {
		t.Errorf("expected empty label, got '%s'", link.Label)
	}
}

func TestParse_PreformattedBlock(t *testing.T) {
	input := "```go\nfunc main() {}\nfmt.Println()\n```\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}

	pre, ok := nodes[0].(model.Preformatted)
	if !ok {
		t.Fatalf("expected Preformatted node, got %T", nodes[0])
	}

	if pre.AltText != "go" {
		t.Errorf("expected alt text 'go', got '%s'", pre.AltText)
	}

	if len(pre.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(pre.Lines))
	}
}

func TestParse_PreformattedBlock_Unclosed(t *testing.T) {
	input := "```\nline1\nline2\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}

	pre, ok := nodes[0].(model.Preformatted)
	if !ok {
		t.Fatalf("expected Preformatted node, got %T", nodes[0])
	}

	if len(pre.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(pre.Lines))
	}
}

func TestParse_HeadingLevels(t *testing.T) {
	input := "# H1\n## H2\n### H3\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}

	for i, expected := range []int{1, 2, 3} {
		h, ok := nodes[i].(model.Heading)
		if !ok {
			t.Errorf("node %d: expected Heading, got %T", i, nodes[i])
		}
		if h.Level != expected {
			t.Errorf("node %d: expected level %d, got %d", i, expected, h.Level)
		}
	}
}

func TestParse_ListItems(t *testing.T) {
	input := "* Item one\n* Item two\n* Item three\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}

	for i, expected := range []string{"Item one", "Item two", "Item three"} {
		li, ok := nodes[i].(model.ListItem)
		if !ok {
			t.Errorf("node %d: expected ListItem, got %T", i, nodes[i])
		}
		if li.Content != expected {
			t.Errorf("node %d: expected '%s', got '%s'", i, expected, li.Content)
		}
	}
}

func TestParse_QuoteLine(t *testing.T) {
	input := "> This is a quote\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}

	q, ok := nodes[0].(model.Quote)
	if !ok {
		t.Fatalf("expected Quote node, got %T", nodes[0])
	}

	if q.Content != "This is a quote" {
		t.Errorf("expected 'This is a quote', got '%s'", q.Content)
	}
}

func TestParse_TextLine(t *testing.T) {
	input := "Just a normal line\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}

	txt, ok := nodes[0].(model.Text)
	if !ok {
		t.Fatalf("expected Text node, got %T", nodes[0])
	}

	if txt.Content != "Just a normal line" {
		t.Errorf("expected 'Just a normal line', got '%s'", txt.Content)
	}
}

func TestParse_EmptyLine(t *testing.T) {
	input := "\n"
	nodes, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}

	txt, ok := nodes[0].(model.Text)
	if !ok {
		t.Fatalf("expected Text node, got %T", nodes[0])
	}

	if txt.Content != "" {
		t.Errorf("expected empty content, got '%s'", txt.Content)
	}
}

func TestParse_MixedContent(t *testing.T) {
	input := `# Welcome
This is a paragraph.
=> https://gemini.circumlunar.space Gemini Protocol
* Item one
* Item two
> A wise quote
`
	nodes, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 6 {
		t.Fatalf("expected 6 nodes, got %d", len(nodes))
	}

	// 型アサーションで各ノードの型を検証（unexportedメソッドに依存しない）
	if _, ok := nodes[0].(model.Heading); !ok {
		t.Errorf("node 0: expected Heading, got %T", nodes[0])
	}
	if _, ok := nodes[1].(model.Text); !ok {
		t.Errorf("node 1: expected Text, got %T", nodes[1])
	}
	if _, ok := nodes[2].(model.Link); !ok {
		t.Errorf("node 2: expected Link, got %T", nodes[2])
	}
	if _, ok := nodes[3].(model.ListItem); !ok {
		t.Errorf("node 3: expected ListItem, got %T", nodes[3])
	}
	if _, ok := nodes[4].(model.ListItem); !ok {
		t.Errorf("node 4: expected ListItem, got %T", nodes[4])
	}
	if _, ok := nodes[5].(model.Quote); !ok {
		t.Errorf("node 5: expected Quote, got %T", nodes[5])
	}
}
