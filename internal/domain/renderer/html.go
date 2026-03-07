// internal/domain/renderer/html.go
package renderer

import (
	"fmt"
	"html"
	"strings"

	"gemini-bridge/internal/domain/model"
)

// RenderNodes はNodeスライスをHTML文字列に変換する。
func RenderNodes(nodes []model.Node) string {
	var b strings.Builder
	var inList bool

	for _, node := range nodes {
		// リストの開始/終了タグ管理
		isListItem := false
		if _, ok := node.(model.ListItem); ok {
			isListItem = true
		}

		if inList && !isListItem {
			b.WriteString("</ul>\n")
			inList = false
		}

		switch n := node.(type) {
		case model.Text:
			if n.Content == "" {
				b.WriteString("<br>\n")
			} else {
				fmt.Fprintf(&b, "<p>%s</p>\n", html.EscapeString(n.Content))
			}

		case model.Link:
			label := n.Label
			if label == "" {
				label = n.URL
			}
			escURL := html.EscapeString(n.URL)
			escLabel := html.EscapeString(label)
			// gemini:// リンクはHTTPS変換
			displayURL := convertGeminiURL(escURL)
			fmt.Fprintf(&b, "<p class=\"gemini-link\"><a href=\"%s\">%s</a></p>\n",
				displayURL, escLabel)

		case model.Heading:
			esc := html.EscapeString(n.Content)
			slug := slugify(n.Content)
			fmt.Fprintf(&b, "<h%d id=\"%s\">%s</h%d>\n",
				n.Level, slug, esc, n.Level)

		case model.ListItem:
			if !inList {
				b.WriteString("<ul>\n")
				inList = true
			}
			fmt.Fprintf(&b, "  <li>%s</li>\n", html.EscapeString(n.Content))

		case model.Quote:
			fmt.Fprintf(&b, "<blockquote><p>%s</p></blockquote>\n",
				html.EscapeString(n.Content))

		case model.Preformatted:
			b.WriteString("<pre")
			if n.AltText != "" {
				fmt.Fprintf(&b, " aria-label=\"%s\"", html.EscapeString(n.AltText))
			}
			b.WriteString("><code>")
			for i, line := range n.Lines {
				if i > 0 {
					b.WriteString("\n")
				}
				b.WriteString(html.EscapeString(line))
			}
			b.WriteString("</code></pre>\n")
		}
	}

	if inList {
		b.WriteString("</ul>\n")
	}

	return b.String()
}

// convertGeminiURL は gemini:// URLを https:// に変換する。
func convertGeminiURL(url string) string {
	if strings.HasPrefix(url, "gemini://") {
		return "https://" + url[9:]
	}
	return url
}

// slugify は見出しテキストからURL-safeなスラグを生成する。
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r > 127 {
			return r
		}
		if r == ' ' {
			return '-'
		}
		return -1
	}, s)
	return s
}
