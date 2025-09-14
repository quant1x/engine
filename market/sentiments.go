package market

import (
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/num"
)

const (
	kSentimentCriticalValue = float64(61.80) // 情绪临界值
)

type SentimentType = int

const (
	SentimentZero SentimentType = iota // 情绪一般
	SentimentHigh SentimentType = 1    // 情绪高涨
	SentimentLow  SentimentType = -1   // 情绪低迷
)

// IndexSentiment 情绪指数, 50%为平稳, 低于50%为情绪差, 高于50%为情绪好
func IndexSentiment(codes ...string) (sentiment float64, consistent int) {
	// 默认上证指数
	securityCode := "sh000001"
	if len(codes) > 0 {
		securityCode = exchange.CorrectSecurityCode(codes[0])
	}
	// 非指数不判断情绪
	if !exchange.AssertIndexBySecurityCode(securityCode) {
		return
	}
	if len(securityCode) != 8 {
		return
	}
	tdxApi := gotdx.GetTdxApi()
	hq, err := tdxApi.GetSnapshot([]string{securityCode})
	if err != nil && len(hq) != len(codes) {
		logger.Errorf("获取即时行情数据失败", err)
		return
	}

	return SnapshotSentiment(hq[0])
}

// SecuritySentiment 计算证券情绪
func SecuritySentiment[E ~int | ~int64 | ~float32 | ~float64](up, down E) (sentiment float64, consistent int) {
	sentiment = 100 * float64(up) / float64(up+down)
	if num.Float64IsNaN(sentiment) {
		sentiment = 0
		return
	}
	consistent = SentimentZero
	if sentiment >= kSentimentCriticalValue {
		consistent = SentimentHigh
	} else if sentiment < 100-kSentimentCriticalValue {
		consistent = SentimentLow
	}
	return
}

// SnapshotSentiment 情绪指数, 50%为平稳, 低于50%为情绪差, 高于50%为情绪好
func SnapshotSentiment(snapshot quotes.Snapshot) (sentiment float64, consistent int) {
	if exchange.AssertIndexByMarketAndCode(snapshot.Market, snapshot.Code) {
		// 指数和板块按照上涨和下跌家数计算
		return SecuritySentiment(snapshot.IndexUp, snapshot.IndexDown)
	}
	// 个股按照内外盘计算
	return SecuritySentiment(snapshot.BVol, snapshot.SVol)
}
