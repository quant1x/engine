package factors

import (
	"gitee.com/quant1x/engine/datasource/dfcf"
)

// ComputeFreeCapital 计算自由流通股本
func ComputeFreeCapital(holderList []dfcf.CirculatingShareholder, capital float64) (top10Capital, freeCapital, capitalChanged, increaseRatio, reductionRatio float64) {
	increase := 0
	reduce := 0
	for k, holder := range holderList {
		top10Capital += float64(holder.HoldNum)
		capitalChanged += float64(holder.HoldNumChange)
		if holder.HoldNumChange >= 0 {
			increase += holder.HoldNumChange
		} else {
			reduce += holder.HoldNumChange
		}
		if k >= 10 {
			continue
		}
		if holder.FreeHoldNumRatio >= 1.00 && holder.IsHoldOrg == "1" {
			capital -= float64(holder.HoldNum)
		}
	}
	increaseRatio = 100.0000 * (float64(increase) / top10Capital)
	reductionRatio = 100.0000 * (float64(reduce) / top10Capital)
	return top10Capital, capital, capitalChanged, increaseRatio, reductionRatio
}
