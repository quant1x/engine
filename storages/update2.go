package storages

import (
	"gitee.com/quant1x/engine/features"
	"gitee.com/quant1x/engine/flash"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/progressbar"
)

// UpdateBaseData 更新基础数据
func UpdateBaseData(barIndex *int, cacheDate, featureDate string) {
	moduleName := "更新基础数据"
	allCodes := market.GetCodeList()
	cacheList := flash.DataSetList()
	cacheCount := len(cacheList)
	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", cacheCount)
	for _, dataSet := range cacheList {
		codeCount := len(allCodes)
		barCode := progressbar.NewBar(*barIndex+1, "执行["+dataSet.Name()+"]", codeCount)
		for _, code := range allCodes {
			data := dataSet.Clone(cacheDate, code).(features.DataSet)
			data.Update(cacheDate, featureDate)
			barCode.Add(1)
		}
		barCache.Add(1)
	}
}
