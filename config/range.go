package config

import (
	"gitee.com/quant1x/pandas/stat"
	"strings"
	_ "unsafe"
)

type ValueType interface {
	~int | ~float64 | ~string
}

func ParseRange[T ValueType](text string) ValueRange[T] {
	text = strings.TrimSpace(text)
	arr := timeRangeRegexp.Split(text, -1)
	if len(arr) != 2 {
		panic(ErrTimeFormat)
	}
	var begin, end T
	begin = stat.GenericParse[T](strings.TrimSpace(arr[0]))
	end = stat.GenericParse[T](strings.TrimSpace(arr[1]))
	if begin > end {
		begin, end = end, begin
	}
	r := ValueRange[T]{
		begin: begin,
		end:   end,
	}
	return r
}

// ValueRange 数值范围
type ValueRange[T ValueType] struct {
	begin T // 最小值
	end   T // 最大值
}

// In 检查是否包含在范围内
func (r ValueRange[T]) In(v T) bool {
	if v < r.begin || v > r.end {
		return false
	}
	return true
}