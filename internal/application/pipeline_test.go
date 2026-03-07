package application

import (
	"io"
	"testing"

	"gemini-bridge/internal/domain/model"
	"gemini-bridge/internal/infrastructure"
)

// memoryWriter implements port.ContentWriter for testing.
type memoryWriter struct {
	files map[string][]byte
}

func newMemoryWriter() *memoryWriter {
	return &memoryWriter{files: make(map[string][]byte)}
}

func (w *memoryWriter) Write(path string, content io.Reader) error {
	data, err := io.ReadAll(content)
	if err != nil {
		return err
	}
	w.files[path] = data
	return nil
}

func (w *memoryWriter) WriteBytes(path string, data []byte) error {
	w.files[path] = data
	return nil
}

// memoryMetaStore implements port.MetadataStore for testing.
type memoryMetaStore struct {
	posts     []model.PostMeta
	siteIndex []model.PostMeta
}

func (s *memoryMetaStore) SavePostMeta(meta model.PostMeta) error {
	s.posts = append(s.posts, meta)
	return nil
}

func (s *memoryMetaStore) SaveSiteIndex(posts []model.PostMeta) error {
	s.siteIndex = posts
	return nil
}

func TestBuildPipeline_Execute(t *testing.T) {
	// Create a temp directory with test .gmi files
	sourceDir := t.TempDir()
	if err := createTestGmiFile(sourceDir, "hello-world.gmi"); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	writer := newMemoryWriter()
	metaStore := &memoryMetaStore{}

	config := BuildConfig{
		SourceDir: sourceDir,
		OutputDir: t.TempDir(),
		SiteTitle: "Test Blog",
		SiteURL:   "https://example.com",
		Author:    "Test Author",
	}

	pipeline := NewBuildPipeline(writer, metaStore, config)
	if err := pipeline.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Check that HTML was generated
	if _, ok := writer.files["posts/hello-world/index.html"]; !ok {
		t.Error("expected posts/hello-world/index.html to be generated")
	}

	// Check that gemtext was copied
	if _, ok := writer.files["posts/hello-world/index.gmi"]; !ok {
		t.Error("expected posts/hello-world/index.gmi to be generated")
	}

	// Check that index.html was generated
	if _, ok := writer.files["index.html"]; !ok {
		t.Error("expected index.html to be generated")
	}

	// Check that atom feed was generated
	if _, ok := writer.files["feed/atom.xml"]; !ok {
		t.Error("expected feed/atom.xml to be generated")
	}

	// Check post metadata
	if len(metaStore.posts) != 1 {
		t.Fatalf("expected 1 post meta, got %d", len(metaStore.posts))
	}
	if metaStore.posts[0].Title != "Hello World" {
		t.Errorf("expected title 'Hello World', got '%s'", metaStore.posts[0].Title)
	}

	// Check site index
	if len(metaStore.siteIndex) != 1 {
		t.Fatalf("expected 1 post in site index, got %d", len(metaStore.siteIndex))
	}
}

func TestBuildPipeline_SkipsDraft(t *testing.T) {
	sourceDir := t.TempDir()
	if err := createTestDraftFile(sourceDir, "draft-post.gmi"); err != nil {
		t.Fatalf("failed to create draft file: %v", err)
	}

	writer := newMemoryWriter()
	metaStore := &memoryMetaStore{}

	config := BuildConfig{
		SourceDir: sourceDir,
		OutputDir: t.TempDir(),
		SiteTitle: "Test Blog",
		SiteURL:   "https://example.com",
		Author:    "Test Author",
	}

	pipeline := NewBuildPipeline(writer, metaStore, config)
	if err := pipeline.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Draft should not generate post HTML
	if _, ok := writer.files["posts/draft-post/index.html"]; ok {
		t.Error("expected draft post to be skipped")
	}
}

func TestBuildPipeline_FileSystemWriter(t *testing.T) {
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	if err := createTestGmiFile(sourceDir, "test-post.gmi"); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	writer := infrastructure.NewFileSystemWriter(outputDir)
	metaStore := infrastructure.NewJsonMetadataStore(outputDir)

	config := BuildConfig{
		SourceDir: sourceDir,
		OutputDir: outputDir,
		SiteTitle: "Test Blog",
		SiteURL:   "https://example.com",
		Author:    "Test Author",
	}

	pipeline := NewBuildPipeline(writer, metaStore, config)
	if err := pipeline.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func createTestGmiFile(dir, name string) error {
	content := "```yaml\ntitle: Hello World\ndate: 2024-01-15\ntags: tech, gemini\nlang: ja\n```\n# Hello World\nThis is a test post.\n=> https://example.com Example Link\n* Item one\n* Item two\n> A wise quote\n"
	return writeTestFile(dir, name, content)
}

func createTestDraftFile(dir, name string) error {
	content := "```yaml\ntitle: Draft Post\ndraft: true\n```\n# Draft Post\nThis is a draft.\n"
	return writeTestFile(dir, name, content)
}

func writeTestFile(dir, name, content string) error {
	return infrastructure.NewFileSystemWriter(dir).WriteBytes(name, []byte(content))
}
