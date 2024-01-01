package trader

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/concurrent"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"sync"
)

const (
	qmtPositionsPath     = "qmt"           // 持仓缓存路径
	qmtPositionsFilename = "positions.csv" // 持仓数据文件名
)

// Position 持仓
type Position struct {
	AccountType     int     `name:"账户类型" dataframe:"account_type"`     // 账户类型
	AccountId       string  `name:"资金账户" dataframe:"account_id"`       // 资金账号
	StrategyCode    int     `name:"策略编码" dataframe:"strategy_code"`    // 策略编码
	OrderFlag       string  `name:"订单标识" dataframe:"order_flag"`       // 订单标识
	SecurityCode    string  `name:"证券代码" dataframe:"stock_code"`       // 证券代码, 例如"sh600000"
	Volume          int     `name:"持仓数量" dataframe:"volume"`           // 持仓数量,股票以'股'为单位, 债券以'张'为单位
	CanUseVolume    int     `name:"可卖数量" dataframe:"can_use_volume"`   // 可用数量, 股票以'股'为单位, 债券以'张'为单位
	OpenPrice       float64 `name:"开仓价" dataframe:"open_price"`        // 开仓价
	MarketValue     float64 `name:"市值" dataframe:"market_value"`       // 市值
	FrozenVolume    int     `name:"冻结数量" dataframe:"frozen_volume"`    // 冻结数量
	OnRoadVolume    int     `name:"在途股份" dataframe:"on_road_volume"`   // 在途股份
	YesterdayVolume int     `name:"昨夜拥股" dataframe:"yesterday_volume"` // 昨夜拥股
	AvgPrice        float64 `name:"成本价" dataframe:"avg_price"`         // 成本价
	CreateTime      string  `name:"创建时间" dataframe:"create_time"`      // 创建时间
	LastOrderId     string  `name:"前订单ID" dataframe:"last_order_id"`   // 前订单ID
	BuyTime         string  `name:"买入时间" dataframe:"buy_time"`         // 买入时间
	BuyPrice        float64 `name:"买入价格" dataframe:"buy_price"`        // 买入价格
	BuyVolume       int     `name:"买入数量" dataframe:"buy_volume"`       // 买入数量
	SellTime        string  `name:"卖出时间" dataframe:"sell_time"`        // 卖出时间
	SellPrice       float64 `name:"卖出价格" dataframe:"sell_price"`       // 卖出价格
	SellVolume      int     `name:"卖出数量" dataframe:"sell_volume"`      // 卖出数量
	CancelTime      string  `name:"撤单时间" dataframe:"cancel_time"`      // 撤单时间
	UpdateTime      string  `name:"更新时间" dataframe:"update_time"`      // 更新时间
}

// Key 用证券代码作为关键字
func (p *Position) Key() string {
	return fmt.Sprintf("%s", p.SecurityCode)
}

func (p *Position) Sync(other PositionDetail) bool {
	err := api.Copy(p, &other)
	if err != nil {
		return false
	}
	if len(p.CreateTime) == 0 && p.YesterdayVolume > 0 {
		// 如果创建时间等于空且昨夜拥股大于0, 则持股日期往前推一天
		today := trading.Today()
		dates := trading.LastNDate(today, 1)
		frontDate := dates[0] + " 00:00:00"
		p.CreateTime = frontDate
	}
	return true
}

// MergeFromOrder 订单合并到持仓
func (p *Position) MergeFromOrder(order OrderDetail) bool {
	// 1. 或者成交量为0, 直接返回
	if order.TradedVolume == 0 {
		return false
	}
	// 2. 获取快照
	snapshot := models.GetStrategySnapshot(p.SecurityCode)
	plus := order.OrderType == STOCK_BUY
	// 3. 缓存持仓和订单成本
	// 3.1 计算当前持仓的买入成本
	openValue := p.OpenPrice * float64(p.Volume)
	// 3.2 计算订单的买入成本
	orderValue := order.TradedPrice * float64(order.TradedVolume)
	// 4. 更新持仓量, TODO: 持仓量边界保护未实现
	if plus {
		// 增加持仓股份
		p.Volume += order.TradedVolume
		// 增加在途股份
		p.OnRoadVolume += order.TradedVolume
		// 更新开仓价
		p.OpenPrice = (openValue + orderValue) / float64(p.Volume)
		// 更新买入信息
		p.BuyTime = order.OrderTime
		p.BuyPrice = order.TradedPrice
		p.BuyVolume = order.TradedVolume
	} else {
		// 减少持仓
		p.Volume -= order.TradedVolume
		if p.Volume < 0 {
			p.Volume = 0
		}
		// 减少可用
		p.CanUseVolume -= order.TradedVolume
		if p.CanUseVolume < 0 {
			p.CanUseVolume = 0
		}
		// 更新卖出信息
		p.SellTime = order.OrderTime
		p.SellPrice = order.TradedPrice
		p.SellVolume = order.TradedVolume
		if p.Volume > 0 {
			p.OpenPrice = (openValue - orderValue) / float64(p.Volume)
		}
	}
	// 5. 更新市值
	p.MarketValue = snapshot.Price * float64(p.Volume)
	// 6. 修改 更新时间
	p.UpdateTime = order.OrderTime
	return true
}

var (
	accountType = SECURITY_ACCOUNT // 账户类型
	accountId   = "888xxxxxxx"     // 账户ID
)

var (
	periodicOnce coroutine.PeriodicOnce
	rwMutex      sync.RWMutex
	mapPositions = concurrent.NewTreeMap[string, *Position]()
)

// 持仓缓存路径
func getPositionsPath() string {
	path := fmt.Sprintf("%s/%s", cache.GetRootPath(), qmtPositionsPath)
	return path
}

// 持仓缓存文件名
func positionsFilename() string {
	filename := fmt.Sprintf("%s/%s-%s", getPositionsPath(), traderParameter.AccountId, qmtPositionsFilename)
	return filename
}

// 加载本地的持仓数据
func lazyLoadLocalPositions() {
	filename := positionsFilename()
	var list []Position
	err := api.CsvToSlices(filename, &list)
	if err != nil || len(list) == 0 {
		logger.Errorf("%s 没有有效数据, error=%+v", filename, err)
		return
	}
	for _, v := range list {
		code := v.SecurityCode
		mapPositions.Put(code, &v)
	}
}

// SyncPositions 同步持仓
func SyncPositions() {
	periodicOnce.Do(lazyLoadLocalPositions)
	list, err := QueryHolding()
	if err != nil {
		return
	}
	for _, v := range list {
		securityCode := proto.CorrectSecurityCode(v.StockCode)
		position, found := mapPositions.Get(securityCode)
		if !found {
			position = &Position{
				AccountType:  v.AccountType,
				AccountId:    v.AccountId,
				SecurityCode: securityCode,
			}
		}
		ok := position.Sync(v)
		if ok {
			mapPositions.Put(securityCode, position)
		}
	}
}

// UpdatePositions 更新持仓
func UpdatePositions() {
	periodicOnce.Do(lazyLoadLocalPositions)
	list, err := QueryOrders()
	if err != nil {
		return
	}
	for _, v := range list {
		securityCode := proto.CorrectSecurityCode(v.StockCode)
		position, found := mapPositions.Get(securityCode)
		if !found {
			position = &Position{
				AccountType:  v.AccountType,
				AccountId:    v.AccountId,
				SecurityCode: securityCode,
			}
		}
		ok := position.MergeFromOrder(v)
		if ok {
			mapPositions.Put(securityCode, position)
		}
	}
}

// CacheSync 缓存同步
func CacheSync() {
	methodName := "CacheSync"
	periodicOnce.Do(lazyLoadLocalPositions)
	length := mapPositions.Size()
	list := make([]Position, 0, length)
	mapPositions.Each(func(key string, value *Position) {
		list = append(list, *value)
	})
	cacheFilename := positionsFilename()
	err := api.SlicesToCsv(cacheFilename, list)
	if err != nil {
		logger.Errorf("services.trader:%s, error:%+v", methodName, err)
	}
}
