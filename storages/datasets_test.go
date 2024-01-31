package storages

import (
	"gitee.com/quant1x/engine/cache"
	"testing"
)

func TestBaseDataUpdate(t *testing.T) {
	barIndex := 1
	date := "2024-01-31"
	plugins := cache.PluginsWithName(cache.PluginMaskBaseData, "wide")
	BaseDataUpdate(barIndex, date, plugins, cache.OpUpdate)

}
