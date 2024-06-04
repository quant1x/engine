package tracker

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pkg/tablewriter"
	"os"
	"sort"
	"time"
)

var (
	__allStocks = []string{}
	//__topSectors = []SectorInfo{}
	// 缓存板块前三的个股
	__block2Top = map[string][]string{}
	// 缓存个股所在的板块列表
	__stock2Block = map[string][]SectorInfo{}
	// 缓存个股的排名
	__stock2Rank = map[string]factors.QuoteSnapshot{}
	// 个股对应板块
	__mapBlockData = map[string]SectorInfo{}
)

// AllScan 全市场扫描
func AllScan(barIndex *int, model models.Strategy, marketData market.MarketData) {
	today := exchange.IndexToday()
	dates := exchange.TradingDateRange(exchange.MARKET_CN_FIRST_DATE, today)
	days := len(dates)
	currentlyDay := dates[days-1]
	updateTime := "15:00:59"
	todayIsTradeDay := false
	if today == currentlyDay {
		todayIsTradeDay = true
		now := time.Now()
		nowTime := now.Format(exchange.CN_SERVERTIME_FORMAT)
		if nowTime < exchange.CN_TradingStartTime {
			currentlyDay = dates[days-2]
		} else if nowTime >= exchange.CN_TradingStartTime && nowTime <= exchange.CN_TradingStopTime {
			updateTime = now.Format(exchange.TimeOnly)
		}
	}
	// 全市场扫描
	tradeRule := config.GetStrategyParameterByCode(model.Code())
	if tradeRule == nil {
		logger.Errorf("strategy[%d]: trade rule not found", model.Code())
		return
	}
	var stockCodes []string
	needFilter := false
	// 检查板块扫描的结果是否存在股票列表
	if len(__allStocks) > 0 {
		// 确定前排板块中的前排个股
		__allStocks = api.Unique(__allStocks)
		stockCodes = __allStocks
		needFilter = true
	} else {
		// 板块扫描结果没有输出个股列表
		// 查看规则配置
		if tradeRule == nil {
			// 如果规则没有配置, 则取全部有效的代码列表
			stockCodes = marketData.GetStockCodeList()
			needFilter = true
		} else {
			// 如果规则配置有效, 股票代码列表从规则中获取
			stockCodes = marketData.GetStrategyStockCodeList(tradeRule)
		}
	}
	if needFilter && tradeRule != nil {
		stockCodes = tradeRule.Filter(stockCodes)
	}
	stockSnapshots := []factors.QuoteSnapshot{}
	stockCount := len(stockCodes)
	*barIndex++
	bar := progressbar.NewBar(*barIndex, "执行[全市场扫描]", stockCount)
	for start := 0; start < stockCount; start++ {
		bar.Add(1)
		code := stockCodes[start]
		securityCode := exchange.CorrectSecurityCode(code)
		if exchange.AssertIndexBySecurityCode(securityCode) {
			continue
		}
		v := models.GetTickFromMemory(securityCode)
		if v != nil {
			snapshot := models.QuoteSnapshotFromProtocol(*v)
			stockSnapshots = append(stockSnapshots, snapshot)
		}
	}
	if len(stockSnapshots) == 0 {
		return
	}
	if todayIsTradeDay && updateTime >= StoreSnapshotTimeBegin && updateTime < StoreSnapshotTimeEnd {
		var cacheSnapshots []factors.Misc
		for _, v := range stockSnapshots {
			securityCode := exchange.GetSecurityCode(v.Market, v.Code)
			snapshot := factors.GetL5Misc(securityCode)
			if snapshot == nil {
				snapshot = &factors.Misc{
					Date: currentlyDay,
					Code: securityCode,
				}
			}
			//snapshot.OpenBiddingDirection, snapshot.OpenVolumeDirection = v.CheckDirection()
			cacheSnapshots = append(cacheSnapshots, *snapshot)
			factors.UpdateL5Misc(snapshot)
		}
		if len(cacheSnapshots) > 0 {
			factors.RefreshL5Misc()
		}
	}
	// 过滤不符合条件的个股
	stockSnapshots = api.Filter(stockSnapshots, func(snapshot factors.QuoteSnapshot) bool {
		return model.Filter(tradeRule.Rules, snapshot) == nil
	})
	// 排序
	sortedStatus := model.Sort(stockSnapshots)
	if sortedStatus == models.SortDefault || sortedStatus == models.SortNotExecuted {
		sort.Slice(stockSnapshots, func(i, j int) bool {
			a := stockSnapshots[i]
			b := stockSnapshots[j]
			if a.OpenTurnZ > b.OpenTurnZ {
				return true
			}
			return a.OpenTurnZ == b.OpenTurnZ && a.OpeningChangeRate > b.OpeningChangeRate
		})
	}

	// 输出二维表格
	tbl := tablewriter.NewWriter(os.Stdout)
	tbl.SetHeader(tags.GetHeadersByTags(models.Statistics{}))
	topN := models.CountTopN
	if topN <= 0 {
		topN = tradeRule.Total
	}
	if topN > len(stockSnapshots) {
		topN = len(stockSnapshots)
	}
	votingResults := []models.Statistics{}
	stockSnapshots = stockSnapshots[:topN]
	now := time.Now()
	orderCreateTime := now.Format(cache.TimeStampMilli)
	for _, v := range stockSnapshots {
		ticket := models.Statistics{
			Date:                 currentlyDay,              // 日期
			Code:                 v.SecurityCode,            // 证券代码
			Name:                 v.Name,                    // 证券名称
			Active:               int(v.Active),             // 活跃度
			LastClose:            v.LastClose,               // 昨收
			Open:                 v.Open,                    // 开盘价
			OpenRaise:            v.OpeningChangeRate,       // 开盘涨幅
			Price:                v.Price,                   // 现价
			UpRate:               v.ChangeRate,              // 涨跌幅
			OpenPremiumRate:      v.PremiumRate,             // 集合竞价买入溢价率
			OpenVolume:           v.OpenVolume,              // 集合竞价-开盘量, 单位是股
			TurnZ:                v.OpenTurnZ,               // 开盘换手率z
			QuantityRatio:        v.OpenQuantityRatio,       // 开盘量比
			AveragePrice:         v.Amount / float64(v.Vol), // 均价线
			ChangePower:          v.ChangePower,             // 涨跌力度
			AverageBiddingVolume: v.AverageBiddingVolume,    // 委托均量
			UpdateTime:           orderCreateTime,           // 更新时间
		}
		if v.Open < v.LastClose {
			ticket.Tendency += "低开"
		} else if v.Open == v.LastClose {
			ticket.Tendency += "平开"
		} else {
			ticket.Tendency += "高开"
		}
		if ticket.AveragePrice < v.Open {
			ticket.Tendency += ",回落"
		} else {
			ticket.Tendency += ",拉升"
		}
		if v.Price > ticket.AveragePrice {
			ticket.Tendency += ",强势"
		} else {
			ticket.Tendency += ",弱势"
		}

		bs, ok := __stock2Block[ticket.Code]
		if ok {
			tb := bs[0]
			if block, ok := __mapBlockData[tb.Code]; ok {
				ticket.BlockName = block.Name
				ticket.BlockRate = block.ChangeRate
				ticket.BlockTop = block.Rank
				shot, ok1 := __stock2Rank[ticket.Code]
				if ok1 {
					ticket.BlockRank = shot.TopNo
				}
			}
		}
		votingResults = append(votingResults, ticket)
	}
	gtP1 := 0 // 存在溢价
	gtP2 := 0 // 超过1%
	gtP3 := 0
	gtP4 := 0
	gtP5 := 0
	yields := 0.00
	for _, v := range votingResults {
		rate := num.NetChangeRate(v.Open, v.Price)
		if rate > 0 {
			gtP1 += 1
		}
		if rate >= 1.00 {
			gtP2 += 1
		}
		if rate >= 2.00 {
			gtP3 += 1
		}
		if rate >= 3.00 {
			gtP4 += 1
		}
		if rate >= 5.00 {
			gtP5 += 1
		}
		yields += rate
		tbl.Append(tags.GetValuesByTags(v))
	}
	yields /= float64(len(votingResults))
	fmt.Println() // 输出一个换行
	tbl.Render()
	count := len(stockSnapshots)
	fmt.Println()

	fmt.Println(currentlyDay + " " + updateTime + ", 胜率统计:")
	fmt.Printf("\t==> 胜    率: %d/%d, %.2f%%, 收益率: %.2f%%\n", gtP1, count, 100*float64(gtP1)/float64(count), yields)
	fmt.Printf("\t==> 溢价超1%%: %d/%d, %.2f%%\n", gtP2, count, 100*float64(gtP2)/float64(count))
	fmt.Printf("\t==> 溢价超2%%: %d/%d, %.2f%%\n", gtP3, count, 100*float64(gtP3)/float64(count))
	fmt.Printf("\t==> 溢价超3%%: %d/%d, %.2f%%\n", gtP4, count, 100*float64(gtP4)/float64(count))
	fmt.Printf("\t==> 溢价超5%%: %d/%d, %.2f%%\n", gtP5, count, 100*float64(gtP5)/float64(count))
	fmt.Println()
	// 存储
	storages.OutputStatistics(model, currentlyDay, votingResults)
}
