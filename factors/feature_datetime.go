package factors

import (
	"gitee.com/quant1x/engine/cache"
	"time"
)

// GetTimestamp 时间戳
//
//	格式: YYYY-MM-DD hh:mm:ss.SSS
func GetTimestamp() string {
	now := time.Now()
	return now.Format(cache.TimeStampMilli)
}
