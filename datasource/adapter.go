package datasource

import (
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/securities"
	"sync"
)

// Adapter 数据源适配器
type TickDataProvider interface {
	GetTickFromMemory(securityCode string) *quotes.Snapshot
	QuoteSnapshotFromProtocol(quote quotes.Snapshot) factors.QuoteSnapshot
	SyncAllSnapshots()
	Next() bool
}

type RealTickDataAdapter struct {
}

func (a *RealTickDataAdapter) GetTickFromMemory(securityCode string) *quotes.Snapshot {
	return models.GetTickFromMemory(securityCode)
}

func (a *RealTickDataAdapter) QuoteSnapshotFromProtocol(quote quotes.Snapshot) factors.QuoteSnapshot {
	return models.QuoteSnapshotFromProtocol(quote)
}

func (a *RealTickDataAdapter) SyncAllSnapshots() {
	models.SyncAllSnapshots()
}

func (a *RealTickDataAdapter) Next() bool {
	//models.SyncAllSnapshots()
	return true
}

////////////////////

type BacktestingTickDataAdapter struct {
	dates              []string
	currentDateIndex   int
	currentMinuteIndex int
	snapshots          map[string]quotes.Snapshot
	snapshotsMutex     sync.RWMutex
	securityCodes      []string
}

func NewBacktestingTickDataAdapter(startDate, endDate string) *BacktestingTickDataAdapter {
	dates := exchange.TradeRange(startDate, endDate)
	securityCodes := securities.AllCodeList()
	//securityCodes = securityCodes[:10]
	return &BacktestingTickDataAdapter{
		dates:              dates,
		currentDateIndex:   0,
		currentMinuteIndex: 0,
		snapshots:          make(map[string]quotes.Snapshot),
		securityCodes:      securityCodes,
	}
}

func (a *BacktestingTickDataAdapter) GetTickFromMemory(securityCode string) *quotes.Snapshot {
	a.snapshotsMutex.RLock()
	defer a.snapshotsMutex.RUnlock()

	snapshot, ok := a.snapshots[securityCode]
	if !ok {
		return nil
	}

	return &snapshot
}

func (a *BacktestingTickDataAdapter) QuoteSnapshotFromProtocol(quote quotes.Snapshot) factors.QuoteSnapshot {
	return models.QuoteSnapshotFromProtocol(quote)
}

func (a *BacktestingTickDataAdapter) clearSnapshots() {
	a.snapshotsMutex.Lock()
	a.snapshots = make(map[string]quotes.Snapshot)
	a.snapshotsMutex.Unlock()
}

func (a *BacktestingTickDataAdapter) Next() bool {
	if len(a.dates) == 0 {
		return false
	}

	a.currentMinuteIndex++
	if a.currentMinuteIndex == 240 {
		a.currentDateIndex++
		a.currentMinuteIndex = 0

		if a.currentDateIndex >= len(a.dates) {
			return false
		}
	}
	a.clearSnapshots()
	return true
}

func (a *BacktestingTickDataAdapter) SyncAllSnapshots() {
	modName := "同步快照数据"

	count := len(a.securityCodes)

	progressManager := utils.NewProgressBarManager(modName, count)
	progressManager.Start()
	defer progressManager.Wait()

	base.ClearCachedKLines()
	for _, securityCode := range a.securityCodes {
		progressManager.Update(1)
		currentDate := a.dates[a.currentDateIndex]
		currentKLine := base.CheckoutKLine(securityCode, currentDate)

		if currentKLine == nil {
			continue
		}

		minKLines := base.LoadMinutes(securityCode, currentDate)
		if len(minKLines) == 0 {
			continue
		}

		currentMinKLine := minKLines[a.currentMinuteIndex]

		a.snapshotsMutex.Lock()
		a.snapshots[securityCode] = *base.CombineKLinesToSnapshot(securityCode, a.currentMinuteIndex, currentKLine, &currentMinKLine)
		a.snapshotsMutex.Unlock()
	}
}
