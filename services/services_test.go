package services

import (
	"testing"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gotdx"
)

func TestGlobalReset(t *testing.T) {
	_ = cleanExpiredStateFiles()
	gotdx.ReOpen()
	date := cache.DefaultCanUpdateDate()
	factors.SwitchDate(date)
}
