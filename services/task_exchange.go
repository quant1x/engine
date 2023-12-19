package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/features"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/runtime"
	"golang.org/x/exp/maps"
	"time"
)

// 更新快照
func jobUpdateExchangeAndSnapshot() {
	funcName, _, _ := runtime.Caller()
	now := time.Now()
	updateInRealTime, _ := trading.CanUpdateInRealtime()
	// 09:15:00~09:27:00, 14:57:00~15:01:00之间更新数据
	if updateInRealTime && trading.CheckCallAuctionTime(now) {
		realtimeUpdateExchangeAndSnapshot()
	} else {
		if runtime.Debug() {
			realtimeUpdateExchangeAndSnapshot()
		}
		logger.Infof("%s, 非集合竞价时段", funcName)
	}
}

var (
	snapshotDate = cache.DefaultCanReadDate()
	mapSnapshot  = map[string][]quotes.Snapshot{}
)

func resetSnapshotCache() {
	date := cache.DefaultCanReadDate()
	if date > snapshotDate {
		maps.Clear(mapSnapshot)
		snapshotDate = date
		cachel5.SwitchDate(snapshotDate)
	}
}

// realtimeUpdateExchangeAndSnapshot 更新快照缓存
func realtimeUpdateExchangeAndSnapshot() {
	resetSnapshotCache()
	moduleName := "执行[同步exchange]"
	logger.Infof("%s: begin", moduleName)
	allCodes := market.GetCodeList()
	count := len(allCodes)
	currentDate := cache.DefaultCanReadDate()
	logger.Infof("%s: begin-1", moduleName)
	bar := progressbar.NewBar(barIndexUpdateExchangeAndSnapshot, moduleName, count)
	for _, securityCode := range allCodes {
		bar.Add(1)
		if api.StartsWith(securityCode, []string{"88"}) {
			securityCode = "sh" + securityCode
		}
		//logger.Infof("%s: begin-1-2", moduleName)
		v := models.GetTickFromMemory(securityCode)
		if v == nil || v.Date != currentDate {
			// 如果snapshot缓存无效, 或者日期不是当前日期, 跳过
			continue
		}
		timestamp := time.Now()
		// 1. 修订日期
		v.Date = currentDate
		//cacheSnapshots := []flash.Exchange{}
		securityCode := v.SecurityCode
		exchange := features.GetL5Exchange(securityCode)
		if exchange == nil {
			exchange = &features.Exchange{
				Date: currentDate,
				Code: securityCode,
			}
		} else {
			exchange.Date = currentDate
		}
		// 2. 计算开盘和收盘的成交量
		exchange.OpenVolume = int64(v.OpenVolume)
		exchange.CloseVolume = int64(v.CloseVolume)
		// 计算开盘换手z和收盘换手z
		f10 := smart.GetL5F10(securityCode)
		if f10 != nil {
			exchange.OpenTurnZ = f10.TurnZ(exchange.OpenVolume)
			exchange.CloseTurnZ = f10.TurnZ(exchange.CloseVolume)
		}
		// 3. 计算快照扩展数据
		if trading.CheckCallAuctionOpen(timestamp) {
			// 3.1 早盘情绪有时效性
			// 计算早盘竞价方向
			exchange.OpenBiddingDirection, exchange.OpenVolumeDirection = v.CheckDirection()
			// 3.2 计算早盘情绪
			exchange.OpenSentiment, exchange.OpenConsistent = market.SnapshotSentiment(*v)
		} else {
			// 3.3盘 中及盘后的数据的计算都没有问题
			// 计算收盘竞价方向
			exchange.CloseBiddingDirection, exchange.CloseVolumeDirection = v.CheckDirection()
			// 3.4 计算收盘情绪
			exchange.CloseSentiment, exchange.CloseConsistent = market.SnapshotSentiment(*v)
		}
		// 4. 竞价上午竞价观察
		if trading.CheckCallAuctionOpen(timestamp) {
			// 4.1 竞价开盘
			if exchange.BidOpen == 0 {
				exchange.BidOpen = v.Ask1
			}
			// 4.2 竞价结束
			exchange.BidClose = v.Price
			// 4.3 竞价最高
			if exchange.BidHigh == 0 || exchange.BidHigh < v.Ask1 {
				exchange.BidHigh = v.Ask1
			}
			// 4.4 竞价最低
			if exchange.BidLow == 0 || exchange.BidLow > v.Ask1 {
				exchange.BidLow = v.Ask1
			}
			// 4.4 竞价匹配量
			exchange.BidMatched = float64(v.BidVol1)
			// 4.5 竞价未匹配量
			if v.BidVol2 == 0 {
				exchange.BidUnmatched = float64(v.AskVol2)
				exchange.BidDirection = -1
			}
			if v.AskVol2 == 0 {
				exchange.BidUnmatched = float64(v.BidVol2)
				exchange.BidDirection = 1
			}
		}
		// 5. 缓存数据
		//cacheSnapshots = append(cacheSnapshots, *exchange)

		// 6. 更新内存中的数据
		features.UpdateL5Exchange(exchange)
		// 7. 刷新缓存
		cacheList, ok := mapSnapshot[securityCode]
		if !ok {
			cacheList = []quotes.Snapshot{}
		}
		if len(cacheList) > 0 {
			lastDay := cacheList[len(cacheList)-1].Date
			lastServerTime := cacheList[len(cacheList)-1].ServerTime
			if currentDate == lastDay && v.ServerTime <= lastServerTime {
				// 时间戳在缓存之前, 忽略
				continue
			}
		}
		cacheList = append(cacheList, *v)
		if len(cacheList) > 0 {
			mapSnapshot[securityCode] = cacheList
		}
		//logger.Infof("%s: begin-1-3", moduleName)
	}
	logger.Infof("%s: begin-2", moduleName)
	// 刷新exchange快照本地cache
	features.RefreshL5Exchange()
	logger.Infof("%s: begin-3", moduleName)
	timestamp := time.Now()
	if trading.CheckCallAuctionOpenFinished(timestamp) || trading.CheckCallAuctionCloseFinished(timestamp) {
		// 早盘和尾盘集合竞价结束后刷新缓存文件
		for _, listSnapshot := range mapSnapshot {
			if len(listSnapshot) == 0 {
				continue
			}
			// 获取第一条记录
			first := listSnapshot[0]
			securityCode := first.SecurityCode
			filename := cache.SnapshotFilename(securityCode)
			cacheList := []quotes.Snapshot{}
			err := api.CsvToSlices(filename, &cacheList)
			if err == nil && len(cacheList) > 0 {
				// 缓存中最后一条记录
				rows := len(cacheList)
				last := cacheList[rows-1]
				// 日期
				lastDay := last.Date
				// 时间戳
				lastServerTime := last.ServerTime
				for _, v := range listSnapshot {
					if currentDate == lastDay && v.ServerTime <= lastServerTime {
						// 时间戳在缓存之前, 忽略
						continue
					}
					cacheList = append(cacheList, v)
				}
			} else {
				// 如果缓存文件不存在, 用缓存数据
				cacheList = listSnapshot
			}
			if len(cacheList) > 0 {
				_ = api.SlicesToCsv(filename, cacheList)
			}
		}
	}
	logger.Infof("%s: begin-4", moduleName)
	logger.Infof("%s: end", moduleName)
}