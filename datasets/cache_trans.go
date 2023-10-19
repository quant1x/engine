package datasets

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/gotdx/quotes"
)

// TransactionRecord 成交记录
//
//	最短3秒内的合并统计数据, 与行情数据保持一致
//	不可以当作tick数据来使用
type TransactionRecord struct {
	cache.Scheme
}

func init() {
	scheme := cache.DataScheme("", "", mapDataSets[BaseTransaction])
	_ = cache.Register(&TransactionRecord{Scheme: scheme})
}

func (r *TransactionRecord) Clone(date string, code string) DataSet {
	scheme := cache.DataScheme(date, code, mapDataSets[BaseTransaction])
	var dest = TransactionRecord{Scheme: scheme}
	return &dest
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
	base.UpdateTickStartDate(date)
	base.GetTickAll(r.GetSecurityCode())
}

func (r *TransactionRecord) Repair(date string) {
	//base.GetTickAll(r.Code)
	base.GetTickData(r.GetSecurityCode(), date)
}

func (r *TransactionRecord) Increase(snapshot quotes.Snapshot) {
	//TODO implement me
	panic("implement me")
}
