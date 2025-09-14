package trader

import (
	"testing"

	"gitee.com/quant1x/engine/models"
)

func TestCacheSync(t *testing.T) {
	barIndex := 1
	models.SyncAllSnapshots(&barIndex)
	//UpdatePositions()
	SyncPositions()
	CacheSync()
}
