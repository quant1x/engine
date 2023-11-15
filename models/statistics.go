package models

// Statistics 0号策略数据, 订单结构
type Statistics struct {
	Date                 string  `name:"日期" dataframe:"date"`
	Code                 string  `name:"证券代码" dataframe:"code"`
	Name                 string  `name:"证券名称" dataframe:"name"`
	TurnZ                float64 `name:"开盘换手Z%" dataframe:"turnz"`
	QuantityRatio        float64 `name:"开盘量比" dataframe:"quantity_ratio"`
	Tendency             string  `name:"趋势" dataframe:"tendency"`
	LastClose            float64 `name:"昨收" dataframe:"last_close"`
	Open                 float64 `name:"开盘价" dataframe:"open"`
	OpenRaise            float64 `name:"开盘涨幅%" dataframe:"open_raise"`
	Price                float64 `name:"现价" dataframe:"price"`
	UpRate               float64 `name:"涨跌幅%" dataframe:"up_rate"`
	OpenPremiumRate      float64 `name:"浮动溢价率%" dataframe:"open_premium_rate"`
	NextPremiumRate      float64 `name:"隔日溢价率%" dataframe:"next_premium_rate"`
	BlockName            string  `name:"板块名称" dataframe:"block_name"`
	BlockRate            float64 `name:"板块涨幅%" dataframe:"block_rate"`
	BlockTop             int     `name:"板块排名" dataframe:"block_top"`
	BlockRank            int     `name:"个股排名" dataframe:"block_rank"`
	OpenVolume           int     `name:"开盘量" dataframe:"open_volume"`
	AveragePrice         float64 `name:"均价线" dataframe:"average_price"`
	Active               int     `name:"活跃度" dataframe:"active"`
	ChangePower          float64 `name:"力度" dataframe:"change_power"`
	AverageBiddingVolume int     `name:"委托均量" dataframe:"average_bidding_volume"` // 委托均量
	UpdateTime           string  `name:"时间戳" dataframe:"update_time"`
}
