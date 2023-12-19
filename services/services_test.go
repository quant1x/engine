package services

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gotdx"
	"testing"
)

func TestGlobalReset(t *testing.T) {
	_ = cleanExpiredStateFiles()
	gotdx.ReOpen()
	date := cache.DefaultCanUpdateDate()
	factors.SwitchDate(date)
}
