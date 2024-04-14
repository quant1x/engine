package tracker

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pkg/tablewriter"
	"os"
	"slices"
	"sort"
)

// ScanSectorForTick 扫描板块
func ScanSectorForTick(barIndex *int) []string {
	// 不分板块类型, 所有的板块放在一起排序
	allBlocks := scanBlockByTypeForTick(barIndex, securities.BK_GAINIAN)
	// 扫描板块内个股排名
	blockCount := len(allBlocks)
	fmt.Println()
	bar := progressbar.NewBar(*barIndex, "执行[板块个股涨幅扫描]", blockCount)
	for i := 0; i < blockCount; i++ {
		bar.Add(1)
		block := &allBlocks[i]
		topCode := SecurityUnknown
		topName := SecurityUnknown
		topRate := 0.00
		stockCodes := block.StockCodes
		var stockSnapshots []factors.QuoteSnapshot
		for _, s := range stockCodes {
			stockCode := s
			securityCode := exchange.CorrectSecurityCode(stockCode)
			tmpBlocks, _ := __stock2Block[securityCode]
			if len(tmpBlocks) == 0 {
				tmpBlocks = make([]SectorInfo, 0)
			}
			tmpBlocks = append(tmpBlocks, *block)
			__stock2Block[securityCode] = tmpBlocks
			snapshot := models.GetStrategySnapshot(securityCode)
			if snapshot == nil {
				continue
			}
			stockSnapshots = append(stockSnapshots, *snapshot)
		}
		if len(stockSnapshots) == 0 {
			continue
		}

		sort.Slice(stockSnapshots, func(i, j int) bool {
			a := stockSnapshots[i]
			b := stockSnapshots[j]
			return StockSort(a, b)
		})

		var stockList []string
		//stockTopNInSector := rules.RuleParameters.StockTopNInSector
		stockCountOfSector := len(stockSnapshots)
		for j := 0; j < stockCountOfSector; /* && j < stockTopNInSector*/ j++ {
			si := stockSnapshots[j]
			stockCode := si.SecurityCode
			if market.IsNeedIgnore(stockCode) {
				continue
			}
			stockList = append(stockList, stockCode)
		}
		__block2Top[block.Code] = stockList
		stockTopList := slices.Clone(stockSnapshots)
		sort.Slice(stockTopList, func(i, j int) bool {
			a := stockTopList[i]
			b := stockTopList[j]
			return a.ChangeRate > b.ChangeRate
		})
		topStock := stockTopList[0]
		topCode = topStock.SecurityCode
		f10 := factors.GetL5F10(topCode)
		if f10 != nil {
			topName = f10.SecurityName
		}
		topRate = num.NetChangeRate(topStock.LastClose, topStock.Price)
		total := 0
		limits := 0
		ling := 0
		up := 0
		down := 0
		for j := 0; j < len(stockSnapshots); j++ {
			gp := stockSnapshots[j]
			total += 1
			zfLimit := exchange.MarketLimit(gp.SecurityCode)
			lastClose := num.Decimal(gp.LastClose)
			zhangting := num.Decimal(lastClose * float64(1.000+zfLimit))
			price := num.Decimal(gp.Price)
			if price >= zhangting {
				limits += 1
			}
			if price > lastClose {
				up++
			} else if price < lastClose {
				down++
			} else {
				ling += 1
			}
			gp.TopNo = j
			_, ok := __stock2Rank[gp.SecurityCode]
			if !ok {
				__stock2Rank[gp.SecurityCode] = gp
			}
		}
		for j, v := range allBlocks {
			if v.Code == block.Code {
				//allBlocks[j].Name = name
				allBlocks[j].TopCode = topCode
				allBlocks[j].TopName = topName
				allBlocks[j].TopRate = topRate
				//allBlocks[j].TopNo = j
				allBlocks[j].Count = total
				allBlocks[j].LimitUpNum = limits
				allBlocks[j].UpCount = up
				allBlocks[j].NoChangeNum = ling
				allBlocks[j].DownCount = down
				__mapBlockData[v.Code] = allBlocks[j]
			}
		}
	}

	// 输出 板块排行表格
	lastBlocks := api.Filter(allBlocks, tickSectorFilter)
	bn := len(lastBlocks)
	//if bn >= rules.RuleParameters.SectorsTopN {
	//	bn = rules.RuleParameters.SectorsTopN
	//}
	topBlocks := lastBlocks[:bn]
	blkTable := tablewriter.NewWriter(os.Stdout)
	blkHeaders := tags.GetHeadersByTags(SectorInfo{})
	blkTable.SetHeader(blkHeaders)
	var blkValues [][]string
	var targets []string
	for _, block := range topBlocks {
		values := tags.GetValuesByTags(block)
		blkValues = append(blkValues, values)
		// 是否从前排板块中检索个股
		//if rules.RuleParameters.SectorsFilter {
		targets = append(targets, block.StockCodes...)
		//}
	}
	blkTable.AppendBulk(blkValues)
	fmt.Println()
	blkTable.Render()
	targets = api.Unique(targets)
	return targets
}
