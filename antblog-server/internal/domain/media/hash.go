package media

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256Hex 计算文件字节数据的 SHA256 哈希（十六进制字符串，公开供应用层调用）
func Sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
