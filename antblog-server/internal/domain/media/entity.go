// Package media 媒体资源领域层 —— 纯业务实体，零外部依赖。
package media

import "time"

// ─── 枚举 ────────────────────────────────────────────────────────────────────

// MediaType 媒体资源类型
type MediaType string

const (
	TypeImage MediaType = "image" // 图片（jpeg/png/gif/webp）
	TypeFile  MediaType = "file"  // 通用文件（未来扩展）
)

// ─── 聚合根 ──────────────────────────────────────────────────────────────────

// Media 媒体资源聚合根
type Media struct {
	ID           uint64
	UploaderID   uint64    // 上传者用户 ID
	ArticleID    *uint64   // 关联文章 ID（nil = 未绑定）
	OriginalName string    // 原始文件名（含扩展名）
	StoragePath  string    // 服务器存储相对路径（如 uploads/2024/01/xxx.jpg）
	URL          string    // 对外访问 URL
	MimeType     string    // MIME 类型，如 image/jpeg
	FileSize     int64     // 字节数
	Width        int       // 图片宽度（px）；非图片为 0
	Height       int       // 图片高度（px）；非图片为 0
	Hash         string    // SHA256（去重、防重复上传）
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ─── 业务方法 ────────────────────────────────────────────────────────────────

// BindArticle 将媒体资源绑定到文章
func (m *Media) BindArticle(articleID uint64) {
	m.ArticleID = &articleID
	m.UpdatedAt = time.Now()
}

// UnbindArticle 解绑文章
func (m *Media) UnbindArticle() {
	m.ArticleID = nil
	m.UpdatedAt = time.Now()
}

// IsImage 是否为图片资源
func (m *Media) IsImage() bool {
	switch m.MimeType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

// MediaType 返回资源类型
func (m *Media) Type() MediaType {
	if m.IsImage() {
		return TypeImage
	}
	return TypeFile
}
