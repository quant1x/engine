package services

import (
	"gitee.com/quant1x/engine/models"
	"testing"
)

func TestRealtimeUpdateExchangeAndSnapshot(t *testing.T) {
	models.SyncAllSnapshots(nil)
	realtimeUpdateMiscAndSnapshot()
}
