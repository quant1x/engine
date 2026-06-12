package services

import (
	"testing"

	"github.com/quant1x/data/level1"
	"github.com/quant1x/engine/cache"
	"github.com/quant1x/engine/factors"
)

func TestGlobalReset(t *testing.T) {
	_ = cleanExpiredStateFiles()
	level1.ReOpen()
	date := cache.DefaultCanUpdateDate()
	factors.SwitchDate(date)
}
