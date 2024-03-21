package trader

import (
	"testing"
	"time"
)

func TestProhibitTradingToBlackList(t *testing.T) {
	code := "sh000001"
	ProhibitTradingToBlackList(code)
	time.Sleep(1 * time.Second)
	ProhibitTradingToBlackList("sh000002")
	time.Sleep(20 * time.Second)
}

func TestAddCodeToBlackList(t *testing.T) {
	code := "sh603230"
	secureType := FreeTrading
	AddCodeToBlackList(code, secureType)
}
