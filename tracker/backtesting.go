package tracker

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/engine/strategies"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/num"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pkg/tablewriter"
	"os"
	"sort"
)

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
}

// BackTesting 回测
func BackTesting(countDays, countTopN int) {
	currentlyDay := trading.GetCurrentlyDay()
	dates := trading.TradeRange(proto.MARKET_CH_FIRST_LISTTIME, currentlyDay)
	scope := api.RangeFinite(-countDays)
	s, e, err := scope.Limits(len(dates))
	if err != nil {
		fmt.Println(err)
		return
	}
	allResult := []models.Statistics{}
	gcs := []GoodCase{}
	dates = dates[s : e+1]
	codes := market.GetCodeList()
	mapStock := map[string][]cache.SecurityFeature{}
	for i, date := range dates {
		// 切换策略数据的缓存日期
		cachel5.SwitchDate(date)
		samples := []SampleFeature{}
		total := len(codes)
		bar := progressbar.NewBar(1, "执行["+date+"涨幅扫描]", total)
		for _, securityCode := range codes {
			bar.Add(1)
			if !proto.AssertStockBySecurityCode(securityCode) {
				continue
			}
			features, ok := mapStock[securityCode]
			if !ok {
				filename := cache.FeatureFilename(securityCode)
				err := api.CsvToSlices(filename, &features)
				if err != nil {
					continue
				}
				mapStock[securityCode] = features
			}
			length := len(features)
			pos := length - countDays + i
			if pos < 0 {
				continue
			}
			feature := features[pos]
			snapshot := models.FeatureToSnapshot(feature, securityCode)
			if !strategies.RuleFilter(snapshot) {
				continue
			}

			// 获取证券名称
			securityName := "unknown"
			f10 := smart.GetL5F10(securityCode)
			if f10 != nil {
				securityName = f10.SecurityName
			}

			// 下一个交易日开盘价
			nextOpen := feature.Close
			diffDays := 1
			if pos+diffDays < length {
				nextFeature := features[pos+diffDays]
				nextOpen = nextFeature.Open
			}

			turn := SampleFeature{
				Name:              securityName,
				SecurityCode:      securityCode,
				OpenQuantityRatio: snapshot.OpenQuantityRatio,
				OpenTurnZ:         feature.OpenTurnZ,
				OpenChangeRate:    num.NetChangeRate(feature.LastClose, feature.Open),
				LastClose:         feature.LastClose,
				Open:              feature.Open,
				Price:             feature.Close,
				UpRate:            num.NetChangeRate(feature.LastClose, feature.Close),
				OpenPremiumRate:   num.NetChangeRate(feature.Open, feature.Close),
				NextPremiumRate:   num.NetChangeRate(feature.Open, nextOpen),
			}
			//basicInfo, err := security.GetBasicInfo(securityCode)
			//if err == nil && basicInfo != nil {
			//	turn.Name = basicInfo.Name
			//}
			samples = append(samples, turn)
		}
		sort.Slice(samples, func(i, j int) bool {
			a := samples[i]
			b := samples[j]
			if a.OpenTurnZ > b.OpenTurnZ {
				return true
			}
			return a.OpenTurnZ == b.OpenTurnZ && a.OpenChangeRate > b.OpenChangeRate
		})

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
				Date:            date,                // 日期
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
			}
			results = append(results, zs)
		}
		gtP1 := 0 // 存在溢价
		gtP2 := 0 // 超过1%
		gtP3 := 0 // 超过2%
		gtP4 := 0 // 超过3%
		gtP5 := 0 // 超过5%
		yields := float64(0.00)
		for _, v := range results {
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
		yields /= float64(len(results))
		fmt.Println() // 输出一个换行
		tbl.Render()
		count := len(samples)
		gc := GoodCase{
			Date:   date,
			Num:    count,
			Yields: yields,
			GtP1:   100 * float64(gtP1) / float64(count),
			GtP2:   100 * float64(gtP2) / float64(count),
			GtP3:   100 * float64(gtP3) / float64(count),
			GtP4:   100 * float64(gtP4) / float64(count),
			GtP5:   100 * float64(gtP5) / float64(count),
		}
		gcs = append(gcs, gc)
		fmt.Println(date + ", 胜率统计:")
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
	fmt.Printf("%s - %s 合计:\n", dates[0], dates[len(dates)-1])
	today := cache.Today()
	dfTotal := pandas.LoadStructs(gcs)
	if dfTotal.Nrow() > 0 {
		winningRate := dfTotal.Col("浮动收益率%").Mean()
		winningAverage := dfTotal.Col("胜率率%").Mean()
		fmt.Printf("\t==> 平均 浮动溢价率:%.4f%%, 平均 胜率率: %.4f%%\n", winningRate, winningAverage)
		filename := fmt.Sprintf("%s/total-%s-%d.csv", storages.GetResultCachePath(), today, countTopN)
		_ = dfTotal.WriteCSV(filename)
	}
	dfRecords := pandas.LoadStructs(allResult)
	if dfRecords.Nrow() > 0 {
		fudu := dfRecords.Col("open_premium_rate").Mean()
		geri := dfRecords.Col("next_premium_rate").Mean()
		fmt.Printf("\t==> 平均 浮动溢价率:%.4f%%, 平均 隔日溢价率: %.4f%%\n", fudu, geri)
		colNames := tags.GetHeadersByTags(allResult[0])
		_ = dfRecords.SetNames(colNames...)
		filename := fmt.Sprintf("%s/backtesting-%s-%d.csv", storages.GetResultCachePath(), today, countTopN)
		_ = dfRecords.WriteCSV(filename)
	}
	//fmt.Println("\n")
}
