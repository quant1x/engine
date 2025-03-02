package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/pkg/runewidth"
	"gitee.com/quant1x/pkg/tablewriter"
	"os"
	"sync"
	"time"
)

// FeaturesBackTest FeaturesUpdate 特征-数据有效性验证
func FeaturesBackTest(barIndex *int, cacheDate, featureDate string, plugins []cache.DataAdapter, op cache.OpKind) MetricCallback {
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
	var metrics []cache.AdapterMetric
	for _, adapter := range adapters {
		logger.Infof("%s: %s, begin", moduleName, adapter.Name())

		wgAdapter.Add(1)
		// 加载指定日期的特征
		adapter.Checkout(cacheDate)
		var sb cache.ScoreBoard

		barCode := progressbar.NewBar(*barIndex+1, "执行["+adapter.Name()+"]", codeCount)
		for _, code := range allCodes {
			now := time.Now()
			passed := false
			raw := adapter.Element(code)
			kind := adapter.Kind()
			sb.From(cache.GetDataAdapter(kind))
			// 判断是否实现验证接口
			feature, ok := raw.(cache.Validator)
			if ok {
				err := feature.Check(featureDate)
				if err == nil {
					passed = true
				}
			}
			sb.Add(1, time.Since(now), passed)
			barCode.Add(1)
			//wg.Add(1)
			//go updateStockFeature(wg, barCode, feature, code, cacheDate, featureDate, op, mapFeature, &sb, now)
		}
		//wg.Wait()
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
	// 输出衡量性能的指标列表
	mcb := func() {
		fmt.Println()
		metricCount := len(metrics)
		if metricCount > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(tags.GetHeadersByTags(cache.ScoreBoard{}))
			for _, v := range metrics {
				table.Append(tags.GetValuesByTags(v))
			}
			table.Render()
		}
	}
	return mcb
}
