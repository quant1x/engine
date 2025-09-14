package tracker

import (
	"fmt"
	"os"
	"time"

	"github.com/quant1x/engine/factors"
	"github.com/quant1x/engine/models"
	"github.com/quant1x/engine/storages"
	"github.com/quant1x/exchange"
	"github.com/quant1x/num"
	"github.com/quant1x/pkg/tablewriter"
	"github.com/quant1x/x/tags"
)

// OutputTable 输出表格
func OutputTable(model models.Strategy, stockSnapshots []factors.QuoteSnapshot) {
	today := exchange.IndexToday()
	dates := exchange.TradeRange(exchange.MARKET_CN_FIRST_DATE, today)
	days := len(dates)
	currentlyDay := dates[days-1]
	updateTime := "15:00:59"
	//todayIsTradeDay := false
	if today == currentlyDay {
		//todayIsTradeDay = true
		now := time.Now()
		nowTime := now.Format(exchange.CN_SERVERTIME_FORMAT)
		if nowTime < exchange.CN_TradingStartTime {
			currentlyDay = dates[days-2]
		} else if nowTime >= exchange.CN_TradingStartTime && nowTime <= exchange.CN_TradingStopTime {
			updateTime = now.Format(exchange.TimeOnly)
		}
	}
	// 控制台输出表格
	tbl := tablewriter.NewWriter(os.Stdout)
	tbl.SetHeader(tags.GetHeadersByTags(models.Statistics{}))
	orderCount := models.CountTopN
	if orderCount > len(stockSnapshots) {
		orderCount = len(stockSnapshots)
	}
	votingResults := []models.Statistics{}
	//stockSnapshots = stockSnapshots[:orderCount]
	orderCreateTime := factors.GetTimestamp()
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
			AveragePrice:         v.Amount / float64(v.Vol), // 均价
			Speed:                v.Rate,                    // 涨速
			ChangePower:          v.ChangePower,             // 涨跌力度
			AverageBiddingVolume: v.AverageBiddingVolume,    // 委比
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
