package services

import (
	"testing"

	"github.com/quant1x/engine/models"
)

func TestRealtimeUpdateExchangeAndSnapshot(t *testing.T) {
	models.SyncAllSnapshots(nil)
	realtimeUpdateMiscAndSnapshot()
}
