package utils

import "time"

// Timestamp 毫秒时间戳
func Timestamp() (timestamp int64) {
	now := time.Now()
	timestamp = now.UnixMilli()
	return timestamp
}
