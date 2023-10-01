package market

import (
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/securities"
	"sync"
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
		securityCode := proto.CorrectSecurityCode(v)
		mapSubnewStock[securityCode] = true
	}
}

// IsSubNewStock 是否次新股
func IsSubNewStock(code string) bool {
	onceSubNew.Do(loadSubNewStock)
	securityCode := proto.CorrectSecurityCode(code)
	if v, ok := mapSubnewStock[securityCode]; ok {
		return v
	}
	return false
}
