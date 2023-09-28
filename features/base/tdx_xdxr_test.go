package base

import (
	"fmt"
	"testing"
)

func TestGetCacheXdxrList(t *testing.T) {
	code := "sz002043"
	list := GetCacheXdxrList(code)
	fmt.Println(list)
}
