package rules

import (
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/exception"
	"gitee.com/quant1x/gox/num"
)

func init() {
	err := Register(RuleF10{})
	if err != nil {
		panic(err)
	}
}

var (
	ErrF10IgnoreStock                 = exception.New(errorRuleF10+0, "忽略的个股")
	ErrF10SubNewStock                 = exception.New(errorRuleF10+1, "次新股")
	ErrF10DisableBeijingStockExchange = exception.New(errorRuleF10+2, "禁止北交所")
	ErrF10DisableChiNextBoard         = exception.New(errorRuleF10+3, "禁止创业板")
	ErrF10DisableSciTechBoard         = exception.New(errorRuleF10+4, "禁止科创板")
	ErrF10PriceRange                  = exception.New(errorRuleF10+5, "股价控制")
	ErrF10RangeOfCapital              = exception.New(errorRuleF10+6, "非流通盘范围")
	ErrF10RangeOfSafetyCode           = exception.New(errorRuleF10+7, "非安全分范围")
	ErrF10RangeOfBasicEPS             = exception.New(errorRuleF10+8, "非每股收益范围")
	ErrF10RangeOfBPS                  = exception.New(errorRuleF10+9, "非净增长范围")
	ErrF10RangeOfMarketCap            = exception.New(errorRuleF10+10, "非市值范围")
)

// RuleF10 基本面规则
type RuleF10 struct{}

func (r RuleF10) Kind() Kind {
	return KRuleF10
}

func (r RuleF10) Name() string {
	return "基本面"
}

func (r RuleF10) Exec(snapshot factors.QuoteSnapshot) error {
	// 基础过滤规则, 检测F10基本面
	securityCode := snapshot.Code
	// 1. 去掉需要忽略的个股
	if market.IsNeedIgnore(securityCode) {
		return ErrF10IgnoreStock
	}
	// 2. 去掉科创板, 已知有688和689开头的9号公司
	//if strings.HasPrefix(securityCode, "sh68") {
	//	return ErrF10DisableSciTechBoard
	//}
	if api.StartsWith(securityCode, RuleParameters.IgnoreCodes) {
		return ErrF10IgnoreStock
	}
	// 3. 去掉次新股
	if market.IsSubNewStock(securityCode) {
		return ErrF10SubNewStock
	}
	// 4. 股价控制
	if num.IsNaN(snapshot.LastClose) || !RuleParameters.Price.Validate(snapshot.LastClose) {
		return ErrF10PriceRange
	}
	// 10. F10数据
	f10 := factors.GetL5F10(securityCode)
	if f10 != nil {
		// 10.1 流通股本控制
		if f10.Capital != 0 && !RuleParameters.Capital.Validate(f10.Capital/Billion) {
			return ErrF10RangeOfCapital
		}
		// 10.1.1 市值控制
		marketValue := f10.TotalCapital * snapshot.LastClose / Billion
		if !RuleParameters.MarketCap.Validate(marketValue) {
			return ErrF10RangeOfMarketCap
		}
		// 10.2 安全分太低
		if f10.SafetyScore != 0 && float64(f10.SafetyScore) < RuleParameters.SafetyScoreMin {
			return ErrF10RangeOfSafetyCode
		}
		// 10.3 季报不理想
		if f10.BasicEPS != 0 && f10.BasicEPS < 0 {
			return ErrF10RangeOfBasicEPS
		}
		// 10.4 净增长小于0
		if f10.BPS != 0 && f10.BPS < 0 {
			return ErrF10RangeOfBPS
		}
		//// 10.5 处理季报有增减持数据, 两个季度前十大流通股总数对比
		//reportDate, _ := api.ParseTime(f10.UpdateDate)
		//after := reportDate.AddDate(0, 2, 0).After(time.Now())
		//// 两月内减持的剔掉, 或者减持统计超过1%
		//if after && (f10.Top10Capital < 0 || f10.ReductionRatio < -1.00) {
		//	return false
		//}
		//// 10.6. 处理上市公司公告
		//if f10.Reduce > 0 || f10.Increase > 0 || f10.Risk > 0 {
		//	return false
		//}
	}
	// 规则通过
	return nil
}
