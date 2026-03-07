// internal/port/writer.go
package port

import "io"

// ContentWriter はビルド成果物の書き出し先を抽象化するインターフェース。
// ファイルシステム、R2直接アップロード、テスト用バッファなどの実装を差し替え可能にする。
type ContentWriter interface {
	// Write は指定パスにコンテンツを書き出す。
	Write(path string, content io.Reader) error

	// WriteBytes は指定パスにバイト列を書き出す。
	WriteBytes(path string, data []byte) error
}
