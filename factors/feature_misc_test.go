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
	date := "2024-01-16"
	cacheDate, featureDate := cache.CorrectDate(date)
	misc := NewMisc(code, date)
	misc.Update(code, cacheDate, featureDate, false)
	data, _ := json.Marshal(misc)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
