package storages

import (
	"testing"

	"github.com/quant1x/engine/cache"
)

func TestBaseDataUpdate(t *testing.T) {
	barIndex := 1
	date := "2024-01-31"
	plugins := cache.PluginsWithName(cache.PluginMaskBaseData, "wide")
	DataSetUpdate(barIndex, date, plugins, cache.OpUpdate)
}
