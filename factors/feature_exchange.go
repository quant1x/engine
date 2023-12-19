package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/datasource/dfcf"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
)

const (
	cacheL5KeyExchange = "exchange"
)

// Exchange 上一个交易日的数据快照
type Exchange struct {
	cache.DataSummary     `dataframe:"-"`
	Date                  string  `name:"日期" dataframe:"日期"`           // 数据日期
	Code                  string  `name:"证券代码" dataframe:"证券代码"`       // 证券代码
	Shape                 uint64  `name:"K线形态" dataframe:"K线形态"`       // K线形态
	MV3                   float64 `name:"前3日分钟均量" dataframe:"前3日分钟均量"` // 前3日分钟均量
	MA3                   float64 `name:"3日均线" dataframe:"3日均线"`       // 3日均价
	MV5                   float64 `name:"前5日分钟均量" dataframe:"前5日分钟均量"` // 前5日每分钟均量, 量比(QuantityRelativeRatio)需要
	MA5                   float64 `name:"5日均线" dataframe:"5日均线"`       // 5日均价
	MA10                  float64 `name:"10日均线" dataframe:"10日均线"`     // 10日均价
	MA20                  float64 `name:"20日均线" dataframe:"20日均线"`     // 生命线(MA20)/20日线
	FundFlow              float64 `name:"资金流向" dataframe:"资金流向"`       // 资金流向, 暂时无用
	VolumeRatio           float64 `name:"成交量比" dataframe:"成交量比"`       // 成交量放大比例, 相邻的两个交易日进行比对
	TurnoverRate          float64 `name:"换手率" dataframe:"换手率"`         // 换手率
	AmplitudeRatio        float64 `name:"振幅" dataframe:"振幅"`           // 振幅
	BidOpen               float64 `name:"竞价开盘" dataframe:"竞价开盘"`       // 竞价开盘价
	BidClose              float64 `name:"竞价结束" dataframe:"竞价结束"`       // 竞价结束
	BidHigh               float64 `name:"竞价最高" dataframe:"竞价最高"`       // 竞价最高
	BidLow                float64 `name:"竞价最低" dataframe:"竞价最低"`       // 竞价最低
	BidMatched            float64 `name:"竞价匹配量" dataframe:"竞价匹配量"`     // 竞价匹配量
	BidUnmatched          float64 `name:"竞价未匹配" dataframe:"竞价未匹配"`     // 竞价未匹配量
	BidDirection          int     `name:"竞价方向" dataframe:"竞价方向"`       // 竞价方向
	OpenBiddingDirection  int     `name:"开盘竞价" dataframe:"开盘竞价"`       // 竞价方向, 交易当日集合竞价开盘时更新
	OpenVolumeDirection   int     `name:"开盘竞量" dataframe:"开盘竞量"`       // 委托量差, 交易当日集合竞价开盘时更新
	CloseBiddingDirection int     `name:"收盘竞价" dataframe:"收盘竞价"`       // 竞价方向, 交易当日集合竞价收盘时更新
	CloseVolumeDirection  int     `name:"收盘竞量" dataframe:"收盘竞量"`       // 委托量差, 交易当日集合竞价收盘时更新
	OpenVolume            int64   `name:"开盘量" dataframe:"开盘量"`         // 开盘量
	OpenTurnZ             float64 `name:"开盘换手z" dataframe:"开盘换手z"`     // 开盘换手z
	CloseVolume           int64   `name:"收盘量" dataframe:"收盘量"`         // TODO:快照数据实际上有好几条, 应该用当日成交记录修订
	CloseTurnZ            float64 `name:"收盘换手z" dataframe:"收盘换手z"`     // 收盘换手z
	LastSentiment         float64 `name:"昨日情绪" dataframe:"昨日情绪"`       // 昨日情绪
	LastConsistent        int     `name:"昨日情绪一致" dataframe:"昨日情绪一致"`   // 昨日情绪一致
	OpenSentiment         float64 `name:"开盘情绪值" dataframe:"开盘情绪值"`     // 开盘情绪值, 个股没有
	OpenConsistent        int     `name:"开盘情绪一致" dataframe:"开盘情绪一致"`   // 开盘情绪一致, 个股没有
	CloseSentiment        float64 `name:"收盘情绪值" dataframe:"收盘情绪值"`     // 收盘情绪值
	CloseConsistent       int     `name:"收盘情绪一致" dataframe:"收盘情绪一致"`   // 收盘情绪一致
	AveragePrice          float64 `name:"均价线" dataframe:"均价线"`         // 均价线
	Volume                int64   `name:"成交量" dataframe:"成交量"`         // 成交量
	InnerVolume           int64   `name:"内盘" dataframe:"内盘"`           // 内盘
	OuterVolume           int64   `name:"外盘" dataframe:"外盘"`           // 外盘
	Change5               float64 `name:"5日涨幅" dataframe:"5日涨幅"`       // 5日涨幅
	Change10              float64 `name:"10日涨幅" dataframe:"10日涨幅"`     // 10日涨幅
	InitialPrice          float64 `name:"启动价格" dataframe:"启动价格"`       // 短线底部(Short-Term Bottom),股价最近一次上穿5日均线
	ShortIntensity        float64 `name:"短线强度" dataframe:"短线强度"`       // 短线强度,Strength
	ShortIntensityDiff    float64 `name:"短线强度增幅" dataframe:"短线强度增幅"`   // 短线强度
	MediumIntensity       float64 `name:"中线强度" dataframe:"中线强度"`       // 中线强度
	MediumIntensityDiff   float64 `name:"中线强度增幅" dataframe:"中线强度增幅"`   // 中线强度
	Vix                   float64 `name:"波动率" dataframe:"波动率"`         // 波动率
	State                 uint64  `name:"样本状态" dataframe:"样本状态"`
}

func NewExchange(date, code string) *Exchange {
	summary := __mapFeatures[FeatureExchange]
	v := Exchange{
		DataSummary: summary,
		Date:        date,
		Code:        code,
	}
	return &v
}

func (this *Exchange) GetDate() string {
	return this.Date
}

func (this *Exchange) GetSecurityCode() string {
	return this.Code
}

func (this *Exchange) Factory(date string, code string) Feature {
	v := NewExchange(date, code)
	return v
}

func (this *Exchange) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (this *Exchange) FromHistory(history History) Feature {
	_ = history
	return this
}

func (this *Exchange) Update(code, cacheDate, featureDate string, complete bool) {
	// 1. K线相关
	exchangeKLineExtend(this, code, featureDate)
	// 2. 成交量
	exchangeTurnZ(this, code, cacheDate, featureDate)
	// 3. 情绪
	exchangeSentiment(this, code, cacheDate, featureDate)
	// 4. 资金流向
	exchangeFundFlow(this, code, cacheDate, featureDate)
}

func (this *Exchange) Repair(code, cacheDate, featureDate string, complete bool) {
	// 1. K线相关
	exchangeKLineExtend(this, code, featureDate)
	// 2. 成交量, 使用cacheDate作为特征的缓存日期
	exchangeTurnZ(this, code, cacheDate, cacheDate)
	// 3. 情绪
	exchangeSentiment(this, code, cacheDate, featureDate)
	// 4. 资金流向
	exchangeFundFlow(this, code, cacheDate, featureDate)
}

func (this *Exchange) Increase(snapshot QuoteSnapshot) Feature {
	_ = snapshot
	return this
}

// ValidateSample 验证样本数据
func (this *Exchange) ValidateSample() error {
	if this.State > 0 {
		return nil
	}
	return ErrInvalidFeatureSample
}

// ExchangeKLineExtend 更新Exchange K线相关数据
func exchangeKLineExtend(info *Exchange, securityCode string, featureDate string) {
	cover := NewExchangeKLine(securityCode, featureDate)
	if cover == nil {
		logger.Errorf("code[%s, %s] kline not found", securityCode, featureDate)
		return
	}
	info.Date = cover.Date
	info.Shape = cover.Shape
	info.MV3 = cover.MV3
	info.MA3 = cover.MA3
	info.MV5 = cover.MV5
	info.MA5 = cover.MA5
	info.MA10 = cover.MA10
	info.MA20 = cover.MA20
	info.VolumeRatio = cover.VolumeRatio
	info.TurnoverRate = cover.TurnoverRate
	info.AmplitudeRatio = cover.AmplitudeRatio
	info.AveragePrice = cover.AveragePrice
	info.Change5 = cover.Change5
	info.Change10 = cover.Change10
	info.InitialPrice = cover.InitialPrice

	// 强弱指标
	info.ShortIntensity = cover.ShortIntensity
	info.ShortIntensityDiff = cover.ShortIntensityDiff
	info.MediumIntensity = cover.MediumIntensity
	info.MediumIntensityDiff = cover.MediumIntensityDiff

	// 波动率
	info.Vix = cover.Vix

	// 情绪, 指数和板块的情绪从K线上的up和down获取, 这里不处理个股的情绪
	if proto.AssertIndexBySecurityCode(info.Code) {
		info.LastSentiment = cover.Sentiment
		info.LastConsistent = cover.Consistent
	}
	info.State |= cover.Kind()
}

// 更新 - exchange - 历史成交数据相关, capture,collect
func exchangeTurnZ(info *Exchange, securityCode string, cacheDate, featureDate string) {
	list := base.Transaction(securityCode, featureDate)
	if len(list) > 0 {
		summary := CountInflow(list, securityCode, featureDate)
		// 修正f10的缓存, 应该是缓存日期为准
		f10 := GetL5F10(securityCode, cacheDate)
		if f10 != nil {
			summary.OpenTurnZ = f10.TurnZ(summary.OpenVolume)
			summary.CloseTurnZ = f10.TurnZ(summary.CloseVolume)
		}
		cover := summary
		info.OpenVolume = cover.OpenVolume
		info.OpenTurnZ = cover.OpenTurnZ
		info.CloseVolume = cover.CloseVolume
		info.CloseTurnZ = cover.CloseTurnZ
		info.Volume = cover.InnerVolume + cover.OuterVolume
		info.InnerVolume = cover.InnerVolume
		info.OuterVolume = cover.OuterVolume
	}
}

// 更新 - exchange - 情绪
func exchangeSentiment(info *Exchange, securityCode string, cacheDate, featureDate string) {
	if proto.AssertIndexBySecurityCode(securityCode) {
		// 跳过指数和板块, 只处理个股的情绪值
		return
	}
	list := base.Transaction(securityCode, featureDate)
	if len(list) > 0 {
		cover := CountInflow(list, securityCode, featureDate)
		info.LastSentiment, info.LastConsistent = market.SecuritySentiment(cover.OuterVolume, cover.InnerVolume)
	}
}

// 更新 - exchange - 资金流向
func exchangeFundFlow(info *Exchange, securityCode string, cacheDate, featureDate string) {
	if !proto.AssertStockBySecurityCode(securityCode) {
		return
	}
	beginDate := proto.MARKET_CH_FIRST_LISTTIME
	filename := cache.FundFlowFilename(securityCode)
	cacheList := []dfcf.FundFlow{}
	err := api.CsvToSlices(filename, &cacheList)
	cacheLength := len(cacheList)
	if err == nil && cacheLength > 0 {
		beginDate = cacheList[cacheLength-1].Date
		cacheList = cacheList[0 : cacheLength-1]
	}
	newList := dfcf.IndividualStocksFundFlow(securityCode, beginDate)
	if len(newList) == 0 {
		return
	}
	list := append(cacheList, newList...)
	_ = api.SlicesToCsv(filename, list)
	last := list[len(list)-1]
	cover := last.Medium
	info.FundFlow = cover
}
