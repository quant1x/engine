package market

import (
	"sync"

	"github.com/quant1x/exchange"
	"github.com/quant1x/gotdx/securities"
)

const (
	// 次新股板块
	kBlockSubNewStock = "880529"
)

var (
	onceSubNew     sync.Once
	mapSubnewStock = map[string]bool{}
)

func loadSubNewStock() {
	blockInfo := securities.GetBlockInfo(kBlockSubNewStock)
	if blockInfo == nil {
		return
	}
	for _, v := range blockInfo.ConstituentStocks {
		securityCode := exchange.CorrectSecurityCode(v)
		mapSubnewStock[securityCode] = true
	}
}

// IsSubNewStock 是否次新股
func IsSubNewStock(code string) bool {
	onceSubNew.Do(loadSubNewStock)
	securityCode := exchange.CorrectSecurityCode(code)
	if v, ok := mapSubnewStock[securityCode]; ok {
		return v
	}
	return false
}
