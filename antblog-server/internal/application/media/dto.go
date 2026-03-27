// Package media 媒体应用层。
package media

import "time"

// ─── 请求 DTO ────────────────────────────────────────────────────────────────

// UploadReq 上传媒体请求（通过 multipart/form-data 提交，字段名 file）
// 注：实际文件数据通过 *multipart.FileHeader 传入 UseCase，DTO 仅描述附加参数。
type UploadReq struct {
	// ArticleID 可选，上传时直接绑定到某篇文章
	ArticleID *uint64 `form:"article_id"`
}

// BindArticleReq 绑定/解绑文章请求
type BindArticleReq struct {
	ArticleID *uint64 `json:"article_id"` // nil = 解绑
}

// AdminListMediaReq 后台媒体列表查询请求
type AdminListMediaReq struct {
	Page       int    `form:"page"        validate:"min=1"`
	PageSize   int    `form:"page_size"   validate:"min=1,max=100"`
	UploaderID uint64 `form:"uploader_id"`
	ArticleID  uint64 `form:"article_id"`
	MimeType   string `form:"mime_type"   validate:"max=128"`
}

// ─── 响应 DTO ────────────────────────────────────────────────────────────────

// MediaResp 媒体资源完整响应
type MediaResp struct {
	ID           uint64     `json:"id"`
	UploaderID   uint64     `json:"uploader_id"`
	ArticleID    *uint64    `json:"article_id"`
	OriginalName string     `json:"original_name"`
	URL          string     `json:"url"`
	MimeType     string     `json:"mime_type"`
	FileSize     int64      `json:"file_size"`
	FileSizeHuman string    `json:"file_size_human"` // 可读文件大小，如 "2.3 MB"
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Hash         string     `json:"hash"`
	CreatedAt    time.Time  `json:"created_at"`
}
