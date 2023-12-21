package config

import (
	"fmt"
	"gitee.com/quant1x/gox/exception"
	"gitee.com/quant1x/pkg/yaml"
	"strings"
	"time"
)

var (
	errnoConfig       = 0
	ErrTimeFormat     = exception.New(errnoConfig+1, "时间格式错误")
	ErrRangeFormat    = exception.New(errnoConfig+2, "数值范围格式错误")
	formatOfTimestamp = time.TimeOnly
)

func getTradingTimestamp() string {
	now := time.Now()
	return now.Format(formatOfTimestamp)
}

// TimeRange 时间范围
type TimeRange struct {
	begin string // 开始时间
	end   string // 结束时间
}

func (this TimeRange) String() string {
	return fmt.Sprintf("{begin: %s, end: %s}", this.begin, this.end)
}

func (this *TimeRange) Parse(text string) error {
	text = strings.TrimSpace(text)
	arr := valueRangeRegexp.Split(text, -1)
	if len(arr) != 2 {
		return ErrTimeFormat
	}
	this.begin = strings.TrimSpace(arr[0])
	this.end = strings.TrimSpace(arr[1])
	if this.begin > this.end {
		this.begin, this.end = this.end, this.begin
	}
	return nil
}

// UnmarshalText 设置默认值调用
func (this *TimeRange) UnmarshalText(text []byte) error {
	//TODO implement me
	panic("implement me")
}

func (this *TimeRange) UnmarshalYAML(node *yaml.Node) error {
	var key, value string
	if len(node.Content) == 0 {
		value = node.Value
	} else if len(node.Content) == 2 {
		key = node.Content[0].Value
		value = node.Content[1].Value
	}
	_ = key
	return this.Parse(value)
}

func (this *TimeRange) IsTrading(timestamp ...string) bool {
	var tm string
	if len(timestamp) > 0 {
		tm = strings.TrimSpace(timestamp[0])
	} else {
		tm = getTradingTimestamp()
	}
	if tm >= this.begin && tm <= this.end {
		return true
	}
	return false
}
