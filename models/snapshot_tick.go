package models

import (
	"sync"

	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/num"
)

var (
	__mutexTicks sync.RWMutex
	__cacheTicks = map[string]quotes.Snapshot{}
)

// GetTickFromMemory 获取快照缓存
func GetTickFromMemory(securityCode string) *quotes.Snapshot {
	__mutexTicks.RLock()
	v, found := __cacheTicks[securityCode]
	__mutexTicks.RUnlock()
	if found {
		return &v
	}
	return nil
}

// GetStrategySnapshot 从缓存中获取快照
func GetStrategySnapshot(securityCode string) *factors.QuoteSnapshot {
	v := GetTickFromMemory(securityCode)
	if v == nil || v.State != quotes.TDX_SECURITY_TRADE_STATE_NORMAL {
		// 非正常交易的记录忽略掉
		return nil
	}
	snapshot := factors.QuoteSnapshot{}
	_ = api.Copy(&snapshot, &v)
	snapshot.Name = securities.GetStockName(securityCode)
	//snapshot.Code = securityCode
	snapshot.OpeningChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Open)
	snapshot.ChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Price)
	f10 := factors.GetL5F10(securityCode)
	if f10 != nil {
		snapshot.Capital = f10.Capital
		snapshot.FreeCapital = f10.FreeCapital
		snapshot.OpenTurnZ = f10.TurnZ(snapshot.OpenVolume)
	}
	history := factors.GetL5History(securityCode)
	if history != nil {
		lastMinuteVolume := history.GetMV5()
		snapshot.OpenQuantityRatio = float64(snapshot.OpenVolume) / lastMinuteVolume
		minuteVolume := float64(snapshot.Vol) / float64(exchange.Minutes(snapshot.Date))
		snapshot.QuantityRatio = minuteVolume / lastMinuteVolume
	}
	snapshot.OpenBiddingDirection, snapshot.OpenVolumeDirection = v.CheckDirection()
	return &snapshot
}

// SyncAllSnapshots 实时更新快照
func SyncAllSnapshots(barIndex *int) {
	modName := "同步快照数据"
	allCodes := securities.AllCodeList()
	count := len(allCodes)
	var bar *progressbar.Bar = nil
	if barIndex != nil {
		bar = progressbar.NewBar(*barIndex, "执行["+modName+"]", count)
	}
	currentDate := exchange.GetCurrentlyDay()
	tdxApi := gotdx.GetTdxApi()
	// 读取配置的并发数
	parallelCount := config.GetDataConfig().Snapshot.Concurrency
	if parallelCount < 1 {
		parallelCount = tdxApi.NumOfServers()
		parallelCount /= 2
		if parallelCount < config.DefaultMinimumConcurrencyForSnapshots {
			parallelCount = config.DefaultMinimumConcurrencyForSnapshots
		}
	}
	var snapshots []quotes.Snapshot
	var wg sync.WaitGroup
	var mutex sync.Mutex
	codeCh := make(chan []string, parallelCount)

	// 启动goroutine来处理快照获取
	for i := 0; i < parallelCount; i++ {
		go func() {
			for subCodes := range codeCh {
				for i := 0; i < quotes.DefaultRetryTimes; i++ {
					list, err := tdxApi.GetSnapshot(subCodes)
					if err != nil {
						logger.Errorf("ZS: 网络异常: %+v, 重试: %d", err, i+1)
						continue
					}
					mutex.Lock()
					for _, v := range list {
						// 修订日期
						v.Date = currentDate
						snapshots = append(snapshots, v)
					}
					mutex.Unlock()

					break
				}
			}
			wg.Done()
		}()
	}

	for start := 0; start < count; start += quotes.TDX_SECURITY_QUOTES_MAX {
		length := count - start
		if length >= quotes.TDX_SECURITY_QUOTES_MAX {
			length = quotes.TDX_SECURITY_QUOTES_MAX
		}
		var subCodes []string
		for i := 0; i < length; i++ {
			securityCode := allCodes[start+i]
			subCodes = append(subCodes, securityCode)
			if barIndex != nil {
				bar.Add(1)
			}
		}
		if len(subCodes) == 0 {
			continue
		}
		codeCh <- subCodes
	}
	// channel 关闭后, 仍然可以读, 一直到读完全部数据
	close(codeCh)

	wg.Add(parallelCount)
	wg.Wait()
	// 如果有进度条
	if bar != nil {
		// 等待进度条结束
		bar.Wait()
	}

	__mutexTicks.Lock()
	for _, v := range snapshots {
		__cacheTicks[v.SecurityCode] = v
	}
	__mutexTicks.Unlock()

	if barIndex != nil {
		*barIndex++
	}
}
