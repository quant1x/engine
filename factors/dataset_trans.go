package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/gotdx/quotes"
)

// TransactionRecord 成交记录
//
//	最短3秒内的合并统计数据, 与行情数据保持一致
//	不可以当作tick数据来使用
type TransactionRecord struct {
	cache.DataSummary
	Date string
	Code string
}

func init() {
	summary := mapDataSets[BaseTransaction]
	_ = cache.Register(&TransactionRecord{DataSummary: summary})
}

func (r *TransactionRecord) Clone(date string, code string) DataSet {
	summary := mapDataSets[BaseTransaction]
	var dest = TransactionRecord{DataSummary: summary, Date: date, Code: code}
	return &dest
}

func (r *TransactionRecord) GetDate() string {
	return r.Date
}

func (r *TransactionRecord) GetSecurityCode() string {
	return r.Code
}

func (r *TransactionRecord) Print(code string, date ...string) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (r *TransactionRecord) Checkout(securityCode, date string) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Check(cacheDate, featureDate string) error {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Filename(date, code string) string {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Update(date string) {
	base.UpdateBeginDateOfHistoricalTradingData(date)
	base.GetAllHistoricalTradingData(r.GetSecurityCode())
}

func (r *TransactionRecord) Repair(date string) {
	//base.GetAllHistoricalTradingData(r.code)
	base.GetHistoricalTradingDataByDate(r.GetSecurityCode(), date)
}

func (r *TransactionRecord) Increase(snapshot quotes.Snapshot) {
	_ = snapshot
}
