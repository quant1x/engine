package tracker

import (
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/num"
)

func tickWeight(snapshot factors.QuoteSnapshot) float64 {
	// 1. 板块涨幅
	weight1 := snapshot.ChangeRate
	// 2. 涨速
	weight2 := snapshot.Rate

	return num.Mean([]float64{weight1, weight2})
}

// SectorSortForTick 板块排序, 盘中
func SectorSortForTick(a, b factors.QuoteSnapshot) bool {
	aWeight := tickWeight(a)
	bWeight := tickWeight(b)
	return aWeight > bWeight
}

// SectorSortForHead 板块排序, 早盘
func SectorSortForHead(a, b factors.QuoteSnapshot) bool {
	if a.OpeningChangeRate > b.OpeningChangeRate {
		return true
	}
	if a.OpeningChangeRate == b.OpeningChangeRate && a.Amount > b.Amount {
		return true
	}
	if a.OpeningChangeRate == b.OpeningChangeRate && a.Amount == b.Amount && a.OpenTurnZ > b.OpenTurnZ {
		return true
	}
	return false
}

// StockSort 个股排序
func StockSort(a, b factors.QuoteSnapshot) bool {
	if a.ChangeRate > b.ChangeRate {
		return true
	}
	if a.ChangeRate == b.ChangeRate && a.OpenTurnZ > b.OpenTurnZ {
		return true
	}
	return false
}
