package factors

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/cache"
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
	misc := NewMisc(code, date)
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
