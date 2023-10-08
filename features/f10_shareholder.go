package features

import (
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/datasets/dfcf"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
)

type top10ShareHolder struct {
	Code           string
	FreeCapital    float64
	Top10Capital   float64
	Top10Change    float64
	ChangeCapital  float64
	IncreaseRatio  float64
	ReductionRatio float64
}

func checkoutShareHolder(securityCode, featureDate string) *top10ShareHolder {
	xdxrs := base.GetCacheXdxrList(securityCode)
	api.SliceSort(xdxrs, func(a, b quotes.XdxrInfo) bool {
		return a.Date > b.Date
	})
	xdxrInfo := checkCapital(xdxrs, featureDate)
	if xdxrInfo != nil && proto.AssertStockBySecurityCode(securityCode) {
		list := dfcf.GetCacheShareHolder(securityCode, featureDate)
		capital := xdxrInfo.HouLiuTong * 10000
		totalCapital := xdxrInfo.HouZongGuBen * 10000
		top10Capital, freeCapital, capitalChanged, increaseRatio, reductionRatio := ComputeFreeCapital(list, capital)
		if freeCapital < 0 {
			top10Capital, freeCapital, capitalChanged, increaseRatio, reductionRatio = ComputeFreeCapital(list, totalCapital)
		}
		frontList := dfcf.GetCacheShareHolder(securityCode, featureDate, 2)
		frontTop10Capital, _, _, _, _ := ComputeFreeCapital(frontList, totalCapital)
		shareHolder := top10ShareHolder{
			Code:           securityCode,
			FreeCapital:    freeCapital,
			Top10Capital:   top10Capital,
			Top10Change:    top10Capital - frontTop10Capital,
			ChangeCapital:  capitalChanged,
			IncreaseRatio:  increaseRatio,
			ReductionRatio: reductionRatio,
		}
		return &shareHolder
	}
	return nil
}
