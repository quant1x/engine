package rules

import (
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gox/exception"
	"gitee.com/quant1x/gox/num"
)

func init() {
	err := Register(RuleBase{})
	if err != nil {
		panic(err)
	}
}

var (
	ErrRangeOfOpeningTurnZ         = exception.New(errorRuleBase+0, "非开盘换手范围")
	ErrRangeOfOpeningQuantityRatio = exception.New(errorRuleBase+1, "非开盘量比范围")
	ErrRangeOfOpeningChangeRate    = exception.New(errorRuleBase+2, "非开盘涨跌幅范围")
)

// RuleBase 基本面规则
type RuleBase struct{}

func (r RuleBase) Kind() Kind {
	return KRuleBase
}

func (r RuleBase) Name() string {
	return "基础规则"
}

func (r RuleBase) Exec(snapshot models.QuoteSnapshot) error {
	// 基础过滤规则
	securityCode := snapshot.Code
	// 1. 开盘换手Z的逻辑
	if num.IsNaN(snapshot.OpenTurnZ) || snapshot.OpenTurnZ < RuleParameters.TurnZMin || snapshot.OpenTurnZ >= RuleParameters.TurnZMax {
		return ErrRangeOfOpeningTurnZ
	}
	// 2. 当日 - 开盘量比
	if num.IsNaN(snapshot.OpenQuantityRatio) || snapshot.OpenQuantityRatio < RuleParameters.QuantityRatioMin || snapshot.OpenQuantityRatio > RuleParameters.QuantityRatioMax {
		return ErrRangeOfOpeningQuantityRatio
	}
	// 3. 当日 - 开盘涨幅
	if num.IsNaN(snapshot.OpeningChangeRate) || snapshot.OpeningChangeRate < RuleParameters.OpenRateMin || snapshot.OpeningChangeRate > RuleParameters.OpenRateMax {
		return ErrRangeOfOpeningChangeRate
	}
	//// 4. 委托量
	//if snapshot.AverageBiddingVolume > RuleParameters.BiddingVolumeMax || snapshot.AverageBiddingVolume < RuleParameters.BiddingVolumeMin {
	//	return false
	//}
	//// 5. 力度-测试
	//if snapshot.ChangePower < 0 {
	//	return false
	//}

	_ = securityCode
	// 规则通过
	return nil
}
