package global

import (
	"gitee.com/quant1x/engine/db"
	"gitee.com/quant1x/engine/market"
	"sync"
)

type Variables struct {
	db         *db.Database
	MarketData *market.MarketData
}

var (
	variables Variables
	once      sync.Once
)

func GetGlobalVariables() *Variables {
	once.Do(func() {
		database := db.InitDatabase()
		marketData := market.NewMarketData(database)
		variables.db = database
		variables.MarketData = marketData
	})
	return &variables
}
