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
	dataManifest
}

func init() {
	_ = cache.Register(&TransactionRecord{dataManifest: dataManifest{kind: BaseTransaction}})
}

func (r *TransactionRecord) Clone(date string, code string) DataSet {
	manifest := dataManifest{Date: date, Code: code, kind: BaseTransaction}
	var dest = TransactionRecord{dataManifest: manifest}
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
	base.GetTickAll(r.Code)
}

func (r *TransactionRecord) Repair(date string) {
	//base.GetTickAll(r.Code)
	base.GetTickData(r.Code, date)
}

func (r *TransactionRecord) Increase(snapshot quotes.Snapshot) {
	//TODO implement me
	panic("implement me")
}
