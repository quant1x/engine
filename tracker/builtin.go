package tracker

import (
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gox/concurrent"
	"sync"
)

const (
	SecurityUnknown = "unknown"
)

// 个股评估
func evaluate(api models.Strategy, wg *sync.WaitGroup, code string, result *concurrent.TreeMap[string, models.ResultInfo]) {
	defer wg.Done()
	api.Evaluate(code, result)
}
