package services

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"testing"
)

func TestTaskSell_getEarlierDate(t *testing.T) {
	v := getEarlierDate(1)
	fmt.Println(v)
}

func TestTaskSell_LastSession(t *testing.T) {
	sellStrategyCode := models.ModelOneSizeFitsAllSells
	// 1. 获取117号策略(卖出)
	sellRule := config.GetTradeRule(sellStrategyCode)
	if sellRule == nil {
		return
	}
	timestamp := "14:55:00"
	v := sellRule.Session.IsTodayLastSession(timestamp)
	fmt.Println(v)
}
