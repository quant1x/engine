package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
)

const (
	cacheL5KeySecuritiesMarginTrading = "rzrq"
)

// SecuritiesMarginTrading 融资融券
type SecuritiesMarginTrading struct {
	cache.DataSummary `dataframe:"-"`
	Date              string  `name:"日期" dataframe:"日期"`     // 数据日期
	Code              string  `name:"证券代码" dataframe:"证券代码"` // 证券代码
	RZYE              float64 `name:"融资余额(元)" dataframe:"RZYE"`
	RQYL              float64 `name:"融券余量(股)" dataframe:"RQYL"`
	RZRQYE            float64 `name:"融资融券余额(元)" dataframe:"RZRQYE"`
	RQYE              float64 `name:"融券余额(元)" dataframe:"RQYE"`
	RQMCL             float64 `name:"融券卖出量(股)" dataframe:"RQMCL"`
	RZRQYECZ          float64 `name:"融资融券余额差值(元)" dataframe:"RZRQYECZ"`
	RZMRE             float64 `name:"融资买入额(元)" dataframe:"RZMRE"`
	SZ                float64 `name:"SZ" dataframe:"SZ"`
	RZYEZB            float64 `name:"融资余额占流通市值比(%)" dataframe:"RZYEZB"`
	RZMRE3D           float64 `name:"3日融资买入额(元)" dataframe:"RZMRE3D"`
	RZMRE5D           float64 `name:"5日融资买入额(元)" dataframe:"RZMRE5D"`
	RZMRE10D          float64 `name:"10日融资买入额(元)" dataframe:"RZMRE10D"`
	RZCHE             float64 `name:"融资偿还额(元)" dataframe:"RZCHE"`
	RZCHE3D           float64 `name:"3日融资偿还额(元)" dataframe:"RZCHE3D"`
	RZCHE5D           float64 `name:"5日融资偿还额(元)" dataframe:"RZCHE5D"`
	RZCHE10D          float64 `name:"10日融资偿还额(元)" dataframe:"RZCHE10D"`
	RZJME             float64 `name:"融资净买额(元)" dataframe:"RZJME"`
	RZJME3D           float64 `name:"3日融资净买额(元)" dataframe:"RZJME3D"`
	RZJME5D           float64 `name:"5日融资净买额(元)" dataframe:"RZJME5D"`
	RZJME10D          float64 `name:"10日融资净买额(元)" dataframe:"RZJME10D"`
	RQMCL3D           float64 `name:"3日融券卖出量(股)" dataframe:"RQMCL3D"`
	RQMCL5D           float64 `name:"5日融券卖出量(股)" dataframe:"RQMCL5D"`
	RQMCL10D          float64 `name:"10日融券卖出量(股)" dataframe:"RQMCL10D"`
	RQCHL             float64 `name:"融券偿还量(股)" dataframe:"RQCHL"`
	RQCHL3D           float64 `name:"3日融券偿还量(股)" dataframe:"RQCHL3D"`
	RQCHL5D           float64 `name:"5日融券偿还量(股)" dataframe:"RQCHL5D"`
	RQCHL10D          float64 `name:"10日融券偿还量(股)" dataframe:"RQCHL10D"`
	RQJMG             float64 `name:"融券净卖出(股)" dataframe:"RQJMG"`
	RQJMG3D           float64 `name:"3日融券净卖出(股)" dataframe:"RQJMG3D"`
	RQJMG5D           float64 `name:"5日融券净卖出(股)" dataframe:"RQJMG5D"`
	RQJMG10D          float64 `name:"10日融券净卖出(股)" dataframe:"RQJMG10D"`
	SPJ               float64 `name:"收盘价" dataframe:"SPJ"`
	ZDF               float64 `name:"涨跌幅" dataframe:"ZDF"`
	RChange3DCP       float64 `name:"3日未识别" dataframe:"RCHANGE3DCP"`
	RChange5DCP       float64 `name:"5日未识别" dataframe:"RCHANGE5DCP"`
	RChange10DCP      float64 `name:"10日未识别" dataframe:"RCHANGE10DCP"`
	KCB               int     `name:"科创板"  dataframe:"KCB"`
	TradeMarketCode   string  `name:"二级市场代码" dataframe:"TRADE_MARKET_CODE"`
	TradeMarket       string  `name:"二级市场" dataframe:"TRADE_MARKET"`
	FinBalanceGr      float64 `name:"FIN_BALANCE_GR" dataframe:"FIN_BALANCE_GR"`
	UpdateTime        string  `name:"更新时间" dataframe:"update_time"` // 更新时间
	State             uint64  `name:"样本状态" dataframe:"样本状态"`        // 样本状态
}

// NewSecuritiesMarginTrading 新建融资融券
func NewSecuritiesMarginTrading(date, code string) *SecuritiesMarginTrading {
	summary := __mapFeatures[FeatureSecuritiesMarginTrading]
	v := SecuritiesMarginTrading{
		DataSummary: summary,
		Date:        date,
		Code:        code,
	}
	return &v
}

func (this *SecuritiesMarginTrading) Factory(date string, code string) Feature {
	v := NewSecuritiesMarginTrading(date, code)
	return v
}

func (this *SecuritiesMarginTrading) GetDate() string {
	return this.Date
}

func (this *SecuritiesMarginTrading) GetSecurityCode() string {
	return this.Code
}

func (this *SecuritiesMarginTrading) Init(ctx context.Context, date string) error {
	MarginTradingTargetInit(this.GetDate())
	_ = ctx
	_ = date
	return nil
}

func (this *SecuritiesMarginTrading) Update(code, cacheDate, featureDate string, whole bool) {
	securityCode := exchange.CorrectSecurityCode(code)
	this.Date = exchange.FixTradeDate(cacheDate)
	this.Code = securityCode
	rzrq, ok := GetMarginTradingTarget(this.GetSecurityCode())
	if !ok {
		return
	}
	api.Copy(this, &rzrq)

	this.UpdateTime = GetTimestamp()
	this.State |= this.Kind()

}

func (this *SecuritiesMarginTrading) Repair(securityCode, cacheDate, featureDate string, whole bool) {
	this.Update(securityCode, cacheDate, featureDate, whole)
}

func (this *SecuritiesMarginTrading) FromHistory(history History) Feature {
	_ = history
	return this
}

func (this *SecuritiesMarginTrading) Increase(snapshot QuoteSnapshot) Feature {
	_ = snapshot
	return this
}

func (this *SecuritiesMarginTrading) ValidateSample() error {
	if this.State > 0 {
		return nil
	}
	return ErrInvalidFeatureSample
}
