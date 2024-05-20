package config

import (
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gox/api"
	"slices"
	"strings"
)

// StrategyParameter 策略参数
type StrategyParameter struct {
	Id                          uint64         `name:"策略编码" yaml:"id" default:"1"`                                     // 策略ID, 默认是1
	Auto                        bool           `name:"是否自动执行" yaml:"auto" default:"false"`                             // 是否自动执行
	Name                        string         `name:"策略名称" yaml:"name"`                                               // 策略名称
	Flag                        string         `name:"订单标识" yaml:"flag"`                                               // 订单标识,分早盘,尾盘和盘中
	Session                     TradingSession `name:"时间范围" yaml:"time" default:"09:30:00~11:30:00,13:00:00~14:56:30"` // 可操作的交易时段
	Weight                      float64        `name:"持仓占比" yaml:"weight" default:"0"`                                 // 策略权重, 默认0, 由系统自动分配
	Total                       int            `name:"订单数上限" yaml:"total" default:"3"`                                 // 订单总数, 默认是3
	PriceCageRatio              float64        `name:"价格笼子比例" yaml:"price_cage_ratio" default:"0.00"`                  // 价格笼子比例, 默认0%
	MinimumPriceFluctuationUnit float64        `name:"价格变动最小单位" yaml:"minimum_price_fluctuation_unit" default:"0.05"`  // 价格最小变动单位, 默认0.05
	FeeMax                      float64        `name:"最大费用" yaml:"fee_max" default:"20000.00"`                         // 可投入资金-最大
	FeeMin                      float64        `name:"最小费用" yaml:"fee_min" default:"10000.00"`                         // 可投入资金-最小
	Sectors                     []string       `name:"板块" yaml:"sectors" default:""`                                   // 板块, 策略适用的板块列表, 默认板块为空, 即全部个股
	IgnoreMarginTrading         bool           `name:"剔除两融" yaml:"ignore_margin_trading" default:"true"`               // 剔除两融标的, 默认是剔除
	HoldingPeriod               int            `name:"持仓周期" yaml:"holding_period" default:"1"`                         // 持仓周期, 默认为1天, 即T+1日触发117号策略
	SellStrategy                uint64         `name:"卖出策略" yaml:"sell_strategy" default:"117"`                        // 卖出策略, 默认117
	FixedYield                  float64        `name:"固定收益率" yaml:"fixed_yield" default:"0"`                           // 固定收益率, 只能和卖出策略绑定
	TakeProfitRatio             float64        `name:"止盈比例" yaml:"take_profit_ratio" default:"15.00"`                  // 止盈比例, 默认15%
	StopLossRatio               float64        `name:"止损比例" yaml:"stop_loss_ratio" default:"-2.00"`                    // 止损比例, 默认-2%
	LowOpeningAmplitude         float64        `name:"低开幅度" yaml:"low_opening_amplitude" default:"0.618"`              // 阳线, 低开幅度
	HighOpeningAmplitude        float64        `name:"高开幅度" yaml:"high_opening_amplitude" default:"0.382"`             // 阴线, 高开幅度
	Rules                       RuleParameter  `name:"规则参数" yaml:"rules"`                                              // 过滤规则
	excludeCodes                []string       `name:"过滤列表"`                                                           //  需要排除的个股
}

func (this *StrategyParameter) QmtStrategyName() string {
	return QmtStrategyNameFromId(this.Id)
}

// Enable 策略是否有效
func (this *StrategyParameter) Enable() bool {
	return this.Auto && this.Id >= 0
}

// BuyEnable 获取可买入状态
func (this *StrategyParameter) BuyEnable() bool {
	return this.Enable() && this.Total > 0
}

// SellEnable 获取可卖出状态
func (this *StrategyParameter) SellEnable() bool {
	return this.Enable()
}

// IsCookieCutterForSell 是否一刀切卖出
func (this *StrategyParameter) IsCookieCutterForSell() bool {
	return this.SellEnable() && this.Total == 0
}

// NumberOfTargets 获得可买入标的总数
func (this *StrategyParameter) NumberOfTargets() int {
	if !this.BuyEnable() {
		return 0
	}
	return this.Total
}

func (this *StrategyParameter) initExclude() {
	if len(this.excludeCodes) > 0 {
		return
	}
	var excludeCodes []string
	for _, v := range this.Sectors {
		sectorCode := strings.TrimSpace(v)
		if strings.HasPrefix(sectorCode, sectorIgnorePrefix) {
			sectorCode = strings.TrimSpace(sectorCode[sectorPrefixLength:])
			blockInfo := securities.GetBlockInfo(sectorCode)
			if blockInfo != nil {
				excludeCodes = append(excludeCodes, blockInfo.ConstituentStocks...)
			}
		}
	}
	excludeCodes = api.Unique(excludeCodes)
	this.excludeCodes = excludeCodes
}

func (this *StrategyParameter) Filter(codes []string) []string {
	this.initExclude()
	// 过滤需要忽略的板块成份股
	newCodeList := api.Filter(codes, func(s string) bool {
		return !slices.Contains(this.excludeCodes, s)
	})
	newCodeList = api.SliceUnique(newCodeList, func(a string, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		} else {
			return 0
		}
	})

	if this.IgnoreMarginTrading {
		// 过滤两融
		marginTradingList := securities.MarginTradingList()
		newCodeList = api.Filter(newCodeList, func(s string) bool {
			if slices.Contains(marginTradingList, s) {
				return false
			}
			return true
		})
	}
	return newCodeList
}

// StockList 取得可以交易的证券代码列表
func (this *StrategyParameter) StockList() []string {
	var codes []string
	for _, v := range this.Sectors {
		sectorCode := strings.TrimSpace(v)
		if !strings.HasPrefix(sectorCode, sectorIgnorePrefix) {
			blockInfo := securities.GetBlockInfo(sectorCode)
			if blockInfo != nil {
				codes = append(codes, blockInfo.ConstituentStocks...)

			}
		}
	}
	if len(codes) == 0 {
		codes = market.GetStockCodeList()
	}
	codes = this.Filter(codes)
	return codes
}
