package factors

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestFeatureMisc(t *testing.T) {
	code := "sh880652"
	code = "sz300904"
	code = "sh603038"
	code = "sh600178"
	code = "sz300261"
	date := "2024-03-27"
	cacheDate, featureDate := cache.CorrectDate(date)
	misc := NewMisc(date, code)
	misc.Update(code, cacheDate, featureDate, false)
	fmt.Println(misc.Shape & KLineShapeDoji)
	data, _ := json.Marshal(misc)
	text := api.Bytes2String(data)
	fmt.Println(text)
}

func TestMisc_MarginTradingTargets(t *testing.T) {
	date := cache.DefaultCanReadDate()
	MarginTradingTargetInit(date)
	code := "000099"
	v, ok := GetMarginTradingTarget(code)
	fmt.Println(v, ok)
}

func TestMisc_AuctionWeaknessToStrength(t *testing.T) {
	date := cache.DefaultCanReadDate()
	code := "300107"
	code = "603679"
	date = "2024-04-19"
	code = exchange.CorrectSecurityCode(code)
	misc := GetL5Misc(code, date)
	v := misc.AuctionWeaknessToStrength()
	fmt.Println(v)
}

func TestMisc_AuctionStrengthToWeakness(t *testing.T) {
	date := cache.DefaultCanReadDate()
	code := "000612"
	//code = "000099"
	date = "2024-04-19"
	code = exchange.CorrectSecurityCode(code)
	misc := GetL5Misc(code, date)
	v := misc.AuctionStrengthToWeakness()
	fmt.Println(v)
}
