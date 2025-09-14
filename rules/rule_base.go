package rules

import (
	"errors"
	"fmt"

	"github.com/quant1x/engine/config"
	"github.com/quant1x/engine/factors"
	"github.com/quant1x/num"
	"github.com/quant1x/x/exception"
	"github.com/quant1x/x/logger"
)

func init() {
	err := RegisterFunc(KRuleBase, "基础规则", ruleBase)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
}

var (
	ErrRangeOfOpeningTurnZ          = exception.New(errorRuleBase+0, "非开盘换手范围")
	ErrRangeOfOpeningQuantityRatio  = exception.New(errorRuleBase+1, "非开盘量比范围")
	ErrRangeOfOpeningChangeRate     = exception.New(errorRuleBase+2, "非开盘涨跌幅范围")
	ErrRangeOfFundFlow              = exception.New(errorRuleBase+3, "非资金流出范围")
	ErrHistoryNotExist              = exception.New(errorRuleBase+4, "没有找到history数据")
	ErrRiskOfGapDown                = exception.New(errorRuleBase+5, "开盘存在向下跳空缺口")
	ErrExchangeNotExist             = exception.New(errorRuleBase+6, "没有找到history数据")
	ErrRangeOfChangeRate            = exception.New(errorRuleBase+7, "非实时涨跌幅范围")
	ErrRangeOfFinancingBalanceRatio = exception.New(errorRuleBase+8, "融资余额占比过大")
)

// 判断是否冗详模式输出错误信息
func throwException(err error, ruleParameter config.RuleParameter, value float64) error {
	if !ruleParameter.Verbose {
		return err
	} else {
		return errors.New(fmt.Sprintf("%s, %f", err.Error(), value))
	}
}

// ruleBase 基础规则
func ruleBase(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error {
	// 基础过滤规则
	securityCode := snapshot.SecurityCode
	// 1. 开盘换手Z的逻辑
	if num.IsNaN(snapshot.OpenTurnZ) || !ruleParameter.OpenTurnZ.Validate(snapshot.OpenTurnZ) {
		return throwException(ErrRangeOfOpeningTurnZ, ruleParameter, snapshot.OpenTurnZ)
	}
	// 2. 当日 - 开盘量比
	if num.IsNaN(snapshot.OpenQuantityRatio) || !ruleParameter.OpenQuantityRatio.Validate(snapshot.OpenQuantityRatio) {
		return throwException(ErrRangeOfOpeningQuantityRatio, ruleParameter, snapshot.OpenQuantityRatio)
	}
	// 3. 当日 - 开盘涨幅
	if num.IsNaN(snapshot.OpeningChangeRate) || !ruleParameter.OpenChangeRate.Validate(snapshot.OpeningChangeRate) {
		return ErrRangeOfOpeningChangeRate
	}
	// 3.1 当日 - 涨幅
	if num.IsNaN(snapshot.ChangeRate) || !ruleParameter.ChangeRate.Validate(snapshot.ChangeRate) {
		return ErrRangeOfChangeRate
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
		// 6.2 检查融资余额占比
		if misc.RZYEZB > 0 && misc.RZYEZB >= ruleParameter.FinancingBalanceRatio {
			return throwException(ErrRangeOfFinancingBalanceRatio, ruleParameter, misc.RZYEZB)
		}
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
