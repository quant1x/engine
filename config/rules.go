package config

// RuleParameter 规则参数
type RuleParameter struct {
	SectorsFilter               bool    `yaml:"sectors_filter" default:"false"`                  // 是否启用板块过滤, false代表全市场扫描
	SectorsTopN                 int     `yaml:"sectors_top_n" default:"3"`                       // 最多关联多少个板块, 默认3个
	StockTopNInSector           int     `yaml:"stock_top_n_in_sector" default:"5"`               // 板块内个股排名前N
	PriceMin                    float64 `yaml:"price_min" default:"2.00"`                        // 2.00 股价最低
	PriceMax                    float64 `yaml:"price_max" default:"30.00"`                       // 30.00 股价最高
	MaximumIncreaseWithin5days  float64 `yaml:"maximum_increase_within_5_days" default:"20.00"`  // 20.00 5日累计最大涨幅
	MaximumIncreaseWithin10days float64 `yaml:"maximum_increase_within_10_days" default:"70.00"` // 70.00 10日累计最大涨幅
	MaxReduceAmount             float64 `yaml:"max_reduce_amount" default:"-1000"`               // -1000 最大流出1000万
	TurnZMax                    float64 `yaml:"turn_z_max" default:"200.00"`                     // 200.00 换手最大值
	TurnZMin                    float64 `yaml:"turn_z_min" default:"1.50"`                       // 1.50 换手最小值
	OpenRateMax                 float64 `yaml:"open_rate_max"  default:"2.00"`                   // 2.00 最大涨幅
	OpenRateMin                 float64 `yaml:"open_rate_min" default:"-2.00"`                   // -2.00 最低涨幅
	QuantityRatioMax            float64 `yaml:"quantity_ratio_max" default:"9.99"`               // 9.99 最大开盘量比
	QuantityRatioMin            float64 `yaml:"quantity_ratio_min"  default:"1.00"`              // 1.00 最小开盘量比
	SafetyScoreMin              float64 `yaml:"safety_score_min" default:"80"`                   // 80 通达信安全分最小值
	VolumeRatioMax              float64 `yaml:"volume_ratio_max" default:"3.82"`                 // 1.800 成交量放大不能超过1.8
	CapitalMin                  float64 `yaml:"capital_min" default:"2"`                         // 2 * 亿 流通股本最小值
	CapitalMax                  float64 `yaml:"capital_max" default:"20"`                        // 20 * 亿 流通股本最大值
	VixMax                      float64 `yaml:"vix_max" default:"100"`                           // 波动率最大值100
	VixMin                      float64 `yaml:"vix_min" default:"0"`                             // 波动率最小值0
	TurnoverRateMax             float64 `yaml:"turnover_rate_max" default:"20.00"`               // 换手率最大20%
	TurnoverRateMin             float64 `yaml:"turnover_rate_min" default:"1.00"`                // 换手率最小1%
	AmplitudeRatioMax           float64 `yaml:"amplitude_ratio_max" default:"15"`                // 振幅 最大
	AmplitudeRatioMin           float64 `yaml:"amplitude_ratio_min" default:"0"`                 // 振幅 最小
	BiddingVolumeMax            int     `yaml:"bidding_volume_max" default:"5000"`               // 5档行情委托平均最大值
	BiddingVolumeMin            int     `yaml:"bidding_volume_min" default:"100"`                // 5档行情委托平均最小值
	SentimentHigh               float64 `yaml:"sentiment_high" default:"61.8"`                   // 情绪值最高
	SentimentLow                float64 `yaml:"sentiment_low" default:"38.2"`                    // 情绪值最低
}
