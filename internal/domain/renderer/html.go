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
			// gemini:// リンクはHTTPS変換（エスケープ前に実施）
			safeURL := sanitizeURL(convertGeminiURL(n.URL))
			escURL := html.EscapeString(safeURL)
			escLabel := html.EscapeString(label)
			fmt.Fprintf(&b, "<p class=\"gemini-link\"><a href=\"%s\">%s</a></p>\n",
				escURL, escLabel)

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

// sanitizeURL は安全でないURLスキームを除去する。
func sanitizeURL(url string) string {
	lower := strings.ToLower(strings.TrimSpace(url))
	// 許可するスキーム
	if strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "gemini://") ||
		strings.HasPrefix(lower, "mailto:") ||
		strings.HasPrefix(lower, "/") {
		return url
	}
	// スキームが無い相対パスは許可
	if !strings.Contains(strings.SplitN(lower, "/", 2)[0], ":") {
		return url
	}
	// javascript: 等の危険なスキームは除去
	return ""
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
