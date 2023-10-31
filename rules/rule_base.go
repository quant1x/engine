package rules

import (
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gox/exception"
	"gitee.com/quant1x/gox/num"
	"strings"
)

func init() {
	err := Register(RuleF10{})
	if err != nil {
		panic(err)
	}
}

var (
	ErrIgnoreStock                         = exception.New(errorRuleBase+0, "忽略的个股")
	ErrSubNewStock                         = exception.New(errorRuleBase+1, "次新股")
	ErrScienceAndTechnologyInnovationBoard = exception.New(errorRuleBase+2, "科创板")
	ErrPriceRange                          = exception.New(errorRuleBase+3, "股价控制")
	ErrRangeOfOpeningTurnZ                 = exception.New(errorRuleBase+4, "非开盘换手范围")
	ErrRangeOfOpeningQuantityRatio         = exception.New(errorRuleBase+5, "非开盘量比范围")
	ErrRangeOfOpeningChangeRate            = exception.New(errorRuleBase+6, "非开盘涨跌幅范围")
	ErrF10RangeOfCapital                   = exception.New(errorRuleBase+7, "非流通盘范围")
	ErrF10RangeOfSafetyCode                = exception.New(errorRuleBase+8, "非安全分范围")
	ErrF10RangeOfBasicEPS                  = exception.New(errorRuleBase+9, "非每股收益范围")
	ErrF10RangeOfBPS                       = exception.New(errorRuleBase+10, "非净增长范围")
)

type RuleF10 struct{}

func (RuleF10) Kind() Kind {
	return RuleBaseF10
}

func (RuleF10) Name() string {
	return "基础规则"
}

func (RuleF10) Exec(snapshot models.QuoteSnapshot) error {
	// 基础过滤规则, 检测F10基本面
	securityCode := snapshot.Code
	//if securityCode == "sz002211" {
	//	fmt.Println(securityCode)
	//}
	//currentDate := trading.GetCurrentlyDay()
	// 1. 去掉需要忽略的个股
	if market.IsNeedIgnore(securityCode) {
		return ErrIgnoreStock
	}
	// 2. 去掉科创板, 已知有688和689开头的9号公司
	if strings.HasPrefix(securityCode, "sh68") {
		return ErrScienceAndTechnologyInnovationBoard
	}
	// 3. 去掉次新股
	if market.IsSubNewStock(securityCode) {
		return ErrSubNewStock
	}
	// 4. 股价控制
	if num.IsNaN(snapshot.LastClose) || snapshot.LastClose < RuleParameters.PriceMin || snapshot.LastClose > RuleParameters.PriceMax {
		return ErrPriceRange
	}
	// 5. 开盘换手Z的逻辑
	if num.IsNaN(snapshot.OpenTurnZ) || snapshot.OpenTurnZ < RuleParameters.TurnZMin || snapshot.OpenTurnZ >= RuleParameters.TurnZMax {
		return ErrRangeOfOpeningTurnZ
	}
	// 6. 当日 - 开盘量比
	if num.IsNaN(snapshot.QuantityRatio) || snapshot.QuantityRatio < RuleParameters.QuantityRatioMin || snapshot.QuantityRatio > RuleParameters.QuantityRatioMax {
		return ErrRangeOfOpeningQuantityRatio
	}
	// 7. 当日 - 开盘涨幅
	if num.IsNaN(snapshot.OpeningChangeRate) || snapshot.OpeningChangeRate < RuleParameters.OpenRateMin || snapshot.OpeningChangeRate > RuleParameters.OpenRateMax {
		return ErrRangeOfOpeningChangeRate
	}
	// 8. 委托量
	//if snapshot.AverageBiddingVolume > RuleParameters.BiddingVolumeMax || snapshot.AverageBiddingVolume < RuleParameters.BiddingVolumeMin {
	//	return false
	//}
	//// 9. 力度-测试
	//if snapshot.ChangePower < 0 {
	//	return false
	//}
	// 10. F10数据
	f10 := smart.GetL5F10(securityCode)
	if f10 != nil {
		// 10.1 流通股本控制
		if f10.Capital != 0 && (f10.Capital < RuleParameters.CapitalMin || f10.Capital > RuleParameters.CapitalMax) {
			return ErrF10RangeOfCapital
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
