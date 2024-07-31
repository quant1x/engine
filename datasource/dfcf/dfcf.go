package dfcf

// 原始的接口状态信息
type eastMoneyStatus struct {
	Version string `json:"version"` // 版本
	Success bool   `json:"success"` // 是否成功
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}

// 原始api结果结构体
type eastMoneyData[T any] struct {
	Pages int `json:"pages"`
	Count int `json:"count"`
	Data  []T `json:"data"`
}

// 原始的api响应结果
type rawResult[T any] struct {
	eastMoneyStatus
	eastMoneyData[T] `json:"result"`
}
