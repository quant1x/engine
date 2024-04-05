package models

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/num"
	"sync"
	"time"
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
	snapshot.Code = securityCode
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
func SyncAllSnapshots() {
	start := time.Now()

	modName := "同步快照数据"
	allCodes := securities.AllCodeList()
	count := len(allCodes)

	progressManager := utils.NewProgressBarManager(modName, count)
	progressManager.Start()
	defer progressManager.Wait()

	currentDate := exchange.GetCurrentlyDay()
	tdxApi := gotdx.GetTdxApi()
	parallelCount := tdxApi.NumOfServers() / 2
	if parallelCount < 2 {
		parallelCount = 2
	}

	var wg sync.WaitGroup
	codeCh := make(chan []string, parallelCount)

	for i := 0; i < parallelCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for subCodes := range codeCh {
				for i := 0; i < quotes.DefaultRetryTimes; i++ {
					list, err := tdxApi.GetSnapshot(subCodes)
					if err != nil {
						logger.Errorf("ZS: 网络异常: %+v, 重试: %d", err, i+1)
						continue
					}

					for _, v := range list {
						v.Date = currentDate
						__mutexTicks.Lock()
						__cacheTicks[v.SecurityCode] = v
						__mutexTicks.Unlock()
					}

					break
				}
			}
		}()
	}

	for start := 0; start < count; start += quotes.TDX_SECURITY_QUOTES_MAX {
		length := count - start
		if length > quotes.TDX_SECURITY_QUOTES_MAX {
			length = quotes.TDX_SECURITY_QUOTES_MAX
		}
		subCodes := make([]string, 0, length)
		for i := 0; i < length; i++ {
			securityCode := allCodes[start+i]
			subCodes = append(subCodes, securityCode)
		}
		progressManager.Update(length)

		codeCh <- subCodes
	}

	close(codeCh)
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %v\n", elapsed)
}
