package services

import (
	"testing"

	"github.com/quant1x/engine/cache"
	"github.com/quant1x/engine/factors"
	"github.com/quant1x/gotdx"
)

func TestGlobalReset(t *testing.T) {
	_ = cleanExpiredStateFiles()
	gotdx.ReOpen()
	date := cache.DefaultCanUpdateDate()
	factors.SwitchDate(date)
}
