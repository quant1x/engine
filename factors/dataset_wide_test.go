package factors

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestKLine(t *testing.T) {
	code := "002882"
	securityCode := exchange.CorrectSecurityCode(code)
	var wides []SecurityFeature
	filename := cache.KLineFilename(securityCode)
	err := api.CsvToSlices(filename, &wides)
	if err != nil || len(wides) == 0 {
		return
	}
	wides = wides[len(wides)-1:]
	df := pandas.LoadStructs(wides)
	fmt.Println(df)
}

func TestKLineWide(t *testing.T) {
	code := "002882"
	code = "600178"
	securityCode := exchange.CorrectSecurityCode(code)
	var wides []SecurityFeature
	filename := cache.FeatureFilename(securityCode)
	err := api.CsvToSlices(filename, &wides)
	if err != nil || len(wides) == 0 {
		return
	}
	//wides = wides[len(wides)-1:]
	df := pandas.LoadStructs(wides)
	fmt.Println(df)
}
