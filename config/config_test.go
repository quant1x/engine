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
	v := GetStrategyParameterByCode(strategyCode)
	fmt.Println(v)
}
