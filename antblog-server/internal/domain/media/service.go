package media

import (
	"mime"
	"path/filepath"
	"strings"

	apperrors "antblog/pkg/errors"
)

// IStorageDriver 存储驱动接口（由 infrastructure/storage 层实现，领域层仅依赖此接口）
type IStorageDriver interface {
	// Save 将文件数据持久化，返回存储相对路径和访问 URL
	Save(originalName string, data []byte) (storagePath, url string, err error)

	// Delete 按存储路径删除文件
	Delete(storagePath string) error
}

// IDomainService 媒体领域服务接口
type IDomainService interface {
	// ValidateUpload 校验上传参数（文件大小、MIME 类型白名单）
	ValidateUpload(originalName string, data []byte, allowedTypes []string, maxSize int64) error

	// BuildMedia 构建媒体实体（含 SHA256 计算）
	BuildMedia(uploaderID uint64, originalName string, data []byte, storagePath, url, mimeType string, width, height int) (*Media, error)

	// DetectMimeType 从文件头字节检测 MIME 类型（更可靠，不依赖扩展名）
	DetectMimeType(data []byte, originalName string) string
}

// DomainService 媒体领域服务实现
type DomainService struct{}

// NewDomainService 创建媒体领域服务
func NewDomainService() IDomainService {
	return &DomainService{}
}

func (s *DomainService) ValidateUpload(originalName string, data []byte, allowedTypes []string, maxSize int64) error {
	if len(data) == 0 {
		return apperrors.ErrInvalidParams("文件内容不能为空")
	}
	if int64(len(data)) > maxSize {
		return apperrors.New(apperrors.CodeMediaTooLarge,
			apperrors.Message(apperrors.CodeMediaTooLarge))
	}

	detected := s.DetectMimeType(data, originalName)
	for _, t := range allowedTypes {
		if strings.EqualFold(t, detected) {
			return nil
		}
	}
	return apperrors.New(apperrors.CodeMediaInvalidType,
		apperrors.Message(apperrors.CodeMediaInvalidType))
}

func (s *DomainService) BuildMedia(
	uploaderID uint64,
	originalName string,
	data []byte,
	storagePath, url, mimeType string,
	width, height int,
) (*Media, error) {
	if storagePath == "" || url == "" {
		return nil, apperrors.ErrInvalidParams("存储路径和访问 URL 不能为空")
	}
	return &Media{
		UploaderID:   uploaderID,
		OriginalName: originalName,
		StoragePath:  storagePath,
		URL:          url,
		MimeType:     mimeType,
		FileSize:     int64(len(data)),
		Width:        width,
		Height:       height,
		Hash:         Sha256Hex(data), // 使用公开函数
	}, nil
}

// DetectMimeType 检测 MIME 类型：优先从文件头魔数判断，降级到扩展名
func (s *DomainService) DetectMimeType(data []byte, originalName string) string {
	if len(data) >= 4 {
		// JPEG: FF D8 FF
		if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			return "image/jpeg"
		}
		// PNG: 89 50 4E 47
		if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			return "image/png"
		}
		// GIF: GIF87a / GIF89a
		if len(data) >= 6 && (string(data[:6]) == "GIF87a" || string(data[:6]) == "GIF89a") {
			return "image/gif"
		}
		// WebP: RIFF????WEBP
		if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
			return "image/webp"
		}
	}

	// 降级：根据扩展名判断
	ext := strings.ToLower(filepath.Ext(originalName))
	if mt := mime.TypeByExtension(ext); mt != "" {
		return strings.Split(mt, ";")[0]
	}
	return "application/octet-stream"
}
