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
	filename := cache.WideFilename(securityCode)
	err := api.CsvToSlices(filename, &wides)
	if err != nil || len(wides) == 0 {
		return
	}
	//wides = wides[len(wides)-1:]
	df := pandas.LoadStructs(wides)
	fmt.Println(df)
}

func TestGetKLineAll(t *testing.T) {
	code := "600178"
	securityCode := exchange.CorrectSecurityCode(code)
	df := GetKLineAll(securityCode)
	fmt.Println(df)
}

func Test_pullWideByDate(t *testing.T) {
	code := "880424"
	code = "sh000001"
	date := "20231009"
	securityCode := exchange.CorrectSecurityCode(code)
	data := pullWideByDate(securityCode, date)
	fmt.Println(data)
}

func TestWideTableValuate(t *testing.T) {
	code := "881432"
	date := "20240116"
	lines := CheckoutWideKLines(code, date)
	df := pandas.LoadStructs(lines)
	fmt.Println(df)
}
