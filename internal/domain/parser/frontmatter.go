// internal/domain/parser/frontmatter.go
package parser

import (
	"strings"
	"time"

	"gemini-bridge/internal/domain/model"
)

// ExtractFrontMatter はNodeスライスの先頭からFront Matterを抽出する。
// 先頭のPreformattedノードのAltTextが "yaml" の場合にFront Matterとして解釈する。
// 戻り値は、FrontMatter構造体と、Front Matterを除いた残りのNodeスライス。
func ExtractFrontMatter(nodes []model.Node) (model.FrontMatter, []model.Node) {
	if len(nodes) == 0 {
		return model.FrontMatter{}, nodes
	}

	pre, ok := nodes[0].(model.Preformatted)
	if !ok || !strings.EqualFold(pre.AltText, "yaml") {
		return model.FrontMatter{}, nodes
	}

	fm := parseFrontMatterLines(pre.Lines)
	return fm, nodes[1:]
}

// parseFrontMatterLines は key: value 形式の行をパースする。
// 外部YAMLライブラリに依存せず、サポートするフィールドに限定した軽量パーサー。
func parseFrontMatterLines(lines []string) model.FrontMatter {
	var fm model.FrontMatter

	for _, line := range lines {
		key, value, ok := parseKeyValue(line)
		if !ok {
			continue
		}

		switch key {
		case "title":
			fm.Title = value
		case "date":
			if t, err := time.Parse("2006-01-02", value); err == nil {
				fm.Date = t
			}
		case "slug":
			fm.Slug = value
		case "tags":
			for _, tag := range strings.Split(value, ",") {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					fm.Tags = append(fm.Tags, tag)
				}
			}
		case "lang":
			fm.Language = value
		case "description":
			fm.Description = value
		case "draft":
			fm.Draft = value == "true"
		}
	}

	// デフォルト値の設定
	if fm.Language == "" {
		fm.Language = "ja"
	}

	return fm
}

func parseKeyValue(line string) (key, value string, ok bool) {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	return key, value, key != ""
}
