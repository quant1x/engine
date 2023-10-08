package storages

import (
	"gitee.com/quant1x/engine/datasets"
	"gitee.com/quant1x/engine/flash"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/progressbar"
	"github.com/mattn/go-runewidth"
	"strings"
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

// RepairBaseData 更新基础数据
func RepairBaseData(barIndex *int, cacheDate, featureDate string) {
	const useGoroutine = false
	moduleName := "更新基础数据"
	allCodes := market.GetCodeList()
	dataSetList := flash.DataSetList()
	dataSetCount := len(dataSetList)
	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", dataSetCount)
	var wg sync.WaitGroup
	maxWidth := 0
	for _, dataSet := range dataSetList {
		width := runewidth.StringWidth(dataSet.Name())
		if width > maxWidth {
			maxWidth = width
		}
	}
	for sequence, dataSet := range dataSetList {
		_ = dataSet.Init(barIndex, featureDate)
		codeCount := len(allCodes)
		//format := fmt.Sprintf("%%%ds", maxWidth)
		//title := fmt.Sprintf(format, dataSet.Name())
		width := runewidth.StringWidth(dataSet.Name())
		title := strings.Repeat(" ", maxWidth-width) + dataSet.Name()
		barNo := *barIndex + 1
		if useGoroutine {
			barNo += sequence
		}
		barCode := progressbar.NewBar(barNo, "执行["+title+"]", codeCount)
		wg.Add(1)
		if useGoroutine {
			go repairDateSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate)
		} else {
			repairDateSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate)
		}
	}
	wg.Wait()
}
