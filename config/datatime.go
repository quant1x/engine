package config

import (
	"fmt"
	"gitee.com/quant1x/gox/exception"
	"regexp"
	"strings"
	"time"
)

var (
	errnoConfig       = 0
	ErrTimeFormat     = exception.New(errnoConfig+1, "时间格式错误")
	formatOfTimestamp = time.TimeOnly
)

var (
	timeRangePattern = "[~-]\\s*"
	timeRangeRegexp  = regexp.MustCompile(timeRangePattern)
)

func getTimestamp() string {
	now := time.Now()
	return now.Format(formatOfTimestamp)
}

func ParseTimeRange(text string) TimeRange {
	text = strings.TrimSpace(text)
	arr := timeRangeRegexp.Split(text, -1)
	if len(arr) != 2 {
		panic(ErrTimeFormat)
	}
	var begin, end string
	begin = strings.TrimSpace(arr[0])
	end = strings.TrimSpace(arr[1])
	if begin > end {
		begin, end = end, begin
	}
	return TimeRange{
		begin: begin,
		end:   end,
	}
}

// TimeRange 时间范围
type TimeRange struct {
	begin string // 开始时间
	end   string // 结束时间
}

func (r TimeRange) String() string {
	return fmt.Sprintf("{begin: %s, end: %s}", r.begin, r.end)
}

func (r TimeRange) IsTrading(timestamp ...string) bool {
	var tm string
	if len(timestamp) > 0 {
		tm = strings.TrimSpace(timestamp[0])
	} else {
		tm = getTimestamp()
	}
	if tm >= r.begin && tm <= r.end {
		return true
	}
	return false
}

var (
	tradingSessionPattern = "[,]\\s*"
	tradingSessionRegexp  = regexp.MustCompile(tradingSessionPattern)
)

func ParseTradingSession(text string) TradingSession {
	text = strings.TrimSpace(text)
	arr := tradingSessionRegexp.Split(text, -1)
	var sessions []TimeRange
	for _, v := range arr {
		tr := ParseTimeRange(v)
		sessions = append(sessions, tr)
	}
	return TradingSession{sessions: sessions}
}

// TradingSession 交易时段
type TradingSession struct {
	sessions []TimeRange
}

// IsTrading 是否交易时段
func (s TradingSession) IsTrading(timestamp ...string) bool {
	var tm string
	if len(timestamp) > 0 {
		tm = strings.TrimSpace(timestamp[0])
	} else {
		tm = getTimestamp()
	}
	for _, timeRange := range s.sessions {
		if timeRange.IsTrading(tm) {
			return true
		}
	}
	return false
}

func (s TradingSession) String() string {
	builder := strings.Builder{}
	builder.WriteByte('[')
	var arr []string
	for _, timeRange := range s.sessions {
		arr = append(arr, timeRange.String())
	}
	builder.WriteString(strings.Join(arr, ","))
	builder.WriteByte(']')
	return builder.String()
}
