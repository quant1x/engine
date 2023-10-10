package storages

import (
	"gitee.com/quant1x/engine/cache"
	"testing"
)

func TestPluginsRepairBase(t *testing.T) {
	barIndex := 1
	date := "20230928"
	cacheDate, featureDate := cache.CorrectDate(date)
	RepairBaseData(&barIndex, cacheDate, featureDate)
}

func TestPluginsRepairFeatures(t *testing.T) {
	barIndex := 1
	date := "20231009"
	cacheDate, featureDate := cache.CorrectDate(date)
	RepairFeatures(&barIndex, cacheDate, featureDate)
}