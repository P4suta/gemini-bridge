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
	data, err := io.ReadAll(content)
	if err != nil {
		return err
	}
	return w.WriteBytes(path, data)
}

func (w *FileSystemWriter) WriteBytes(path string, data []byte) error {
	fullPath := filepath.Join(w.baseDir, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0644)
}
