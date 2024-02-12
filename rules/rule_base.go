package rules

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gox/exception"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/num"
)

func init() {
	err := RegisterFunc(KRuleBase, "基础规则", ruleBase)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
}

var (
	ErrRangeOfOpeningTurnZ         = exception.New(errorRuleBase+0, "非开盘换手范围")
	ErrRangeOfOpeningQuantityRatio = exception.New(errorRuleBase+1, "非开盘量比范围")
	ErrRangeOfOpeningChangeRate    = exception.New(errorRuleBase+2, "非开盘涨跌幅范围")
	ErrRangeOfFundFlow             = exception.New(errorRuleBase+3, "非资金流出范围")
	ErrHistoryNotExist             = exception.New(errorRuleBase+4, "没有找到history数据")
	ErrRiskOfGapDown               = exception.New(errorRuleBase+5, "开盘存在向下跳空缺口")
	ErrExchangeNotExist            = exception.New(errorRuleBase+4, "没有找到history数据")
)

// ruleBase 基础规则
func ruleBase(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error {
	// 基础过滤规则
	securityCode := snapshot.Code
	// 1. 开盘换手Z的逻辑
	if num.IsNaN(snapshot.OpenTurnZ) || !ruleParameter.OpenTurnZ.Validate(snapshot.OpenTurnZ) {
		return ErrRangeOfOpeningTurnZ
	}
	// 2. 当日 - 开盘量比
	if num.IsNaN(snapshot.OpenQuantityRatio) || !ruleParameter.OpenQuantityRatio.Validate(snapshot.OpenQuantityRatio) {
		return ErrRangeOfOpeningQuantityRatio
	}
	// 3. 当日 - 开盘涨幅
	if num.IsNaN(snapshot.OpeningChangeRate) || !ruleParameter.OpenChangeRate.Validate(snapshot.OpeningChangeRate) {
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
	misc := factors.GetL5Misc(securityCode)
	if misc == nil {
		//return ErrExchangeNotExist
	} else {
		//// 6.1 资金流向
		//if exchange.FundFlow != 0 && (exchange.FundFlow/TenThousand) < RuleParameters.MaxReduceAmount {
		//	return ErrRangeOfFundFlow
		//}
	}
	// 7. 历史数据
	history := factors.GetL5History(securityCode)
	if history == nil {
		return ErrHistoryNotExist
	} else {
		// 7.1 开盘存在跳空缺口
		if !ruleParameter.GapDown && history.LOW >= snapshot.Open {
			return ErrRiskOfGapDown
		}
	}
	// 规则通过
	return nil
}
