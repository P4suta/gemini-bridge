// internal/application/pipeline.go
package application

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"gemini-bridge/internal/domain/feed"
	"gemini-bridge/internal/domain/model"
	"gemini-bridge/internal/domain/parser"
	"gemini-bridge/internal/domain/renderer"
	"gemini-bridge/internal/port"
)

// BuildPipeline はgemtextからの静的サイト生成を統括する。
type BuildPipeline struct {
	writer    port.ContentWriter
	metaStore port.MetadataStore
	config    BuildConfig
}

func NewBuildPipeline(
	writer port.ContentWriter,
	metaStore port.MetadataStore,
	config BuildConfig,
) *BuildPipeline {
	return &BuildPipeline{
		writer:    writer,
		metaStore: metaStore,
		config:    config,
	}
}

// Execute はビルドパイプライン全体を実行する。
func (p *BuildPipeline) Execute() error {
	// 1. 全gemtextファイルを収集
	gmiFiles, err := p.collectGemtextFiles()
	if err != nil {
		return fmt.Errorf("failed to collect gemtext files: %w", err)
	}

	log.Printf("Found %d gemtext files", len(gmiFiles))

	// 2. 各ファイルをパース・変換
	var posts []model.PostMeta
	for _, path := range gmiFiles {
		meta, err := p.processFile(path)
		if err != nil {
			return fmt.Errorf("failed to process %s: %w", path, err)
		}
		if meta.Slug != "" {
			posts = append(posts, meta)
		}
	}

	// 3. 日付降順にソート
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	// 4. インデックスページ生成
	if err := p.generateIndex(posts); err != nil {
		return fmt.Errorf("failed to generate index: %w", err)
	}

	// 5. Atomフィード生成
	site := model.Site{
		Title:    p.config.SiteTitle,
		BaseURL:  p.config.SiteURL,
		Author:   p.config.Author,
		Language: "ja",
		Posts:    posts,
	}
	if err := p.generateFeed(site); err != nil {
		return fmt.Errorf("failed to generate feed: %w", err)
	}

	// 6. メタデータインデックス保存
	if err := p.metaStore.SaveSiteIndex(posts); err != nil {
		return fmt.Errorf("failed to save site index: %w", err)
	}

	log.Printf("Build complete: %d posts generated", len(posts))
	return nil
}

func (p *BuildPipeline) processFile(path string) (model.PostMeta, error) {
	// ファイル読み込み
	raw, err := os.ReadFile(path)
	if err != nil {
		return model.PostMeta{}, err
	}

	// パース
	nodes, err := parser.Parse(strings.NewReader(string(raw)))
	if err != nil {
		return model.PostMeta{}, err
	}

	// Front Matter抽出
	fm, contentNodes := parser.ExtractFrontMatter(nodes)

	// ドラフトはスキップ
	if fm.Draft {
		log.Printf("Skipping draft: %s", path)
		return model.PostMeta{}, nil
	}

	// スラグの導出
	slug := fm.Slug
	if slug == "" {
		base := filepath.Base(path)
		slug = strings.TrimSuffix(base, filepath.Ext(base))
	}

	// HTML変換
	htmlContent := renderer.RenderNodes(contentNodes)

	// 全文のワードカウント
	var wordCount int
	for _, node := range contentNodes {
		switch n := node.(type) {
		case model.Text:
			wordCount += utf8.RuneCountInString(n.Content)
		case model.Heading:
			wordCount += utf8.RuneCountInString(n.Content)
		case model.ListItem:
			wordCount += utf8.RuneCountInString(n.Content)
		case model.Quote:
			wordCount += utf8.RuneCountInString(n.Content)
		}
	}

	// ハッシュ計算（変更検知用）
	hash := sha256.Sum256(raw)
	hashStr := fmt.Sprintf("%x", hash[:8])

	// HTML出力
	htmlPath := fmt.Sprintf("posts/%s/index.html", slug)
	if err := p.writer.WriteBytes(htmlPath, []byte(htmlContent)); err != nil {
		return model.PostMeta{}, err
	}

	// 生gemtext出力（Content Negotiation用）
	gmiPath := fmt.Sprintf("posts/%s/index.gmi", slug)
	if err := p.writer.WriteBytes(gmiPath, raw); err != nil {
		return model.PostMeta{}, err
	}

	// メタデータ
	meta := model.PostMeta{
		Slug:        slug,
		Title:       fm.Title,
		Date:        fm.Date,
		Tags:        fm.Tags,
		Language:    fm.Language,
		Description: fm.Description,
		WordCount:   wordCount,
		GemtextHash: hashStr,
	}

	// メタデータJSON出力
	if err := p.metaStore.SavePostMeta(meta); err != nil {
		return model.PostMeta{}, err
	}

	log.Printf("Processed: %s -> %s", path, slug)
	return meta, nil
}

func (p *BuildPipeline) collectGemtextFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(p.config.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".gmi") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (p *BuildPipeline) generateIndex(posts []model.PostMeta) error {
	const indexTmpl = `<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.SiteTitle}}</title>
  <link rel="alternate" type="application/atom+xml" href="/feed/atom.xml">
  <link rel="stylesheet" href="/assets/css/style.css">
</head>
<body>
  <header>
    <h1>{{.SiteTitle}}</h1>
    <nav>
      <a href="/">Home</a>
      <a href="/feed/atom.xml">Feed</a>
    </nav>
  </header>
  <main>
    <section class="post-list">
      {{range .Posts}}
      <article>
        <h2><a href="/posts/{{.Slug}}/">{{.Title}}</a></h2>
        <time datetime="{{.Date.Format "2006-01-02"}}">{{.Date.Format "2006年1月2日"}}</time>
        {{if .Description}}<p>{{.Description}}</p>{{end}}
        {{if .Tags}}
        <ul class="tags">
          {{range .Tags}}<li>{{.}}</li>{{end}}
        </ul>
        {{end}}
      </article>
      {{end}}
    </section>
  </main>
  <footer>
    <p>Powered by gemini-bridge</p>
  </footer>
</body>
</html>`

	tmpl, err := template.New("index").Parse(indexTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse index template: %w", err)
	}

	data := struct {
		SiteTitle string
		Posts     []model.PostMeta
	}{
		SiteTitle: p.config.SiteTitle,
		Posts:     posts,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute index template: %w", err)
	}

	return p.writer.WriteBytes("index.html", []byte(buf.String()))
}

func (p *BuildPipeline) generateFeed(site model.Site) error {
	atomXML, err := feed.GenerateAtom(site)
	if err != nil {
		return err
	}
	return p.writer.WriteBytes("feed/atom.xml", []byte(atomXML))
}
