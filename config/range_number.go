package config

import (
	"fmt"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pandas/stat"
	"gitee.com/quant1x/pkg/yaml"
	"regexp"
	"strings"
)

// 值范围正则表达式
var (
	numberRangePattern = "[~]\\s*"
	numberRangeRegexp  = regexp.MustCompile(numberRangePattern)
)

// NumberRange 数值范围
type NumberRange struct {
	begin float64
	end   float64
}

func (this NumberRange) String() string {
	return fmt.Sprintf("{begin: %f, end: %f}", this.begin, this.end)
}

func (this *NumberRange) init() {
	this.begin = stat.MinFloat64
	this.end = stat.MaxFloat64
}

func (this *NumberRange) Parse(text string) error {
	text = strings.TrimSpace(text)
	arr := numberRangeRegexp.Split(text, -1)
	if len(arr) != 2 {
		this.init()
		return ErrRangeFormat
	}
	begin := strings.TrimSpace(arr[0])
	end := strings.TrimSpace(arr[1])
	this.begin = stat.AnyToFloat64(begin)
	this.end = stat.AnyToFloat64(end)
	if this.begin > this.end {
		this.begin, this.end = this.end, this.begin
	}
	return nil
}

// UnmarshalYAML YAML自定义解析
func (this *NumberRange) UnmarshalYAML(node *yaml.Node) error {
	var value string
	if len(node.Content) == 0 {
		value = node.Value
	} else if len(node.Content) == 2 {
		value = node.Content[1].Value
	} else {
		this.init()
		return ErrRangeFormat
	}

	return this.Parse(value)
}

// UnmarshalText 设置默认值调用
func (this *NumberRange) UnmarshalText(text []byte) error {
	return this.Parse(api.Bytes2String(text))
}

// Validate 验证
func (this *NumberRange) Validate(v float64) bool {
	if this.begin == 0 && this.end == 0 {
		return true
	} else if v >= this.begin && v < this.end {
		return true
	}
	return false
}
