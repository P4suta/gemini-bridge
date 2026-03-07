// internal/domain/model/document.go
package model

import "time"

// Document はパース済みのgemtext文書を表す。
type Document struct {
	FrontMatter FrontMatter
	Nodes       []Node
}

// FrontMatter はgemtextファイル冒頭のメタデータ。
// プリフォーマットブロック内にYAML形式で記述する規約とする。
type FrontMatter struct {
	Title       string    `yaml:"title"`
	Date        time.Time `yaml:"date"`
	Slug        string    `yaml:"slug"`
	Tags        []string  `yaml:"tags"`
	Language    string    `yaml:"lang"`
	Description string    `yaml:"description"`
	Draft       bool      `yaml:"draft"`
}

// PostMeta はビルド後に生成される記事メタデータ。
// JSONシリアライズしてR2/D1に格納する。
type PostMeta struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Tags        []string  `json:"tags"`
	Language    string    `json:"lang"`
	Description string    `json:"description"`
	WordCount   int       `json:"wordCount"`
	GemtextHash string    `json:"gemtextHash"`
}

// Site はサイト全体の情報を保持する。
type Site struct {
	Title    string
	Subtitle string
	BaseURL  string
	Author   string
	Language string
	Posts    []PostMeta
}
