package storages

import (
	"gitee.com/quant1x/engine/datasets"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/progressbar"
	"sync"
)

// 更新单个数据集
func repairDateSet(wg *sync.WaitGroup, parent, bar *progressbar.Bar, dataSet datasets.DataSet, cacheDate, featureDate string) {
	allCodes := market.GetCodeList()
	for _, code := range allCodes {
		data := dataSet.Clone(cacheDate, code).(datasets.DataSet)
		data.Update(cacheDate, featureDate)
		bar.Add(1)
	}
	parent.Add(1)
	wg.Done()
}
