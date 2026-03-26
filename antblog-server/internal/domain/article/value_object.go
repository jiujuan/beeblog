package article

import (
	"strings"
	"unicode/utf8"
)

const (
	maxTitleLen   = 256
	maxSummaryLen = 512
	autoSummaryLen = 200 // 自动截取摘要的字符数
)

// Title 文章标题值对象
type Title struct{ value string }

// NewTitle 创建并校验标题
func NewTitle(s string) (Title, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Title{}, errInvalidTitle("标题不能为空")
	}
	if utf8.RuneCountInString(s) > maxTitleLen {
		return Title{}, errInvalidTitle("标题不得超过 256 个字符")
	}
	return Title{value: s}, nil
}

func (t Title) String() string { return t.value }

// AutoSummary 若 summary 为空，则从 markdown content 自动截取前 N 个字符作为摘要
func AutoSummary(content, summary string) string {
	if strings.TrimSpace(summary) != "" {
		return summary
	}
	// 去除 Markdown 标记符（简单处理：去除 # * ` > - 等行首符号）
	lines := strings.Split(content, "\n")
	var sb strings.Builder
	for _, line := range lines {
		line = strings.TrimLeft(line, "#>-*` ")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		sb.WriteString(line)
		sb.WriteRune(' ')
		if utf8.RuneCountInString(sb.String()) >= autoSummaryLen {
			break
		}
	}
	result := strings.TrimSpace(sb.String())
	runes := []rune(result)
	if len(runes) > autoSummaryLen {
		return string(runes[:autoSummaryLen]) + "..."
	}
	return result
}

// CountWords 统计 Markdown 内容的有效字数（去除标记后的 Unicode 字符数）
func CountWords(content string) int {
	// 去除 Markdown 代码块
	inCode := false
	var sb strings.Builder
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "```") {
			inCode = !inCode
			continue
		}
		if inCode {
			continue
		}
		line = strings.TrimLeft(line, "#>-*`| ")
		sb.WriteString(line)
	}
	return utf8.RuneCountInString(strings.TrimSpace(sb.String()))
}

// ─── 私有错误辅助 ────────────────────────────────────────────────────────────

type titleError struct{ msg string }

func (e *titleError) Error() string { return e.msg }

func errInvalidTitle(msg string) error { return &titleError{msg: msg} }
