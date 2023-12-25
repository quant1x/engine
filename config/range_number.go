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
	min float64
	max float64
}

func (this NumberRange) String() string {
	return fmt.Sprintf("{min: %f, max: %f}", this.min, this.max)
}

func (this *NumberRange) init() {
	this.min = stat.MinFloat64
	this.max = stat.MaxFloat64
}

func (this *NumberRange) Max() float64 {
	return this.max
}

func (this *NumberRange) Min() float64 {
	return this.min
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
	this.min = stat.AnyToFloat64(begin)
	this.max = stat.AnyToFloat64(end)
	if this.min > this.max {
		this.min, this.max = this.max, this.min
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
	if this.min == 0 && this.max == 0 {
		return true
	} else if v >= this.min && v < this.max {
		return true
	}
	return false
}
