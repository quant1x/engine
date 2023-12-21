package rules

import (
	"gitee.com/quant1x/engine/factors"
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
	ErrRangeOfFundFlow             = exception.New(errorRuleBase+3, "非资金流出")
)

// RuleBase 基本面规则
type RuleBase struct{}

func (r RuleBase) Kind() Kind {
	return KRuleBase
}

func (r RuleBase) Name() string {
	return "基础规则"
}

func (r RuleBase) Exec(snapshot factors.QuoteSnapshot) error {
	// 基础过滤规则
	securityCode := snapshot.Code
	// 1. 开盘换手Z的逻辑
	if num.IsNaN(snapshot.OpenTurnZ) || !RuleParameters.TurnZ.Validate(snapshot.OpenTurnZ) {
		return ErrRangeOfOpeningTurnZ
	}
	// 2. 当日 - 开盘量比
	if num.IsNaN(snapshot.OpenQuantityRatio) || !RuleParameters.QuantityRatio.Validate(snapshot.OpenQuantityRatio) {
		return ErrRangeOfOpeningQuantityRatio
	}
	// 3. 当日 - 开盘涨幅
	if num.IsNaN(snapshot.OpeningChangeRate) || !RuleParameters.OpenRate.Validate(snapshot.OpeningChangeRate) {
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
	// 6. exchange 过滤
	exchange := factors.GetL5Exchange(securityCode)
	if exchange != nil {
		//// 6.1 资金流向
		//if exchange.FundFlow != 0 && (exchange.FundFlow/TenThousand) < RuleParameters.MaxReduceAmount {
		//	return ErrRangeOfFundFlow
		//}
	}

	_ = securityCode
	// 规则通过
	return nil
}
