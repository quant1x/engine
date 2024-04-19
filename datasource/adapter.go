package datasource

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/securities"
	"sync"
	"time"
)

// Adapter 数据源适配器
type TickDataProvider interface {
	GetTickFromMemory(securityCode string) *quotes.Snapshot
	QuoteSnapshotFromProtocol(quote quotes.Snapshot) factors.QuoteSnapshot
	SyncAllSnapshots()
	Next() bool
	RegisterDaily(func())
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

func (a *RealTickDataAdapter) RegisterDaily(dailyFunc func()) {
}

////////////////////

type BacktestingTickDataAdapter struct {
	dates              []string
	currentDateIndex   int
	currentMinuteIndex int
	snapshots          map[string]quotes.Snapshot
	snapshotsMutex     sync.RWMutex
	securityCodes      []string
	dailyFunc          func()
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

		if a.dailyFunc != nil {
			a.dailyFunc()
		}
	}
	a.clearSnapshots()
	return true
}

func (a *BacktestingTickDataAdapter) RegisterDaily(dailyFunc func()) {
	a.dailyFunc = dailyFunc
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

type Stock struct {
	Date  time.Time
	Price float64
}

type Backtest struct {
	Stocks      map[string][]Stock
	BuyRecords  map[string][]Stock
	SellRecords map[string][]Stock
}

func NewBacktest() *Backtest {
	return &Backtest{
		Stocks:      make(map[string][]Stock),
		BuyRecords:  make(map[string][]Stock),
		SellRecords: make(map[string][]Stock),
	}
}

func (b *Backtest) Buy(stockName string, price float64) {
	stock := Stock{
		Date:  time.Now(),
		Price: price,
	}
	b.BuyRecords[stockName] = append(b.BuyRecords[stockName], stock)
	b.Stocks[stockName] = append(b.Stocks[stockName], stock)
	fmt.Printf("Bought %s at price %.2f\n", stockName, price)
}

func (b *Backtest) SellPreviousDayStocks(sellCondition func(stock Stock) bool) {
	previousDay := time.Now().AddDate(0, 0, -1)
	sellStocks := make(map[string][]Stock)

	for stockName, stocks := range b.Stocks {
		var soldStocks []Stock
		for _, stock := range stocks {
			if stock.Date.Day() == previousDay.Day() && sellCondition(stock) {
				sellStocks[stockName] = append(sellStocks[stockName], stock)
				soldStocks = append(soldStocks, stock)
			}
		}
		b.removeSoldStocks(stockName, soldStocks)
	}

	for stockName, stocks := range sellStocks {
		b.SellRecords[stockName] = append(b.SellRecords[stockName], stocks...)
	}

	printRecords("Sell Records:", b.SellRecords)
}

func (b *Backtest) removeSoldStocks(stockName string, stocks []Stock) {
	for _, stock := range stocks {
		for i, buyStock := range b.Stocks[stockName] {
			if buyStock.Date == stock.Date && buyStock.Price == stock.Price {
				b.Stocks[stockName] = append(b.Stocks[stockName][:i], b.Stocks[stockName][i+1:]...)
				break
			}
		}
	}
}

func printRecords(title string, records map[string][]Stock) {
	fmt.Println(title)
	for stockName, stocks := range records {
		for _, stock := range stocks {
			fmt.Printf("Stock Name: %s, Date: %s, Price: %.2f\n", stockName, stock.Date.Format("2006-01-02"), stock.Price)
		}
	}
}

func main() {
	backtest := NewBacktest()

	// 示例：每天进行买入操作
	backtest.Buy("AAPL", 150.0)
	backtest.Buy("GOOGL", 2500.0)
	backtest.Buy("MSFT", 300.0)

	// 示例：根据条件卖出前一天的股票
	sellCondition := func(stock Stock) bool {
		// 根据自定义条件判断是否卖出股票
		return stock.Price > 200.0
	}

	backtest.SellPreviousDayStocks(sellCondition)
}
