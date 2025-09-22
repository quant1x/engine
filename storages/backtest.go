package storages

import (
	"sync"
	"time"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/pkg/runewidth"
)

// FeaturesBackTest FeaturesUpdate 特征-数据有效性验证
func FeaturesBackTest(barIndex *int, cacheDate, featureDate string, plugins []cache.DataAdapter, op cache.OpKind) []cache.FactorMetrics {
	moduleName := cache.OpMap[op] + "特征数据"
	moduleName += cacheDate
	var adapters []factors.FeatureRotationAdapter
	maxWidth := 0
	for _, plugin := range plugins {
		adapter, ok := plugin.(factors.FeatureRotationAdapter)
		if ok {
			adapters = append(adapters, adapter)
			width := runewidth.StringWidth(adapter.Name())
			if width > maxWidth {
				maxWidth = width
			}
		}
	}
	logger.Infof("%s: all, begin", moduleName)

	var wgAdapter sync.WaitGroup
	cacheCount := len(adapters)
	barAdapter := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", cacheCount)
	allCodes := market.GetCodeList()
	allCodes = allCodes[:]
	codeCount := len(allCodes)
	var metrics []cache.FactorMetrics
	for _, adapter := range adapters {
		logger.Infof("%s: %s, begin", moduleName, adapter.Name())

		wgAdapter.Add(1)
		// 加载指定cacheDate日期的特征, 那么数据本身是featureDate日期
		adapter.Checkout(cacheDate)
		var sb cache.ScoreBoard

		barCode := progressbar.NewBar(*barIndex+1, "执行["+adapter.Name()+"]", codeCount)
		for _, code := range allCodes {
			now := time.Now()
			var hasSignal bool
			var passed bool
			raw := adapter.Element(code)
			kind := adapter.Kind()
			sb.From(cache.GetDataAdapter(kind))
			// 判断是否实现验证接口
			feature, ok := raw.(cache.FactorSignalEvaluator)
			if ok {
				var err error
				hasSignal, err = feature.Check(cacheDate, featureDate)
				if err == nil {
					passed = true
				}
			} else {
				// 未实现Validator接口, 视为仅统计样本, 不产生信号且不算通过
			}
			sb.Add(1, time.Since(now), hasSignal, passed)
			barCode.Add(1)
		}
		barCode.Wait()

		// 适配器进度条+1
		barAdapter.Add(1)
		wgAdapter.Done()
		metrics = append(metrics, sb.Metric())
		logger.Infof("%s: %s, end", moduleName, adapter.Name())
	}
	wgAdapter.Wait()
	barAdapter.Wait()
	logger.Infof("%s: all, end", moduleName)

	return metrics
}
