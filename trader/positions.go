package trader

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/concurrent"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"sync"
)

var (
	accountType = SECURITY_ACCOUNT // 账户类型
	accountId   = "888xxxxxxx"     // 账户ID
)

var (
	periodicOnce coroutine.PeriodicOnce
	rwMutex      sync.RWMutex
	mapPositions = concurrent.NewTreeMap[string, Position]()
)

// 持仓缓存路径
func getPositionsPath() string {
	path := fmt.Sprintf("%s/%s", cache.GetRootPath(), qmtPositionsPath)
	return path
}

// 持仓缓存文件名
func positionsFilename() string {
	filename := fmt.Sprintf("%s/%s", getPositionsPath(), qmtPositionsFilename)
	return filename
}

// 加载持仓数据
func lazyLoadPositions() {
	filename := positionsFilename()
	var list []Position
	err := api.CsvToSlices(filename, &list)
	if err != nil || len(list) == 0 {
		logger.Errorf("%s 没有有效数据, error=%+v", filename, err)
		return
	}
	for _, v := range list {
		code := v.StockCode
		mapPositions.Put(code, v)
	}
}

//func PositionsAdd(direction Direction, model models.Strategy, securityCode string, price float64, volume int) {
//	periodicOnce.Do(lazyLoadPositions)
//	position, ok := mapPositions.Get(securityCode)
//	if direction == BUY {
//
//	}
//}

// Position 持仓
type Position struct {
	AccountType     int     `name:"账户类型" dataframe:"account_type"`     // 账户类型
	AccountId       string  `name:"资金账户" dataframe:"account_id"`       // 资金账号
	StrategyCode    int     `name:"策略编码" dataframe:"strategy_code"`    // 策略编码
	OrderFlag       string  `name:"订单标识" dataframe:"order_flag"`       // 订单标识
	StockCode       string  `name:"证券代码" dataframe:"stock_code"`       // 证券代码, 例如"600000.SH"
	Volume          int     `name:"持仓数量" dataframe:"volume"`           // 持仓数量,股票以'股'为单位, 债券以'张'为单位
	CanUseVolume    int     `name:"可卖数量" dataframe:"can_use_volume"`   // 可用数量, 股票以'股'为单位, 债券以'张'为单位
	OpenPrice       float64 `name:"开仓价" dataframe:"open_price"`        // 开仓价
	MarketValue     float64 `name:"市值" dataframe:"market_value"`       // 市值
	FrozenVolume    int     `name:"冻结数量" dataframe:"frozen_volume"`    // 冻结数量
	OnRoadVolume    int     `name:"在途股份" dataframe:"on_road_volume"`   // 在途股份
	YesterdayVolume int     `name:"昨夜拥股" dataframe:"yesterday_volume"` // 昨夜拥股
	AvgPrice        float64 `name:"成本价" dataframe:"avg_price"`         // 成本价
	CreateTime      string  `name:"创建时间" dataframe:"create_time"`      // 创建时间
	BuyTime         string  `name:"买入时间" dataframe:"buy_time"`         // 买入时间
	BuyPrice        float64 `name:"买入价格" dataframe:"buy_price"`        // 买入价格
	BuyVolume       int     `name:"买入数量" dataframe:"buy_volume"`       // 买入数量
	SellTime        string  `name:"卖出时间" dataframe:"sell_time"`        // 卖出时间
	SellPrice       float64 `name:"卖出价格" dataframe:"sell_price"`       // 卖出价格
	SellVolume      int     `name:"卖出数量" dataframe:"sell_volume"`      // 卖出数量
	CancelTime      string  `name:"撤单时间" dataframe:"cancel_time"`      // 撤单时间
	UpdateTime      string  `name:"更新时间" dataframe:"update_time"`      // 更新时间
}
