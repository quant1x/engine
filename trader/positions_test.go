package trader

import (
	"gitee.com/quant1x/engine/models"
	"testing"
)

func TestCacheSync(t *testing.T) {
	barIndex := 1
	models.SyncAllSnapshots(&barIndex)
	//UpdatePositions()
	SyncPositions()
	CacheSync()
}
