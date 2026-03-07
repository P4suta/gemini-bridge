// internal/domain/parser/gemtext.go
package parser

import (
	"bufio"
	"io"
	"strings"

	"gemini-bridge/internal/domain/model"
)

// Parse はgemtextのReaderを受け取り、Nodeのスライスを返す。
// Front Matterの解析は呼び出し側の責務とし、本関数は純粋なgemtextパースに専念する。
func Parse(r io.Reader) ([]model.Node, error) {
	var nodes []model.Node
	scanner := bufio.NewScanner(r)

	var inPreformatted bool
	var preBlock model.Preformatted

	for scanner.Scan() {
		line := scanner.Text()

		// プリフォーマットトグル判定
		if strings.HasPrefix(line, "```") {
			if inPreformatted {
				// ブロック終了: 蓄積した行をノードとして追加
				nodes = append(nodes, preBlock)
				preBlock = model.Preformatted{}
				inPreformatted = false
			} else {
				// ブロック開始: alt textを取得
				preBlock.AltText = strings.TrimSpace(line[3:])
				inPreformatted = true
			}
			continue
		}

		// プリフォーマットモード中は行をそのまま蓄積
		if inPreformatted {
			preBlock.Lines = append(preBlock.Lines, line)
			continue
		}

		// 通常モード: 行頭の文字列パターンで分岐
		node := parseLine(line)
		nodes = append(nodes, node)
	}

	// ファイル末尾で閉じられていないプリフォーマットブロックの処理
	if inPreformatted {
		nodes = append(nodes, preBlock)
	}

	return nodes, scanner.Err()
}

// parseLine は通常モードの1行をパースして適切なNodeを返す。
func parseLine(line string) model.Node {
	// リンク行: => URL optional-label
	if strings.HasPrefix(line, "=>") {
		return parseLink(line[2:])
	}

	// 見出し行: ###, ##, # の順にチェック（長い方を先に）
	if strings.HasPrefix(line, "###") {
		return model.Heading{Level: 3, Content: strings.TrimSpace(line[3:])}
	}
	if strings.HasPrefix(line, "##") {
		return model.Heading{Level: 2, Content: strings.TrimSpace(line[2:])}
	}
	if strings.HasPrefix(line, "#") {
		return model.Heading{Level: 1, Content: strings.TrimSpace(line[1:])}
	}

	// リストアイテム: * text
	if strings.HasPrefix(line, "* ") {
		return model.ListItem{Content: strings.TrimSpace(line[2:])}
	}

	// 引用行: > text
	if strings.HasPrefix(line, ">") {
		return model.Quote{Content: strings.TrimSpace(line[1:])}
	}

	// その他: テキスト行
	return model.Text{Content: line}
}

// parseLink はリンク行の URL と label を分離する。
func parseLink(raw string) model.Link {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return model.Link{}
	}

	// 最初の空白文字（スペースまたはタブ）でURLとラベルを分割
	idx := strings.IndexAny(trimmed, " \t")
	if idx < 0 {
		return model.Link{URL: trimmed}
	}

	url := trimmed[:idx]
	label := strings.TrimSpace(trimmed[idx+1:])

	return model.Link{URL: url, Label: label}
}
