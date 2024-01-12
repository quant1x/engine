package tracker

import (
	"gitee.com/quant1x/engine/factors"
)

// SectorSort 板块排序
func v1SectorSort(a, b factors.QuoteSnapshot) bool {
	if a.ChangeRate > b.ChangeRate {
		return true
	}
	if a.ChangeRate == b.ChangeRate && a.Amount > b.Amount {
		return true
	}
	if a.ChangeRate == b.ChangeRate && a.Amount == b.Amount && a.OpenTurnZ > b.OpenTurnZ {
		return true
	}
	return false
}

// SectorSort 板块排序
func SectorSort(a, b factors.QuoteSnapshot) bool {
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
