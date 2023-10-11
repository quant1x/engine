package cache

// SecurityFeature 证券特征信息
type SecurityFeature struct {
	Date           string  `json:"Date" array:"0" name:"日期" dataframe:"date,string"`
	Open           float64 `json:"Open" array:"1" name:"开盘" dataframe:"open,float64"`
	Close          float64 `json:"Close" array:"2" name:"收盘" dataframe:"close,float64"`
	High           float64 `json:"High" array:"3" name:"最高" dataframe:"high,float64"`
	Low            float64 `json:"Low" array:"4" name:"最低" dataframe:"low,float64"`
	Volume         int64   `json:"Volume" array:"5" name:"成交量" dataframe:"volume,int64"`
	Amount         float64 `json:"Amount" array:"6" name:"成交额" dataframe:"amount,float64"`
	Up             int     `json:"Up" array:"7" name:"上涨家数" dataframe:"up,int64"`
	Down           int     `json:"Down" array:"8" name:"下跌家数" dataframe:"down,int64"`
	LastClose      float64 `json:"LastClose" array:"9" name:"昨收" dataframe:"last_close,float64"`
	ChangeRate     float64 `json:"ChangeRate" array:"10" name:"涨跌幅" dataframe:"change_rate,float64"`
	OpenVolume     int64   `json:"OpenVolume" array:"11" name:"开盘量" dataframe:"open_volume,int64"`
	OpenTurnZ      float64 `json:"OpenTurnZ" array:"12" name:"开盘换手z" dataframe:"open_turnz,float64"`
	OpenUnmatched  int64   `json:"OpenUnmatched" array:"13" name:"开盘未匹配" dataframe:"open_unmatched,int64"`
	CloseVolume    int64   `json:"CloseVolume" array:"14" name:"收盘量" dataframe:"close_volume,int64"`
	CloseTurnZ     float64 `json:"CloseTurnZ" array:"15" name:"收盘换手z" dataframe:"close_turnz,float64"`
	CloseUnmatched int64   `json:"CloseUnmatched" array:"16" name:"收盘未匹配" dataframe:"close_unmatched,int64"`
	InnerVolume    int64   `json:"InnerVolume" array:"17" name:"内盘" dataframe:"inner_volume,int64"`
	OuterVolume    int64   `json:"OuterVolume" array:"18" name:"外盘" dataframe:"outer_volume,int64"`
	InnerAmount    float64 `json:"InnerAmount" array:"19" name:"流出金额" dataframe:"inner_amount,float64"`
	OuterAmount    float64 `json:"OuterAmount" array:"20" name:"流入金额" dataframe:"outer_amount,float64"`
}

// TurnoverDataSummary 换手数据概要
type TurnoverDataSummary struct {
	OpenVolume     int64   `json:"OpenVolume" array:"0" name:"开盘量" dataframe:"open_volume,int64"`
	OpenTurnZ      float64 `json:"OpenTurnZ" array:"1" name:"开盘换手z" dataframe:"open_turnz,float64"`
	OpenUnmatched  int64   `json:"OpenUnmatched" array:"2" name:"开盘未匹配" dataframe:"open_unmatched,int64"`
	CloseVolume    int64   `json:"CloseVolume" array:"3" name:"收盘量" dataframe:"close_volume,int64"`
	CloseTurnZ     float64 `json:"CloseTurnZ" array:"4" name:"收盘换手z" dataframe:"close_turnz,float64"`
	CloseUnmatched int64   `json:"CloseUnmatched" array:"5" name:"收盘未匹配" dataframe:"close_unmatched,int64"`
	InnerVolume    int64   `json:"InnerVolume" array:"6" name:"内盘" dataframe:"inner_volume,int64"`
	OuterVolume    int64   `json:"OuterVolume" array:"7" name:"外盘" dataframe:"outer_volume,int64"`
	InnerAmount    float64 `json:"InnerAmount" array:"8" name:"流出金额" dataframe:"inner_amount,float64"`
	OuterAmount    float64 `json:"OuterAmount" array:"9" name:"流入金额" dataframe:"outer_amount,float64"`
}
