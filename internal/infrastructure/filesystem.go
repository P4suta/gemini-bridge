// internal/infrastructure/filesystem.go
package infrastructure

import (
	"io"
	"os"
	"path/filepath"
)

// FileSystemWriter はport.ContentWriterのファイルシステム実装。
// SSGのビルド出力をローカルディレクトリに書き出す。
type FileSystemWriter struct {
	baseDir string
}

func NewFileSystemWriter(baseDir string) *FileSystemWriter {
	return &FileSystemWriter{baseDir: baseDir}
}

func (w *FileSystemWriter) Write(path string, content io.Reader) error {
	fullPath := filepath.Join(w.baseDir, path)

	// ディレクトリの自動作成
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, content)
	return err
}

func (w *FileSystemWriter) WriteBytes(path string, data []byte) error {
	fullPath := filepath.Join(w.baseDir, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0644)
}
