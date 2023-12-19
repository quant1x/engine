package base

import (
	"fmt"
	"testing"
)

func TestUpdateXdxrInfo(t *testing.T) {
	code := "sh603158"
	UpdateXdxrInfo(code)
}

func TestGetCacheXdxrList(t *testing.T) {
	code := "sz002043"
	code = "sh000001"
	code = "sh603158"
	list := GetCacheXdxrList(code)
	fmt.Println(list)
}
