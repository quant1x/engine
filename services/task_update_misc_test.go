package services

import (
	"testing"

	"gitee.com/quant1x/engine/models"
)

func TestRealtimeUpdateExchangeAndSnapshot(t *testing.T) {
	models.SyncAllSnapshots(nil)
	realtimeUpdateMiscAndSnapshot()
}
