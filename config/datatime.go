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

func getTradingTimestamp() string {
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
		tm = getTradingTimestamp()
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

// Size 获取时段总数
func (s TradingSession) Size() int {
	return len(s.sessions)
}

// Index 判断timestamp是第几个交易时段
func (s TradingSession) Index(timestamp ...string) int {
	var tm string
	if len(timestamp) > 0 {
		tm = strings.TrimSpace(timestamp[0])
	} else {
		tm = getTradingTimestamp()
	}
	for i, timeRange := range s.sessions {
		if timeRange.IsTrading(tm) {
			return i
		}
	}
	return -1
}

// IsTrading 是否交易时段
func (s TradingSession) IsTrading(timestamp ...string) bool {
	index := s.Index(timestamp...)
	if index < 0 {
		return false
	}
	return true
}

// IsTodayLastSession 当前时段是否今天最后一个交易时段
//
//	备选函数名 IsTodayFinalSession
func (s TradingSession) IsTodayLastSession(timestamp ...string) bool {
	n := s.Size()
	index := s.Index(timestamp...)
	if index+1 < n {
		return false
	}
	return true
}

// CanStopLoss 当前时段是否可以进行止损操作
//
//	如果是3个时段, 止损操作在第2时段, 如果是4个时段, 止损在第3个
//	如果是2个时段, 则是第2个时段, 也就是最后一个时段
func (s TradingSession) CanStopLoss(timestamp ...string) bool {
	n := s.Size()
	index := s.Index(timestamp...)
	// 1个时段, 立即止损
	c1 := n == 1
	// 2个时段, 在第二个时间止损
	c2 := n == 2 && index == 1
	// 3个以上时段, 在倒数第2个时段止损
	c3 := n >= 3 && index+2 == n
	if c1 || c2 || c3 {
		return true
	}
	return false
}

// CanTakeProfit 当前时段是否可以止盈
func (s TradingSession) CanTakeProfit(timestamp ...string) bool {
	_ = timestamp
	return true
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
