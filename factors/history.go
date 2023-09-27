package factors

// History 历史整合数据
type History struct {
	Date       string  // 日期, 数据落地的日期
	Code       string  // 代码
	MA3        float64 // 3日均价
	MV3        float64 // 3日均量
	MA5        float64 // 5日均价
	MV5        float64 // 5日均量
	MA10       float64 // 10日均价
	MV10       float64 // 10日均量
	MA20       float64 // 20日均价
	MV20       float64 // 20日均量
	UpdateTime string  // 更新时间
}
