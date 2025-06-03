package tdxweb

import (
	"fmt"
	"testing"
)

func TestGetSafetyScore(t *testing.T) {
	code := "sh510050"
	code = "sh000001"
	code = "sh600178"
	v := GetSafetyScore(code)
	fmt.Println(v)
}
