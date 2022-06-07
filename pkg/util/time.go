package util

import (
	"time"
)

var DefaultDateTimeFormat = "2006-01-02 15:04:05"
var DefaultDateFormat = "2006-01-02"
var TimeLocal, _ = time.LoadLocation("Asia/Shanghai")

// DBMinTimePrecision 数据库 datetime(6) 字段最小操作精度
var DBMinTimePrecision = time.Microsecond

// DBMaxTimeNanoseconds mysql datetime(6) 字段操作最大的纳秒数(超过此值999_999_499秒数会进位)
var DBMaxTimeNanoseconds = 999_999_000

// 未特别说明的方法，默认处理 TimeLocal 本地时区时间

// ===== FORMAT =====

// FormatTime format string 基础方法
func FormatTime(t time.Time, layout string) string {
	return t.In(TimeLocal).Format(layout)
}

// FormatDate return "2006-01-02" format string
func FormatDate(t time.Time) string {
	return FormatTime(t, DefaultDateFormat)
}

// FormatDatePointer return "2006-01-02" format string
func FormatDatePointer(t *time.Time) string {
	if t == nil {
		return ""
	}

	return FormatTime(*t, DefaultDateFormat)
}

// FormatDateTime return "2006-01-02 15:04:05" format string
func FormatDateTime(t time.Time) string {
	return FormatTime(t, DefaultDateTimeFormat)
}

// FormatDateTimePointer return "2006-01-02 15:04:05" format string
func FormatDateTimePointer(t *time.Time) string {
	if t == nil {
		return ""
	}

	return FormatTime(*t, DefaultDateTimeFormat)
}

// ===== FORMAT END =====

// ===== PARSE =====

// ParseStringToTime parse string to time 基础方法
func ParseStringToTime(timeStr, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, timeStr, TimeLocal)
}

// ParseStringToDate return "2021-01-01 00:00:00 +0800 CST"(local time)
func ParseStringToDate(timeStr string) time.Time {
	t, _ := ParseStringToTime(timeStr, DefaultDateFormat)
	return t
}

// ParseStringToDateTime return "2021-01-01 01:01:01 +0800 CST"(local time)
func ParseStringToDateTime(timeStr string) time.Time {
	t, _ := ParseStringToTime(timeStr, DefaultDateTimeFormat)
	return t
}

// ===== PARSE END =====

// AddDate 时间增减
// 类似于ruby中的时间增减，和 time.AddDate 不同
// 如：
// loc, _ := time.LoadLocation("Asia/Shanghai")
// t := time.Date(2010, 3, 31, 12, 0, 0, 0, loc)
// utils.AddDate(t, 0, 1)
// => 2010-04-30 12:00:00 +0800 CST
// 不会因为4月没有31号，而变成5月1号
func AddDate(t time.Time, years int, months int) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	// firstDayOfMonthAfterAddDate: years 年，months 月后的 那个月份的1号
	firstDayOfMonthAfterAddDate := time.Date(year+years, month+time.Month(months), 1,
		hour, min, sec, t.Nanosecond(), t.Location())
	// firstDayOfMonthAfterAddDate 月份的最后一天
	lastDay := LastDayOfMonth(firstDayOfMonthAfterAddDate)

	// 如果 t 的天 > lastDay，则设置为lastDay
	// 如：t 为 2020-03-31 12:00:00 +0800，增加1个月，为4月31号
	// 但是4月没有31号，则设置为4月最后一天lastDay（30号）
	if day > lastDay {
		day = lastDay
	}

	return time.Date(year+years, month+time.Month(months), day,
		hour, min, sec, t.Nanosecond(), t.Location())
}

// Year returns the year in which local time t occurs.
func Year(t time.Time) int {
	return t.In(TimeLocal).Year()
}

// Day returns the day of the month specified by local time of t.
func Day(t time.Time) int {
	return t.In(TimeLocal).Day()
}

// Hour returns the hour within the day specified by local time of t,
// in the range [0, 23].
func Hour(t time.Time) int {
	return t.In(TimeLocal).Hour()
}

// LastDayOfMonth returns the last day of the month specified by local time of t.
func LastDayOfMonth(t time.Time) int {
	return EndOfMonth(t).Day()
}

// EndOfDay end of day(local time)
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.In(TimeLocal).Date()
	return time.Date(year, month, day, 23, 59, 59, DBMaxTimeNanoseconds, TimeLocal)
}

// BeginningOfMonth beginning of month(local time)
func BeginningOfMonth(t time.Time) time.Time {
	year, month, _ := t.In(TimeLocal).Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, TimeLocal)
}

// EndOfMonth end of month(local time)
func EndOfMonth(t time.Time) time.Time {
	// gorm 1.0 不支持自动转化成数据库精度(应该是直接传给数据库), 2.0 没有问题
	// error: gorm update 2020-02-04 15:59:59.999999999 +0000 UTC to 2020-02-04 16:00:00 +0000 UTC
	return BeginningOfMonth(t).AddDate(0, 1, 0).Add(-DBMinTimePrecision)
}
