// internal/application/config.go
package application

// BuildConfig はSSGビルドの設定を保持する。
type BuildConfig struct {
	SourceDir string // gemtextファイルのルートディレクトリ
	OutputDir string // ビルド出力ディレクトリ
	SiteTitle string
	SiteURL   string
	Author    string
}
