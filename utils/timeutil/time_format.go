package timeutil

import (
	"strconv"
	"time"

	"github.com/Cospk/base-tools/errs"
)

const (
	TimeOffset = 8 * 3600  // 8小时偏移量
	HalfOffset = 12 * 3600 // 半天小时偏移量
)

// GetCurrentTimestampBySecond 获取当前秒级时间戳
func GetCurrentTimestampBySecond() int64 {
	return time.Now().Unix()
}

// UnixSecondToTime 将时间戳转换为 time.Time 类型
func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

// UnixNanoSecondToTime 将纳秒时间戳转换为 time.Time 类型
func UnixNanoSecondToTime(nanoSecond int64) time.Time {
	return time.Unix(0, nanoSecond)
}

// UnixMillSecondToTime 将毫秒时间戳转换为 time.Time 类型
func UnixMillSecondToTime(millSecond int64) time.Time {
	return time.Unix(0, millSecond*1e6)
}

// GetCurrentTimestampByNano 获取当前纳秒级时间戳
func GetCurrentTimestampByNano() int64 {
	return time.Now().UnixNano()
}

// GetCurrentTimestampByMill 获取当前毫秒级时间戳
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetCurDayZeroTimestamp 获取当天0点的时间戳
func GetCurDayZeroTimestamp() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix() - TimeOffset
}

// GetCurDayHalfTimestamp 获取当天12点的时间戳
func GetCurDayHalfTimestamp() int64 {
	return GetCurDayZeroTimestamp() + HalfOffset

}

// GetCurDayZeroTimeFormat 获取当天0点的格式化时间，格式为 "2006-01-02_00-00-00"
func GetCurDayZeroTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp(), 0).Format("2006-01-02_15-04-05")
}

// GetCurDayHalfTimeFormat 获取当天12点的格式化时间，格式为 "2006-01-02_12-00-00"
func GetCurDayHalfTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp()+HalfOffset, 0).Format("2006-01-02_15-04-05")
}

// GetTimeStampByFormat 将字符串转换为 Unix 时间戳
func GetTimeStampByFormat(datetime string) string {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix()
	return strconv.FormatInt(timestamp, 10)
}

// TimeStringFormatTimeUnix 将字符串转换为 Unix 时间戳
func TimeStringFormatTimeUnix(timeFormat string, timeSrc string) int64 {
	tm, _ := time.Parse(timeFormat, timeSrc)
	return tm.Unix()
}

// TimeStringToTime 将字符串转换为 time.Time
func TimeStringToTime(timeString string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", timeString)
	return t, errs.WrapMsg(err, "timeStringToTime failed", "timeString", timeString)
}

// TimeToString 将 time.Time 转换为字符串
func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func GetCurrentTimeFormatted() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetTimestampByTimezone 根据时区获取特定时间戳
func GetTimestampByTimezone(timezone string) (int64, error) {
	// set time zone
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, errs.New("error loading location:", "error:", err)
	}
	// get current time
	currentTime := time.Now().In(location)
	// get timestamp
	timestamp := currentTime.Unix()
	return timestamp, nil
}

func DaysBetweenTimestamps(timezone string, timestamp int64) (int, error) {
	// set time zone
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, errs.New("error loading location:", "error:", err)
	}
	// get current time
	now := time.Now().In(location)
	// timestamp to time
	givenTime := time.Unix(timestamp, 0)
	// calculate duration
	duration := now.Sub(givenTime)
	// change to days
	days := int(duration.Hours() / 24)
	return days, nil
}

// IsSameWeekday judge current day and specific day is the same of a week.
func IsSameWeekday(timezone string, timestamp int64) (bool, error) {
	// set time zone
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return false, errs.New("error loading location:", "error:", err)
	}
	// get current weekday
	currentWeekday := time.Now().In(location).Weekday()
	// change timestamp to weekday
	givenTime := time.Unix(timestamp, 0)
	givenWeekday := givenTime.Weekday()
	// compare two days
	return currentWeekday == givenWeekday, nil
}

func IsSameDayOfMonth(timezone string, timestamp int64) (bool, error) {
	// set time zone
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return false, errs.New("error loading location:", "error:", err)
	}
	// Get the current day of the month
	currentDay := time.Now().In(location).Day()
	// Convert the timestamp to time and get the day of the month
	givenDay := time.Unix(timestamp, 0).Day()
	// Compare the days
	return currentDay == givenDay, nil
}

func IsWeekday(timestamp int64) bool {
	// Convert the timestamp to time
	givenTime := time.Unix(timestamp, 0)
	// Get the day of the week
	weekday := givenTime.Weekday()
	// Check if the day is between Monday (1) and Friday (5)
	return weekday >= time.Monday && weekday <= time.Friday
}

func IsNthDayCycle(timezone string, startTimestamp int64, n int) (bool, error) {
	// set time zone
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return false, errs.New("error loading location:", "error:", err)
	}
	// Parse the start date
	startTime := time.Unix(startTimestamp, 0)
	if err != nil {
		return false, errs.New("invalid start date format:", "error:", err)
	}
	// Get the current time
	now := time.Now().In(location)
	// Calculate the difference in days between the current time and the start time
	diff := now.Sub(startTime).Hours() / 24
	// Check if the difference in days is a multiple of n
	return int(diff)%n == 0, nil
}

// IsNthWeekCycle checks if the current day is part of an N-week cycle starting from a given start timestamp.
func IsNthWeekCycle(timezone string, startTimestamp int64, n int) (bool, error) {
	// set time zone
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return false, errs.New("error loading location:", "error:", err)
	}

	// Get the current time
	now := time.Now().In(location)

	// Parse the start timestamp
	startTime := time.Unix(startTimestamp, 0)
	if err != nil {
		return false, errs.New("invalid start timestamp format:", "error:", err)
	}

	// Calculate the difference in days between the current time and the start time
	diff := now.Sub(startTime).Hours() / 24

	// Convert days to weeks
	weeks := int(diff) / 7

	// Check if the difference in weeks is a multiple of n
	return weeks%n == 0, nil
}

// IsNthMonthCycle checks if the current day is part of an N-month cycle starting from a given start timestamp.
func IsNthMonthCycle(timezone string, startTimestamp int64, n int) (bool, error) {
	// set time zone
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return false, errs.New("error loading location:", "error:", err)
	}

	// Get the current date
	now := time.Now().In(location)

	// Parse the start timestamp
	startTime := time.Unix(startTimestamp, 0)
	if err != nil {
		return false, errs.New("invalid start timestamp format:", "error:", err)
	}

	// Calculate the difference in months between the current time and the start time
	yearsDiff := now.Year() - startTime.Year()
	monthsDiff := int(now.Month()) - int(startTime.Month())
	totalMonths := yearsDiff*12 + monthsDiff

	// Check if the difference in months is a multiple of n
	return totalMonths%n == 0, nil
}
