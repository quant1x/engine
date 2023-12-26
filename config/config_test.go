package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	config, found := LoadConfig()
	fmt.Println(found)
	fmt.Println(config)
	strategyCode := 82
	v := GetTradeRule(strategyCode)
	fmt.Println(v)
}
