package tracker

import (
	"fmt"
	"os"
	"sort"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pkg/tablewriter"
)

// GoodCase good case
type GoodCase struct {
	Date   string  `dataframe:"日期"`
	Num    int     `dataframe:"数量"`
	Yields float64 `dataframe:"浮动收益率%"`
	//NextYields float64 `dataframe:"隔日收益率%"`
	GtP1 float64 `dataframe:"胜率率%"`
	GtP2 float64 `dataframe:"溢价超1%"`
	GtP3 float64 `dataframe:"溢价超2%"`
	GtP4 float64 `dataframe:"溢价超3%"`
	GtP5 float64 `dataframe:"溢价超5%"`
}

// SampleFeature 样本特征
type SampleFeature struct {
	SecurityCode      string
	Name              string
	OpenChangeRate    float64
	OpenTurnZ         float64
	LastClose         float64
	Open              float64
	Price             float64
	UpRate            float64
	OpenPremiumRate   float64
	NextPremiumRate   float64
	OpenQuantityRatio float64 // 量比
	Beta              float64
	Alpha             float64
}

func checkWideOffset(klines []factors.SecurityFeature, date string) (offset int) {
	rows := len(klines)
	offset = 0
	for i := 0; i < rows; i++ {
		klineDate := klines[rows-1-i].Date
		if klineDate < date {
			return -1
		} else if klineDate == date {
			break
		} else {
			offset++
		}
	}
	if offset+1 >= rows {
		return -1
	}
	return
}

// BackTesting 执行策略回测，计算指定策略在指定日期范围内的表现
//
// 参数:
//
//	strategyNo: 策略编号
//	countDays: 回测天数
//	countTopN: 每日选取的标的数量
//
// 功能:
//  1. 获取指定日期范围内的交易日数据
//  2. 加载策略模型和交易规则
//  3. 对每个交易日执行以下操作:
//     - 加载标的特征数据
//     - 过滤符合条件的标的
//     - 计算各项指标(Alpha/Beta/溢价率等)
//     - 输出每日回测结果
//  4. 汇总统计整体回测表现
//
// 输出:
//  1. 控制台输出每日回测结果表格
//  2. 生成CSV文件保存详细回测数据
//  3. 计算并输出平均收益率、胜率等汇总指标
func BackTesting(strategyNo uint64, countDays, countTopN int) {
	currentlyDay := exchange.GetCurrentlyDay()
	dates := exchange.TradingDateRange(exchange.MARKET_CH_FIRST_LISTTIME, currentlyDay)
	scope := api.RangeFinite(-countDays)
	s, e, err := scope.Limits(len(dates))
	if err != nil {
		fmt.Println(err)
		return
	}
	model, err := models.CheckoutStrategy(strategyNo)
	if err != nil {
		fmt.Println(err)
		return
	}
	//TODO: 这里应该要取策略的规则参数
	tradeRule := config.GetStrategyParameterByCode(strategyNo)
	if tradeRule == nil {
		return
	}
	backTestingParameter := config.GetDataConfig().BackTesting
	var allResult []models.Statistics
	var gcs []GoodCase
	dates = dates[s : e+1]
	codes := market.GetCodeList()
	mapStock := map[string][]factors.SecurityFeature{}
	for _, date := range dates {
		testDate := date
		// 切换策略数据的缓存日期
		factors.SwitchDate(testDate)
		var marketPrices []float64
		var stockSnapshots []factors.QuoteSnapshot
		total := len(codes)
		//pos := 0
		bar := progressbar.NewBar(1, "执行["+testDate+"涨幅扫描]", total)
		for _, securityCode := range codes {
			bar.Add(1)
			if !exchange.AssertStockBySecurityCode(securityCode) && securityCode != backTestingParameter.TargetIndex {
				continue
			}
			//features := factors.CheckoutWideTableByDate(securityCode, date)
			features, ok := mapStock[securityCode]
			if !ok {
				filename := cache.WideFilename(securityCode)
				err := api.CsvToSlices(filename, &features)
				if err != nil {
					continue
				}
				mapStock[securityCode] = features
			}
			if securityCode == backTestingParameter.TargetIndex && len(marketPrices) == 0 {
				for _, m := range features {
					marketPrices = append(marketPrices, m.ChangeRate)
				}
			}
			length := len(features)
			offset := checkWideOffset(features, testDate)
			if offset < 0 {
				continue
			}

			// 索引防御与窗口计算（基于 feature 在宽表中的位置）
			idx := length - 1 - offset // 当前 testDate 在 features 中的下标
			if idx < 0 || idx >= length {
				continue
			}
			start := idx - (countDays - 1)
			if start < 0 {
				start = 0
			}
			windowLen := idx - start + 1
			if windowLen <= 0 {
				continue
			}

			// 如果还没有基准序列，懒加载目标指数（只加载 targetIndex，不会加载其它无关标的）
			if len(marketPrices) == 0 && backTestingParameter.TargetIndex != "" {
				ti := backTestingParameter.TargetIndex
				if feats, ok := mapStock[ti]; ok && len(feats) > 0 {
					for _, m := range feats {
						marketPrices = append(marketPrices, m.ChangeRate)
					}
				} else {
					var feats2 []factors.SecurityFeature
					fn := cache.WideFilename(ti)
					if err := api.CsvToSlices(fn, &feats2); err == nil && len(feats2) > 0 {
						mapStock[ti] = feats2
						for _, m := range feats2 {
							marketPrices = append(marketPrices, m.ChangeRate)
						}
					}
				}
			}

			// 基准数据长度必须满足窗口长度
			if len(marketPrices) < windowLen {
				// 基准不足，跳过该标的
				continue
			}

			markets := marketPrices[:windowLen]
			prices := make([]float64, windowLen)
			for si := 0; si < windowLen; si++ {
				prices[si] = features[start+si].ChangeRate
			}

			feature := features[idx]
			// 宽表和测试日期没有对齐, 跳过
			if feature.Date != testDate {
				// 停牌导致的日期无法从后往前对齐
				continue
			}
			snapshot := models.FeatureToSnapshot(feature, securityCode)
			// 下一个交易日开盘价
			diffDays := 1
			nextOffset := length - offset - 1 + diffDays
			if nextOffset < length {
				nextFeature := features[nextOffset]
				snapshot.NextOpen = nextFeature.Open
				snapshot.NextClose = nextFeature.Close
				snapshot.NextHigh = nextFeature.High
				snapshot.NextLow = nextFeature.Low
			}
			snapshot.Beta, snapshot.Alpha = exchange.EvaluateYields(prices, markets, config.TraderConfig().DailyRiskFreeRate(testDate))
			snapshot.Beta *= 100
			snapshot.Alpha *= 100
			stockSnapshots = append(stockSnapshots, snapshot)
		}
		bar.Wait()
		if len(stockSnapshots) == 0 {
			continue
		}

		// 过滤不符合条件的个股
		stockSnapshots = api.Filter(stockSnapshots, func(snapshot factors.QuoteSnapshot) bool {
			err := model.Filter(tradeRule.Rules, snapshot)
			//if snapshot.SecurityCode == "sz300956" {
			//	fmt.Printf("%+v, err=%v\n", snapshot, err)
			//	//return true
			//}
			return err == nil
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

		var samples []SampleFeature
		for _, snapshot := range stockSnapshots {
			securityCode := snapshot.SecurityCode
			// 获取证券名称
			securityName := "unknown"
			f10 := factors.GetL5F10(securityCode, testDate)
			if f10 != nil {
				securityName = f10.SecurityName
			}

			sample := SampleFeature{
				Name:              securityName,
				SecurityCode:      securityCode,
				OpenQuantityRatio: snapshot.OpenQuantityRatio,
				OpenTurnZ:         snapshot.OpenTurnZ,
				OpenChangeRate:    num.NetChangeRate(snapshot.LastClose, snapshot.Open),
				LastClose:         snapshot.LastClose,
				Open:              snapshot.Open,
				Price:             snapshot.Price,
				UpRate:            num.NetChangeRate(snapshot.LastClose, snapshot.Price),
				OpenPremiumRate:   num.NetChangeRate(snapshot.Open, snapshot.Price),
				NextPremiumRate:   num.NetChangeRate(snapshot.Open, snapshot.NextOpen),
			}
			switch tradeRule.Flag {
			case models.OrderFlagHead:
				sample.OpenPremiumRate = num.NetChangeRate(snapshot.Open, snapshot.Price)
				sample.NextPremiumRate = num.NetChangeRate(snapshot.Open, snapshot.NextOpen)
			case models.OrderFlagTail:
				sample.OpenPremiumRate = num.NetChangeRate(snapshot.Price, snapshot.Price)
				sample.NextPremiumRate = num.NetChangeRate(snapshot.Price, snapshot.NextClose)
				if snapshot.Price < snapshot.NextClose && snapshot.Price*(1+backTestingParameter.NextPremiumRate+0.005) < snapshot.NextHigh {
					sample.NextPremiumRate = num.NetChangeRate(snapshot.Price, snapshot.Price*(1+backTestingParameter.NextPremiumRate))
				}
			case models.OrderFlagTick:
				sample.OpenPremiumRate = num.NetChangeRate(snapshot.Price, snapshot.Price)
				sample.NextPremiumRate = num.NetChangeRate(snapshot.Price, snapshot.NextClose)
			}
			sample.Beta = snapshot.Beta
			sample.Alpha = snapshot.Alpha
			samples = append(samples, sample)
		}

		// 单日回测结果
		// 检查有效记录最大数
		topN := countTopN
		if topN > len(samples) {
			topN = len(samples)
		}

		tbl := tablewriter.NewWriter(os.Stdout)
		tbl.SetHeader(tags.GetHeadersByTags(models.Statistics{}))
		samples = samples[:topN]
		var results []models.Statistics
		for _, v := range samples {
			zs := models.Statistics{
				Date:            testDate,            // 日期
				Code:            v.SecurityCode,      // 证券代码
				Name:            v.Name,              // 证券名称
				OpenRaise:       v.OpenChangeRate,    // 开盘涨幅
				TurnZ:           v.OpenTurnZ,         // 开盘换手率z
				QuantityRatio:   v.OpenQuantityRatio, // 开盘量比
				LastClose:       v.LastClose,         // 昨日收盘
				Open:            v.Open,              // 开盘价
				Price:           v.Price,             // 现价
				UpRate:          v.UpRate,            // 涨跌幅
				OpenPremiumRate: v.OpenPremiumRate,   // 集合竞价买入, 溢价率
				NextPremiumRate: v.NextPremiumRate,   // 隔日溢价率
				Beta:            v.Beta,
				Alpha:           v.Alpha,
			}
			switch tradeRule.Flag {
			case models.OrderFlagHead:
				zs.UpdateTime = zs.Date + " 09:27:10.000"
			case models.OrderFlagTail:
				zs.UpdateTime = zs.Date + " 14:56:10.000"
			case models.OrderFlagTick:
				zs.UpdateTime = zs.Date + " 14:56:10.000"
			}

			results = append(results, zs)
		}
		gtP1 := 0 // 存在溢价
		gtP2 := 0 // 超过1%
		gtP3 := 0 // 超过2%
		gtP4 := 0 // 超过3%
		gtP5 := 0 // 超过5%
		yields := 0.00
		for _, v := range results {
			rate := v.NextPremiumRate
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
		yields /= float64(len(results))
		fmt.Println() // 输出一个换行
		tbl.Render()
		count := len(samples)
		gc := GoodCase{
			Date:   testDate,
			Num:    count,
			Yields: yields,
			GtP1:   100 * float64(gtP1) / float64(count),
			GtP2:   100 * float64(gtP2) / float64(count),
			GtP3:   100 * float64(gtP3) / float64(count),
			GtP4:   100 * float64(gtP4) / float64(count),
			GtP5:   100 * float64(gtP5) / float64(count),
		}
		if num.IsNaN(gc.Yields) {
			gc.Yields = 0
		}
		if num.IsNaN(gc.GtP1) {
			gc.GtP1 = 0
		}
		gcs = append(gcs, gc)
		fmt.Println(testDate + ", 胜率统计:")
		fmt.Printf("\t==> 胜    率: %d/%d, %.2f%%, 收益率: %.2f%%\n", gtP1, count, 100*float64(gtP1)/float64(count), yields)
		fmt.Printf("\t==> 溢价超1%%: %d/%d, %.2f%%\n", gtP2, count, 100*float64(gtP2)/float64(count))
		fmt.Printf("\t==> 溢价超2%%: %d/%d, %.2f%%\n", gtP3, count, 100*float64(gtP3)/float64(count))
		fmt.Printf("\t==> 溢价超3%%: %d/%d, %.2f%%\n", gtP4, count, 100*float64(gtP4)/float64(count))
		fmt.Printf("\t==> 溢价超5%%: %d/%d, %.2f%%\n", gtP5, count, 100*float64(gtP5)/float64(count))
		fmt.Println()
		allResult = append(allResult, results...)
		//storages.OutputStatistics("tracker", topN, date, results)
	}

	// 合计输出
	total := len(dates)
	fmt.Printf("\n策略编号: %d, 策略名称: %s, 订单类型: %s\n", model.Code(), model.Name(), tradeRule.Flag)
	fmt.Printf("%s - %s 合计: %d 个交易日\n", dates[0], dates[len(dates)-1], total)
	today := cache.Today()
	dfTotal := pandas.LoadStructs(gcs)
	if dfTotal.Nrow() > 0 {
		winningRate := dfTotal.Col("浮动收益率%").FillNa(0, true).Mean()
		winningAverage := dfTotal.Col("胜率率%").FillNa(0, true).Mean()
		fmt.Printf("\t==> 平均 浮动溢价率:%.4f%%, 平均 胜率率: %.4f%%\n", winningRate, winningAverage)
		filename := fmt.Sprintf("%s/total-%s-%s-%d.csv", storages.GetResultCachePath(), tradeRule.QmtStrategyName(), today, countTopN)
		_ = dfTotal.WriteCSV(filename)
	}
	dfRecords := pandas.LoadStructs(allResult)
	if dfRecords.Nrow() > 0 {
		fudu := dfRecords.Col("open_premium_rate").FillNa(0, true).Mean()
		geri := dfRecords.Col("next_premium_rate").FillNa(0, true).Mean()
		fmt.Printf("\t==> 平均 浮动溢价率:%.4f%%, 平均 隔日溢价率: %.4f%%\n", fudu, geri)
		colNames := tags.GetHeadersByTags(allResult[0])
		_ = dfRecords.SetNames(colNames...)
		filename := fmt.Sprintf("%s/backtesting-%s-%s-%d.csv", storages.GetResultCachePath(), tradeRule.QmtStrategyName(), today, countTopN)
		_ = dfRecords.WriteCSV(filename)
	}
	var winningRate []float64
	var winningAverage []float64
	winningCount := 0
	for _, gc := range gcs {
		if gc.Num < 1 {
			continue
		}
		winningCount++
		if !num.IsNaN(gc.Yields) {
			winningRate = append(winningRate, gc.Yields)
		}
		if !num.IsNaN(gc.GtP1) {
			winningAverage = append(winningAverage, gc.GtP1)
		}
	}
	if winningCount > 0 {
		fmt.Printf("\n")
		fmt.Printf("\t==> 扣除未交易后: %d 个交易日, 策略覆盖交易日率: %d/%d = %.4f%%\n", winningCount, winningCount, total, 100*num.ChangeRate(total, winningCount))
		fmt.Printf("\t==> 平均 浮动溢价率:%.4f%%, 平均 胜率率: %.4f%%\n", num.Sum(winningRate)/float64(winningCount), num.Sum(winningAverage)/float64(winningCount))
	}

	var fudong []float64
	var geri []float64
	for _, result := range allResult {
		if !num.IsNaN(result.OpenPremiumRate) {
			fudong = append(fudong, result.OpenPremiumRate)
		}
		if !num.IsNaN(result.NextPremiumRate) {
			geri = append(geri, result.NextPremiumRate)
		}
	}
	fmt.Printf("\t==> 平均 浮动溢价率:%.4f%%, 平均 隔日溢价率: %.4f%%\n", num.Mean(fudong), num.Mean(geri))
}
