package models

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gox/num"
)

// FeatureToSnapshot 特征缓存数据转快照
func FeatureToSnapshot(feature cache.SecurityFeature, securityCode string) factors.QuoteSnapshot {
	qs := factors.QuoteSnapshot{
		//Market            uint8   // 市场
		//Code              string  `name:"证券代码"`  // 代码
		Code: securityCode,
		//Name              string  `name:"证券名称"`  // 证券名称
		//Active            uint16  `name:"活跃度"`   // 活跃度
		//LastClose         float64 `name:"昨收"`    // 昨收
		LastClose: feature.LastClose,
		//Open              float64 `name:"开盘价"`   // 开盘
		Open: feature.Open,
		//OpeningChangeRate float64 `name:"开盘涨幅%"` // 开盘
		OpeningChangeRate: num.NetChangeRate(feature.LastClose, feature.Open),
		//Price             float64 `name:"现价"`    // 现价
		Price: feature.Close,
		//ChangeRate        float64 `name:"涨跌幅%"`  // 涨跌幅
		ChangeRate: num.NetChangeRate(feature.LastClose, feature.Close),
		//PremiumRate       float64 `name:"溢价率%"`  // 集合竞价买入溢价, hedgeRatio
		PremiumRate: num.NetChangeRate(feature.Open, feature.Close),
		//High              float64 // 最高
		High: feature.High,
		//Low               float64 // 最低
		Low: feature.Low,
		//ServerTime        string  // 时间
		//ReversedBytes0    int     // 保留(时间 ServerTime)
		//ReversedBytes1    int     // 保留
		//Vol               int     // 总量
		Vol: int(feature.Volume),
		//CurVol            int     // 现量
		//Amount            float64 // 总金额
		Amount: feature.Amount,
		//SVol              int     // 内盘
		SVol: int(feature.InnerVolume),
		//BVol              int     // 外盘
		BVol: int(feature.OuterVolume),
		//IndexOpenAmount   int     // 指数-集合竞价成交金额=开盘成交金额
		//StockOpenAmount   int     // 个股-集合竞价成交金额=开盘成交金额
		StockOpenAmount: int(float64(feature.OpenVolume) * feature.Open),
		//OpenVolume        int     `name:"开盘量"` // 集合竞价-开盘量, 单位是股
		OpenVolume: int(feature.OpenVolume),
		//Bid1              float64
		//Ask1              float64
		//BidVol1           int
		//AskVol1           int
		//Bid2              float64
		//Ask2              float64
		//BidVol2           int
		//AskVol2           int
		//Bid3              float64
		//Ask3              float64
		//BidVol3           int
		//AskVol3           int
		//Bid4              float64
		//Ask4              float64
		//BidVol4           int
		//AskVol4           int
		//Bid5              float64
		//Ask5              float64
		//BidVol5           int
		//AskVol5           int
		//ReversedBytes4    uint16  // 保留
		//ReversedBytes5    int     // 保留
		//ReversedBytes6    int     // 保留
		//ReversedBytes7    int     // 保留
		//ReversedBytes8    int     // 保留
		//Rate              float64 // 涨速
		//Active2           uint16  // 活跃度
		//TopNo             int     // 板块排名
		//TopCode           string  // 领涨个股
		//TopName           string  // 领涨个股名称
		//TopRate           float64 // 领涨个股涨幅
		//ZhanTing          int     // 涨停数
		//Ling              int     // 平盘数
		//Count             int     // 总数
		//Capital           float64 `name:"流通盘"`    // 流通盘
		//FreeCapital       float64 `name:"自由流通股本"` // 自由流通股本
		//OpenTurnZ         float64 `name:"开盘换手Z%"` // 开盘换手
		//OpenQuantityRatio     float64 `name:"开盘量比"`
	}
	f10 := factors.GetL5F10(securityCode)
	if f10 != nil {
		qs.Capital = f10.Capital
		qs.FreeCapital = f10.FreeCapital
		qs.OpenTurnZ = 10000 * num.ChangeRate(qs.FreeCapital, qs.OpenVolume)
	}
	extend := factors.GetL5History(securityCode)
	if extend != nil && extend.MV5 > 0 {
		qs.OpenQuantityRatio = num.ChangeRate(extend.MV5, qs.OpenVolume)
	}
	return qs
}
