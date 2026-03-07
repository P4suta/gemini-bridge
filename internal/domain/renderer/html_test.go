package renderer

import (
	"strings"
	"testing"

	"gemini-bridge/internal/domain/model"
)

func TestRenderNodes_TextLine(t *testing.T) {
	nodes := []model.Node{
		model.Text{Content: "Hello, world!"},
	}
	result := RenderNodes(nodes)
	expected := "<p>Hello, world!</p>\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderNodes_EmptyTextLine(t *testing.T) {
	nodes := []model.Node{
		model.Text{Content: ""},
	}
	result := RenderNodes(nodes)
	expected := "<br>\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderNodes_Link(t *testing.T) {
	nodes := []model.Node{
		model.Link{URL: "https://example.com", Label: "Example"},
	}
	result := RenderNodes(nodes)
	if !strings.Contains(result, `<a href="https://example.com">Example</a>`) {
		t.Errorf("expected link HTML, got %q", result)
	}
}

func TestRenderNodes_Link_NoLabel(t *testing.T) {
	nodes := []model.Node{
		model.Link{URL: "https://example.com", Label: ""},
	}
	result := RenderNodes(nodes)
	if !strings.Contains(result, `<a href="https://example.com">https://example.com</a>`) {
		t.Errorf("expected URL as label, got %q", result)
	}
}

func TestRenderNodes_GeminiLink(t *testing.T) {
	nodes := []model.Node{
		model.Link{URL: "gemini://example.com/page", Label: "Gemini Page"},
	}
	result := RenderNodes(nodes)
	if !strings.Contains(result, `href="https://example.com/page"`) {
		t.Errorf("expected gemini:// to be converted to https://, got %q", result)
	}
}

func TestRenderNodes_Headings(t *testing.T) {
	nodes := []model.Node{
		model.Heading{Level: 1, Content: "Title"},
		model.Heading{Level: 2, Content: "Section"},
		model.Heading{Level: 3, Content: "Subsection"},
	}
	result := RenderNodes(nodes)
	if !strings.Contains(result, `<h1 id="title">Title</h1>`) {
		t.Errorf("expected h1, got %q", result)
	}
	if !strings.Contains(result, `<h2 id="section">Section</h2>`) {
		t.Errorf("expected h2, got %q", result)
	}
	if !strings.Contains(result, `<h3 id="subsection">Subsection</h3>`) {
		t.Errorf("expected h3, got %q", result)
	}
}

func TestRenderNodes_ListItems(t *testing.T) {
	nodes := []model.Node{
		model.ListItem{Content: "First"},
		model.ListItem{Content: "Second"},
	}
	result := RenderNodes(nodes)
	if !strings.Contains(result, "<ul>\n") {
		t.Errorf("expected opening <ul>, got %q", result)
	}
	if !strings.Contains(result, "</ul>\n") {
		t.Errorf("expected closing </ul>, got %q", result)
	}
	if !strings.Contains(result, "<li>First</li>") {
		t.Errorf("expected first list item, got %q", result)
	}
	if !strings.Contains(result, "<li>Second</li>") {
		t.Errorf("expected second list item, got %q", result)
	}
}

func TestRenderNodes_ListClosedByNonListItem(t *testing.T) {
	nodes := []model.Node{
		model.ListItem{Content: "Item"},
		model.Text{Content: "After list"},
	}
	result := RenderNodes(nodes)
	// The </ul> should appear before the <p>
	ulClose := strings.Index(result, "</ul>")
	pOpen := strings.Index(result, "<p>After list</p>")
	if ulClose < 0 || pOpen < 0 || ulClose > pOpen {
		t.Errorf("expected </ul> before <p>, got %q", result)
	}
}

func TestRenderNodes_Quote(t *testing.T) {
	nodes := []model.Node{
		model.Quote{Content: "Wise words"},
	}
	result := RenderNodes(nodes)
	expected := "<blockquote><p>Wise words</p></blockquote>\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderNodes_Preformatted(t *testing.T) {
	nodes := []model.Node{
		model.Preformatted{
			AltText: "go",
			Lines:   []string{"func main() {}", "  fmt.Println()"},
		},
	}
	result := RenderNodes(nodes)
	if !strings.Contains(result, `<pre aria-label="go"><code>`) {
		t.Errorf("expected pre with aria-label, got %q", result)
	}
	if !strings.Contains(result, "func main() {}") {
		t.Errorf("expected code content, got %q", result)
	}
}

func TestRenderNodes_Preformatted_NoAltText(t *testing.T) {
	nodes := []model.Node{
		model.Preformatted{
			Lines: []string{"plain text"},
		},
	}
	result := RenderNodes(nodes)
	if !strings.Contains(result, "<pre><code>") {
		t.Errorf("expected pre without aria-label, got %q", result)
	}
}

func TestRenderNodes_XSSEscape(t *testing.T) {
	nodes := []model.Node{
		model.Text{Content: "<script>alert('xss')</script>"},
		model.Heading{Level: 1, Content: "Title & <More>"},
		model.Link{URL: "https://example.com", Label: "Link & <Tag>"},
	}
	result := RenderNodes(nodes)
	if strings.Contains(result, "<script>") {
		t.Errorf("XSS not escaped in text: %q", result)
	}
	if strings.Contains(result, "& <More>") {
		t.Errorf("XSS not escaped in heading: %q", result)
	}
	if strings.Contains(result, "& <Tag>") {
		t.Errorf("XSS not escaped in link label: %q", result)
	}
}
