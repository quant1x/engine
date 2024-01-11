package factors

// SecurityFeature 证券特征信息
type SecurityFeature struct {
	Date           string  `name:"日期" dataframe:"date,string"`
	Open           float64 `name:"开盘" dataframe:"open,float64"`
	Close          float64 `name:"收盘" dataframe:"close,float64"`
	High           float64 `name:"最高" dataframe:"high,float64"`
	Low            float64 `name:"最低" dataframe:"low,float64"`
	Volume         int64   `name:"成交量" dataframe:"volume,int64"`
	Amount         float64 `name:"成交额" dataframe:"amount,float64"`
	Up             int     `name:"上涨家数" dataframe:"up,int64"`
	Down           int     `name:"下跌家数" dataframe:"down,int64"`
	LastClose      float64 `name:"昨收" dataframe:"last_close,float64"`
	ChangeRate     float64 `name:"涨跌幅" dataframe:"change_rate,float64"`
	OpenVolume     int64   `name:"开盘量" dataframe:"open_volume,int64"`
	OpenTurnZ      float64 `name:"开盘换手z" dataframe:"open_turnz,float64"`
	OpenUnmatched  int64   `name:"开盘未匹配" dataframe:"open_unmatched,int64"`
	CloseVolume    int64   `name:"收盘量" dataframe:"close_volume,int64"`
	CloseTurnZ     float64 `name:"收盘换手z" dataframe:"close_turnz,float64"`
	CloseUnmatched int64   `name:"收盘未匹配" dataframe:"close_unmatched,int64"`
	InnerVolume    int64   `name:"内盘" dataframe:"inner_volume,int64"`
	OuterVolume    int64   `name:"外盘" dataframe:"outer_volume,int64"`
	InnerAmount    float64 `name:"流出金额" dataframe:"inner_amount,float64"`
	OuterAmount    float64 `name:"流入金额" dataframe:"outer_amount,float64"`
}

// TurnoverDataSummary 换手数据概要
type TurnoverDataSummary struct {
	OpenVolume     int64   `name:"开盘量" dataframe:"open_volume,int64"`
	OpenTurnZ      float64 `name:"开盘换手z" dataframe:"open_turnz,float64"`
	OpenUnmatched  int64   `name:"开盘未匹配" dataframe:"open_unmatched,int64"`
	CloseVolume    int64   `name:"收盘量" dataframe:"close_volume,int64"`
	CloseTurnZ     float64 `name:"收盘换手z" dataframe:"close_turnz,float64"`
	CloseUnmatched int64   `name:"收盘未匹配" dataframe:"close_unmatched,int64"`
	InnerVolume    int64   `name:"内盘" dataframe:"inner_volume,int64"`
	OuterVolume    int64   `name:"外盘" dataframe:"outer_volume,int64"`
	InnerAmount    float64 `name:"流出金额" dataframe:"inner_amount,float64"`
	OuterAmount    float64 `name:"流入金额" dataframe:"outer_amount,float64"`
}
