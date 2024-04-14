package models

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestFeatureToSnapshot(t *testing.T) {
	code := "600105"
	securityCode := exchange.CorrectSecurityCode(code)
	filename := cache.WideFilename(securityCode)
	features := []factors.SecurityFeature{}
	err := api.CsvToSlices(filename, &features)
	if err != nil {
		fmt.Println(err)
		return
	}
	length := len(features)
	feature := features[length-1]
	v := FeatureToSnapshot(feature, securityCode)
	fmt.Println(v)
}
