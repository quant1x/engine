package tracker

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/concurrent"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pkg/tablewriter"
	"os"
	"slices"
	"sort"
	"sync"
	"time"
)

// ScanAllSectors 扫描板块
func ScanAllSectors(barIndex *int, model models.Strategy) {
	tradeRule := config.GetStrategyParameterByCode(model.Code())
	if tradeRule == nil {
		logger.Errorf("strategy[%d]: trade rule not found", model.Code())
		return
	}
	// 执行板块指数的检测
	typeBlocks := TopBlockWithType(barIndex, tradeRule)
	// 不分板块类型, 所有的板块放在一起排序
	allBlocks := []SectorInfo{}
	for _, v := range typeBlocks {
		allBlocks = append(allBlocks, v...)
	}
	// 扫描板块内个股排名
	blockCount := len(allBlocks)
	fmt.Println()
	bar := progressbar.NewBar(*barIndex, "执行[板块个股涨幅扫描]", blockCount)
	for i := 0; i < blockCount; i++ {
		bar.Add(1)
		block := &allBlocks[i]
		//blockCode := block.Code
		topCode := SecurityUnknown
		topName := SecurityUnknown
		topRate := float64(0.00)
		stockCodes := block.StockCodes
		stockSnapshots := []factors.QuoteSnapshot{}
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
		for j := 0; j < len(stockSnapshots) && j < tradeRule.Rules.StockTopNInSector; j++ {
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
				allBlocks[j].TopCode = topCode
				allBlocks[j].TopName = topName
				allBlocks[j].TopRate = topRate
				allBlocks[j].Count = total
				allBlocks[j].LimitUpNum = limits
				allBlocks[j].UpCount = up
				allBlocks[j].NoChangeNum = ling
				allBlocks[j].DownCount = down
				__mapBlockData[v.Code] = allBlocks[j]
			}
		}
	}
	// TODO 板块再排序
	// 输出 板块排行表格
	var lastBlocks []SectorInfo
	isHead := tradeRule.Flag == models.OrderFlagHead
	if isHead {
		lastBlocks = api.Filter(allBlocks, sectorFilterForHead)
	} else {
		lastBlocks = api.Filter(allBlocks, sectorFilterForTick)
	}
	bn := len(lastBlocks)
	if bn >= tradeRule.Rules.SectorsTopN {
		bn = tradeRule.Rules.SectorsTopN
	}
	topBlocks := lastBlocks[:bn]
	blkTable := tablewriter.NewWriter(os.Stdout)
	blkHeaders := tags.GetHeadersByTags(SectorInfo{})
	blkTable.SetHeader(blkHeaders)
	blkValues := [][]string{}
	for _, block := range topBlocks {
		values := tags.GetValuesByTags(block)
		blkValues = append(blkValues, values)
		// 是否从前排板块中检索个股
		if tradeRule.Rules.SectorsFilter {
			__allStocks = append(__allStocks, block.StockCodes...)
		}
	}
	blkTable.AppendBulk(blkValues)
	fmt.Println()
	blkTable.Render()

	*barIndex++
	// 执行策略
	// 获取全部证券代码
	stockCodes := []string{}
	for bkn, v := range topBlocks {
		bc := v.Code
		sl, ok := __block2Top[bc]
		if ok {
			stockCodes = append(stockCodes, sl...)
		}
		_ = bkn
	}
	stockCodes = api.Unique(stockCodes)

	count := len(stockCodes)
	fmt.Println("")
	bar = progressbar.NewBar(*barIndex, "执行["+model.Name()+"]", count)
	*barIndex++
	mapStock := concurrent.NewTreeMap[string, models.ResultInfo]()
	mainStart := time.Now()
	var wg = sync.WaitGroup{}
	for i, v := range stockCodes {
		securityCode := v
		bar.Add(1)
		wg.Add(1)
		go evaluate(model, &wg, securityCode, mapStock)
		_ = i
	}
	wg.Wait()
	fmt.Println()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tags.GetHeadersByTags(models.ResultInfo{}))

	elapsedTime := time.Since(mainStart) / time.Millisecond
	goals := mapStock.Size()
	message := fmt.Sprintf("总耗时: %.3fs, 总记录: %d, 命中: %d, 平均: %.3f/s\n", float64(elapsedTime)/1000, count, goals, float64(count)/(float64(elapsedTime)/1000))
	fmt.Printf(message)
	logger.Infof(message)
	fmt.Println("")
	// 执行曲线回归
	wg = sync.WaitGroup{}
	bar = progressbar.NewBar(*barIndex, "执行[综合策略]", goals)
	*barIndex++
	rs := make([]models.ResultInfo, 0)
	mapStock.Each(func(key string, value models.ResultInfo) {
		bar.Add(1)
		row := value
		stockCode := row.Code
		bs, ok := __stock2Block[stockCode]
		if ok {
			tb := bs[0]
			if block, ok1 := __mapBlockData[tb.Code]; ok1 {
				row.BlockType = block.Type
				row.BlockName = block.Name
				row.BlockRate = block.ChangeRate
				row.BlockTop = block.Rank
				row.BlockZhangTing = fmt.Sprintf("%d/%d", block.LimitUpNum, block.Count)
				row.BlockDescribe = fmt.Sprintf("%d/%d/%d", block.UpCount, block.DownCount, block.NoChangeNum)
				row.BlockTopName = block.TopName
				row.BlockTopRate = block.TopRate
				shot, ok1 := __stock2Rank[stockCode]
				if ok1 {
					row.BlockRank = shot.TopNo
				}
			}
		}
		predict := func(info models.ResultInfo, rs *[]models.ResultInfo, tbl *tablewriter.Table) {
			defer wg.Done()
			wg.Add(1)
			info.Predict()
			*rs = append(*rs, info)
			tbl.Append(tags.GetValuesByTags(info))
		}
		predict(row, &rs, table)
	})
	wg.Wait()
	fmt.Println("")
	output(model.Code(), rs)
	table.Render()
}

func output(strategyNo uint64, v []models.ResultInfo) {
	df := pandas.LoadStructs(v)
	filename := fmt.Sprintf("%s/%s/%s-%d.csv", cache.GetRootPath(), models.CACHE_STRATEGY_PATH, cache.Today(), strategyNo)
	_ = df.WriteCSV(filename)
}

// 个股评估
func evaluate(api models.Strategy, wg *sync.WaitGroup, code string, result *concurrent.TreeMap[string, models.ResultInfo]) {
	defer wg.Done()
	api.Evaluate(code, result)
}
