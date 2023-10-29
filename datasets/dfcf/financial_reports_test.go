package dfcf

import (
	"fmt"
	"testing"
)

func TestQuarterlyReports(t *testing.T) {
	date := "20230928"

	list, pages, err := QuarterlyReports(date, 1)
	fmt.Println(list)
	fmt.Println(pages)
	fmt.Println(err)
}

func TestGetCacheQuarterlyReportsBySecurityCode(t *testing.T) {
	date := "20231027"
	code := "sz002685"
	v := GetCacheQuarterlyReportsBySecurityCode(code, date)
	fmt.Println(v)
}
