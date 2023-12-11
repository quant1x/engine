package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	strategyCode := 82
	v := GetTradeRule(strategyCode)
	fmt.Println(v)
}
