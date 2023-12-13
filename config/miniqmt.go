package config

import "fmt"

// QmtStrategyNameFromId 通过策略ID返回用于在QMT系统中表示的string类型的策略名称
func QmtStrategyNameFromId(strategyCode int) string {
	return fmt.Sprintf("S%d", strategyCode)
}
