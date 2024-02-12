package utils

import (
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
)

func IndexReverse(s pandas.Series) pandas.Series {
	var indexes []int
	rows := s.Len()
	s.Apply(func(idx int, v any) {
		indexes = append(indexes, rows-idx)
	})
	return pandas.ToSeries(indexes...)
}

// SeriesIndexOf 获取序列第n索引的值
//
// Deprecated: 推荐使用 Float64IndexOf
func SeriesIndexOf(s pandas.Series, n int) float64 {
	v := s.IndexOf(n)
	return num.AnyToFloat64(v)
}

// SeriesChangeRate 计算两个序列的净增长
func SeriesChangeRate(base, v pandas.Series) pandas.Series {
	chg := v.Div(base).Sub(1.00).Mul(100)
	return chg
}

func StringIndexOf(s pandas.Series, n int) string {
	v := s.IndexOf(n)
	return num.AnyToString(v)
}

func BoolIndexOf(s pandas.Series, n int) bool {
	v := s.IndexOf(n)
	return num.AnyToBool(v)
}

func Float64IndexOf(s pandas.Series, n int) float64 {
	v := s.IndexOf(n)
	return num.AnyToFloat64(v)
}

func IntegerIndexOf(s pandas.Series, n int) int {
	v := s.IndexOf(n)
	return int(num.AnyToInt64(v))
}

func Int64IndexOf(s pandas.Series, n int) int64 {
	v := s.IndexOf(n)
	return num.AnyToInt64(v)
}
