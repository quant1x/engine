package models

import (
	"fmt"
	"testing"

	"github.com/quant1x/engine/cache"
	"github.com/quant1x/engine/factors"
	"github.com/quant1x/exchange"
	"github.com/quant1x/x/api"
)

func TestFeatureToSnapshot(t *testing.T) {
	code := "300410"
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
