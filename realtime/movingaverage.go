package realtime

import (
	"gitee.com/quant1x/pandas/formula"
	"gitee.com/quant1x/pandas/stat"
)

// MovingAverage 计算均线范围
func MovingAverage(CLOSE, HIGH, LOW stat.Series, PN int) (ma, half, maMax, maMin stat.Series) {
	half = formula.MA(CLOSE, PN-1)
	r1Half := formula.REF(half, 1)
	maMax = r1Half.Mul(PN - 1).Add(HIGH).Div(PN)
	maMin = r1Half.Mul(PN - 1).Add(LOW).Div(PN)
	ma = r1Half.Mul(PN - 1).Add(CLOSE).Div(PN)
	return ma, half, maMax, maMin
}
