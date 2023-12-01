package config

// OrderParameter 订单参数
type OrderParameter struct {
	AccountId string `yaml:"account_id"`               // 账号ID
	OrderPath string `yaml:"order_path"`               // 订单路径
	TopN      int    `yaml:"top_n" default:"3"`        // 最多输出前多少名个股
	HaveETF   bool   `yaml:"have_etf" default:"false"` // 是否包含ETF
}
