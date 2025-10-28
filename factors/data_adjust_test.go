package factors

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/engine/datasource/base"
)

// 2. 完善接口实现
type testCase struct {
	date  string
	price float64
}

func (tc *testCase) Apply(factor func(float64) float64) {
	tc.price = factor(tc.price) // 实际实现调整逻辑
}

func (tc *testCase) GetDate() string {
	return tc.date
}

func apply0[E any](securityCode string, data []E, startDate string) {
	for i := 0; i < len(data); i++ {
		f, ok := any(&data[i]).(AdjustmentExecutor)
		if !ok {
			continue
		}
		if f.GetDate() > startDate {
			f.Apply(func(p float64) float64 {
				return p * 0.5 // 示例调整因子
			})
		}
	}
	_ = securityCode
	_ = data
	_ = startDate
}

func TestApplyAdjustment0(t *testing.T) {
	code := "sh000001"
	date := "2025-03-10"
	securityCode := code
	//securityCode := exchange.CorrectSecurityCode(code)
	//date = exchange.FixTradeDate(date)
	//klines := base.CheckoutKLines(securityCode, date)
	//fmt.Println(klines)
	testCases := []testCase{
		{
			date:  "2025-03-11",
			price: 1.0,
		},
		{
			date:  "2025-03-12",
			price: 2.0,
		},
	}
	fmt.Println(testCases)
	// 执行测试
	apply0(securityCode, testCases, date)
	fmt.Println("------------------------------------------------------------")
	fmt.Println(testCases)
}

func TestApplyAdjustment(t *testing.T) {
	code := "600580"
	date := "2025-03-10"
	securityCode := exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	klines := base.CheckoutKLines(securityCode, date)
	fmt.Println(klines)
	// 执行测试
	ApplyAdjustment(securityCode, klines, date)
	fmt.Println("------------------------------------------------------------")
	fmt.Println(klines)
}
