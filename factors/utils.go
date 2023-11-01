package factors

import (
	"gitee.com/quant1x/pandas/stat"
)

// SeriesIndexOf 获取序列第n索引的值
func SeriesIndexOf(s stat.Series, n int) float64 {
	v := s.IndexOf(n)
	return stat.AnyToFloat64(v)
}

// SeriesChangeRate 计算两个序列的净增长
func SeriesChangeRate(base, v stat.Series) stat.Series {
	chg := v.Div(base).Sub(1.00).Mul(100)
	return chg
}

func BoolIndexOf(s stat.Series, n int) bool {
	v := s.IndexOf(n)
	return stat.AnyToBool(v)
}

func IntegerIndexOf(s stat.Series, n int) int {
	v := s.IndexOf(n)
	return int(stat.AnyToInt64(v))
}

func StringIndexOf(s stat.Series, n int) string {
	v := s.IndexOf(n)
	return stat.AnyToString(v)
}

func IndexReverse(s stat.Series) stat.Series {
	var indexes []int
	rows := s.Len()
	s.Apply(func(idx int, v any) {
		indexes = append(indexes, rows-idx)
	})
	return stat.NewSeries(indexes...)
}
