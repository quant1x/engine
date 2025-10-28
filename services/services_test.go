package services

import (
	"testing"

	"gitee.com/quant1x/data/level1"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
)

func TestGlobalReset(t *testing.T) {
	_ = cleanExpiredStateFiles()
	level1.ReOpen()
	date := cache.DefaultCanUpdateDate()
	factors.SwitchDate(date)
}
