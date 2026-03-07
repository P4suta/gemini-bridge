// internal/domain/model/node.go
package model

// Node はgemtext行の抽象表現。
// unexportedメソッド sealed() により外部パッケージからの実装を防ぐ。
type Node interface {
	nodeType() NodeType
	sealed()
}

type NodeType int

const (
	NodeText NodeType = iota
	NodeLink
	NodeHeading
	NodeListItem
	NodeQuote
	NodePreformatted
)

// Text は通常のテキスト行。
type Text struct {
	Content string
}

func (Text) nodeType() NodeType { return NodeText }
func (Text) sealed()            {}

// Link はリンク行（=> URL label）。
type Link struct {
	URL   string
	Label string // 空の場合はURLをラベルとして表示
}

func (Link) nodeType() NodeType { return NodeLink }
func (Link) sealed()            {}

// Heading は見出し行（#, ##, ###）。
type Heading struct {
	Level   int // 1, 2, or 3
	Content string
}

func (Heading) nodeType() NodeType { return NodeHeading }
func (Heading) sealed()            {}

// ListItem はリスト項目行（* item）。
type ListItem struct {
	Content string
}

func (ListItem) nodeType() NodeType { return NodeListItem }
func (ListItem) sealed()            {}

// Quote は引用行（> text）。
type Quote struct {
	Content string
}

func (Quote) nodeType() NodeType { return NodeQuote }
func (Quote) sealed()            {}

// Preformatted はプリフォーマットブロック（```で囲まれた範囲）。
// gemtextでは複数行にまたがる唯一のブロック型。
type Preformatted struct {
	AltText string   // ``` に続くalt text（言語ヒント等）
	Lines   []string // ブロック内の各行
}

func (Preformatted) nodeType() NodeType { return NodePreformatted }
func (Preformatted) sealed()            {}
