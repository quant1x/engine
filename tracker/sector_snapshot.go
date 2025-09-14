package tracker

import (
	"fmt"

	"github.com/quant1x/engine/factors"
	"github.com/quant1x/engine/models"
	"github.com/quant1x/gotdx/securities"
	"github.com/quant1x/num"
	"github.com/quant1x/x/api"
	"github.com/quant1x/x/progressbar"
)

// 板块扫描
func scanSectorSnapshots(pbarIndex *int, blockType securities.BlockType, isHead bool) (list []factors.QuoteSnapshot) {
	// 执行板块指数的检测
	blockInfos := securities.BlockList()
	// 获取指定类型的板块代码列表
	var blockCodes []string
	for _, v := range blockInfos {
		if v.Type != blockType {
			continue
		}
		blockCode := v.Code
		blockCodes = append(blockCodes, blockCode)
		blockTypeName, _ := securities.BlockTypeNameByTypeCode(v.Type)
		__mapBlockTypeName[blockCode] = blockTypeName
	}

	blockCount := len(blockCodes)
	fmt.Println()
	btn, ok := securities.BlockTypeNameByTypeCode(blockType)
	if !ok {
		btn = num.AnyToString(blockType)
	}
	bar := progressbar.NewBar(*pbarIndex, "执行[扫描"+btn+"板块指数]", blockCount)
	*pbarIndex++
	for i := 0; i < blockCount; i++ {
		bar.Add(1)
		blockCode := blockCodes[i]
		snapshot := models.GetStrategySnapshot(blockCode)
		if snapshot == nil {
			continue
		}
		list = append(list, *snapshot)
	}
	if isHead {
		api.SliceSort(list, SectorSortForHead)
	} else {
		api.SliceSort(list, SectorSortForTick)
	}
	return list
}
