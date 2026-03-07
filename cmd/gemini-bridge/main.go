// cmd/gemini-bridge/main.go
package main

import (
	"flag"
	"log"
	"os"

	"gemini-bridge/internal/application"
	"gemini-bridge/internal/infrastructure"
)

func main() {
	// CLI引数の定義
	sourceDir := flag.String("source", "content", "Source directory containing .gmi files")
	outputDir := flag.String("output", "dist", "Output directory for generated files")
	siteTitle := flag.String("title", "gemini-bridge", "Site title")
	siteURL := flag.String("url", "https://blog.sakashita.dev", "Site base URL")
	author := flag.String("author", "坂下 康信", "Site author name")
	flag.Parse()

	// ソースディレクトリの存在確認
	if _, err := os.Stat(*sourceDir); os.IsNotExist(err) {
		log.Fatalf("Source directory does not exist: %s", *sourceDir)
	}

	// 手動DI: ポート層インターフェースにインフラ層実装を注入
	writer := infrastructure.NewFileSystemWriter(*outputDir)
	metaStore := infrastructure.NewJsonMetadataStore(*outputDir)

	config := application.BuildConfig{
		SourceDir: *sourceDir,
		OutputDir: *outputDir,
		SiteTitle: *siteTitle,
		SiteURL:   *siteURL,
		Author:    *author,
	}

	pipeline := application.NewBuildPipeline(writer, metaStore, config)

	log.Printf("Starting build: source=%s output=%s", *sourceDir, *outputDir)
	if err := pipeline.Execute(); err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	log.Println("Build succeeded.")
}
