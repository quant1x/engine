package tracker

import (
	"fmt"
	"slices"
	"strings"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/data/level1/securities"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gox/api"
)

// CheckStrategy 检查当前交易日中个股在策略中的执行情况
func CheckStrategy(strategyCode uint64, securityCode, testDate string) {
	fmt.Printf("\n策略检测中...\n")
	// 1. 获取快照
	name := securities.GetStockName(securityCode)
	fmt.Printf("\t=> 证券代码: %s, 证券名称: %s...\n", securityCode, name)
	fmt.Printf("\t=> 1. 获取tick[%s]...\n", securityCode)
	testDate = strings.TrimSpace(testDate)
	var snapshot *factors.QuoteSnapshot
	if len(testDate) == 0 {
		snapshot = models.GetTick(securityCode)
		testDate = cache.DefaultCanReadDate()
	} else {
		testDate = exchange.FixTradeDate(testDate)
		//cacheDate, featureDate := cache.CorrectDate(testDate)
		features := factors.CheckoutWideTableByDate(securityCode, testDate)
		rows := len(features)
		if rows > 0 {
			tick := models.FeatureToSnapshot(features[rows-1], securityCode)
			snapshot = &tick
			factors.SwitchDate(testDate)
		}
	}
	fmt.Printf("\t=> 1. 获取tick[%s], date=%s...\n", securityCode, testDate)
	if snapshot == nil {
		fmt.Printf("\t=> 1. 获取tick[%s]...failed\n", securityCode)
		return
	}
	fmt.Printf("\t=> 1. 获取tick[%s]...success\n", securityCode)

	// 2. 获取策略配置
	fmt.Printf("\t=> 2. 获取策略[%d]配置...\n", strategyCode)
	strategyParameter := config.GetStrategyParameterByCode(strategyCode)
	if strategyParameter == nil {
		fmt.Printf("\t=> 2. 获取策略[%d]配置...not found\n", strategyCode)
		return
	}
	fmt.Printf("\t=> 2. 获取策略[%d]配置...success\n", strategyCode)
	fmt.Printf("\t=> 2. 获取策略[%d]配置, 策略名称=%s\n", strategyCode, strategyParameter.Name)

	// 3. 检测板块及两融匹配
	fmt.Printf("\t=> 3. 检测策略[%d]板块配置...\n", strategyCode)
	fmt.Printf("\t=> 3. 检测策略[%d]板块配置...是否需要剔除两融...\n", strategyCode)
	if strategyParameter.IgnoreMarginTrading {
		fmt.Printf("\t=> 3. 检测策略[%d]板块配置...是否需要剔除两融, 需要\n", strategyCode)
		// 过滤两融
		marginTradingList := securities.MarginTradingList()
		if len(marginTradingList) == 0 {
			fmt.Printf("\t=> 3. 检测策略[%d]板块配置...是否需要剔除两融, 需要, 两融列表为空, 跳过检测\n", strategyCode)
		} else if slices.Contains(marginTradingList, securityCode) {
			fmt.Printf("\t=> 3. 检测策略[%d]板块配置...是否需要剔除两融, 需要, 检测失败: [%s]为两融标的,\n", strategyCode, securityCode)
			return
		}
	} else {
		fmt.Printf("\t=> 3. 检测策略[%d]板块配置...是否需要剔除两融, 不需要\n", strategyCode)
	}
	fmt.Printf("\t=> 3. 检测策略[%d]板块配置...是否需要剔除两融...success\n", strategyCode)

	// 4. 检测板块及两融匹配
	fmt.Printf("\t=> 4. 检测策略[%d]板块是否匹配...\n", strategyCode)
	stockList := strategyParameter.StockList()
	if !slices.Contains(stockList, securityCode) {
		fmt.Printf("\t=> 4. 检测策略[%d]板块是否匹配...失败, %s非策略配置的板块成分股\n", strategyCode, securityCode)
		return
	}
	fmt.Printf("\t=> 4. 检测策略[%d]板块是否匹配...success\n", strategyCode)

	// 5. 获取策略对象
	fmt.Printf("\t=> 5. 获取策略[%d]对象...\n", strategyCode)
	model, err := models.CheckoutStrategy(strategyCode)
	if err != nil {
		fmt.Printf("\t=> 5. 获取策略[%d]对象...失败: %+v\n", strategyCode, err)
		return
	}
	fmt.Printf("\t=> 5. 获取策略[%d]对象...success\n", strategyCode)

	// 6. 执行过滤规则
	fmt.Printf("\t=> 6. 执行策略[%d]过滤规则...\n", strategyCode)
	var rules config.RuleParameter
	_ = api.Copy(&rules, &strategyParameter.Rules)
	//rules = strategyParameter.Rules
	rules.Verbose = true
	v := model.Filter(rules, *snapshot)
	if v == nil {
		fmt.Printf("\t=> 6. 执行策略[%d]过滤规则...passed\n", strategyCode)
	} else {
		fmt.Printf("\t=> 6. 执行策略[%d]过滤规则...failed: %+v\n", strategyCode, v)
	}
}
