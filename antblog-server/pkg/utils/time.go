package utils

import (
	"time"
)

const (
	LayoutDate     = "2006-01-02"
	LayoutDateTime = "2006-01-02 15:04:05"
	LayoutDateTimeT = "2006-01-02T15:04:05"
	LayoutMonth    = "2006-01"
	LayoutYear     = "2006"
)

// Now 返回当前本地时间
func Now() time.Time {
	return time.Now().Local()
}

// FormatDate 格式化为日期字符串 "2006-01-02"
func FormatDate(t time.Time) string {
	return t.Format(LayoutDate)
}

// FormatDateTime 格式化为日期时间字符串 "2006-01-02 15:04:05"
func FormatDateTime(t time.Time) string {
	return t.Format(LayoutDateTime)
}

// ParseDate 解析日期字符串
func ParseDate(s string) (time.Time, error) {
	return time.ParseInLocation(LayoutDate, s, time.Local)
}

// ParseDateTime 解析日期时间字符串
func ParseDateTime(s string) (time.Time, error) {
	return time.ParseInLocation(LayoutDateTime, s, time.Local)
}

// BeginOfDay 返回某天 00:00:00
func BeginOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// EndOfDay 返回某天 23:59:59
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, 999999999, t.Location())
}

// BeginOfMonth 返回当月第一天 00:00:00
func BeginOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 返回当月最后一天 23:59:59
func EndOfMonth(t time.Time) time.Time {
	return BeginOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// DaysInMonth 返回某月天数
func DaysInMonth(t time.Time) int {
	return EndOfMonth(t).Day()
}

// TimeAgo 返回相对时间描述（中文）
// 例：3分钟前、2小时前、昨天、3天前
func TimeAgo(t time.Time) string {
	now := Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "刚刚"
	case diff < time.Hour:
		return intToStr(int(diff.Minutes())) + "分钟前"
	case diff < 24*time.Hour:
		return intToStr(int(diff.Hours())) + "小时前"
	case diff < 48*time.Hour:
		return "昨天"
	case diff < 30*24*time.Hour:
		return intToStr(int(diff.Hours()/24)) + "天前"
	case diff < 365*24*time.Hour:
		return intToStr(int(diff.Hours()/24/30)) + "个月前"
	default:
		return intToStr(int(diff.Hours()/24/365)) + "年前"
	}
}

// UnixToTime 将 Unix 时间戳（秒）转为 time.Time
func UnixToTime(ts int64) time.Time {
	return time.Unix(ts, 0).Local()
}

// TimeToUnix 将 time.Time 转为 Unix 时间戳（秒）
func TimeToUnix(t time.Time) int64 {
	return t.Unix()
}

// IsSameDay 判断两个时间是否同一天
func IsSameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
