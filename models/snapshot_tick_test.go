package models

import "testing"

func TestSyncAllSnapshots(t *testing.T) {
	barIndex := 1
	SyncAllSnapshots(&barIndex)
}
