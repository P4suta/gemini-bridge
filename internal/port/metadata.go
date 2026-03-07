// internal/port/metadata.go
package port

import "gemini-bridge/internal/domain/model"

// MetadataStore は記事メタデータの永続化を抽象化するインターフェース。
type MetadataStore interface {
	// SavePostMeta は個別記事のメタデータを永続化する。
	SavePostMeta(meta model.PostMeta) error

	// SaveSiteIndex はサイト全体の記事インデックスを永続化する。
	SaveSiteIndex(posts []model.PostMeta) error
}
