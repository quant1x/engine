package config

import (
	"fmt"
	"gitee.com/quant1x/gotdx/securities"
	"testing"
	"time"
)

import (
	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	config, found := LoadConfig()
	fmt.Println(found)
	fmt.Println(config)
	strategyCode := 82
	v := GetStrategyParameterByCode(uint64(strategyCode))
	fmt.Println(v)
}

func TestBlocks(t *testing.T) {
	sectorCode := "sh880884"
	blk := securities.GetBlockInfo(sectorCode)
	fmt.Println(len(blk.ConstituentStocks))
}

type TradingPeriodYaml struct {
	StartTime string `yaml:"start_time"`
	EndTime   string `yaml:"end_time"`
	OrderType string `yaml:"order_type"`
}

type MarketYaml struct {
	Name           string              `yaml:"name"`
	Timezone       string              `yaml:"timezone"`
	TradingPeriods []TradingPeriodYaml `yaml:"trading_periods"`
}

type TradeConfigYaml struct {
	MarketType string `yaml:"market_type"`
}

type Config struct {
	Markets []MarketYaml    `yaml:"markets"`
	Trade   TradeConfigYaml `yaml:"trade"`
}

type TradingPeriod struct {
	StartTime time.Time
	EndTime   time.Time
	OrderType string `yaml:"order_type"`
}

type Market struct {
	Name           string
	TradingPeriods []TradingPeriod
}

func MarketYamlToMarket(marketYaml *MarketYaml, date time.Time) Market {
	tradingPeriods := make([]TradingPeriod, 0)

	for _, periodYaml := range marketYaml.TradingPeriods {
		period := tradingPeriodYamlToTradingPeriod(periodYaml, marketYaml.Timezone, date)
		tradingPeriods = append(tradingPeriods, period)
	}

	market := Market{
		Name:           marketYaml.Name,
		TradingPeriods: tradingPeriods,
	}

	return market
}
func tradingPeriodYamlToTradingPeriod(tradingPeriodYaml TradingPeriodYaml, timezone string, date time.Time) TradingPeriod {

	location, err := time.LoadLocation(timezone)
	if err != nil {
		fmt.Println("无法加载时区:", err)
		panic("failed")
	}

	startTime, err := time.ParseInLocation("15:04", tradingPeriodYaml.StartTime, location)
	if err != nil {
		fmt.Println("无法解析开始时间:", err)
		panic("failed")
	}

	endTime, err := time.ParseInLocation("15:04", tradingPeriodYaml.EndTime, location)
	if err != nil {
		fmt.Println("无法解析结束时间:", err)
		panic("failed")
	}

	// 获取今天的日期

	// 组合得到今天的开始时间和结束时间
	startDateTime := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, location)
	endDateTime := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, location)

	tradingPeriod := TradingPeriod{
		StartTime: startDateTime,
		EndTime:   endDateTime,
		OrderType: tradingPeriodYaml.OrderType,
	}
	return tradingPeriod
}

func (m *Market) FindTradingPeriod(timestamp int64, timezone string) *TradingPeriod {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		fmt.Println("无法加载时区:", err)
		panic("failed")
	}

	t := time.Unix(timestamp, 0).In(location)

	for _, period := range m.TradingPeriods {
		if t.After(period.StartTime) && t.Before(period.EndTime) {
			return &period
		}
	}
	return nil
}

func TestYaml(t *testing.T) {
	yamlString := `
markets:
  - name: "股票市场"
    timezone: "Asia/Shanghai"
    trading_periods:
      - start_time: "09:15"
        end_time: "09:20"
        order_type: "可委托可撤单"
      - start_time: "09:20"
        end_time: "09:25"
        order_type: "可委托不可撤单"
      - start_time: "09:25"
        end_time: "09:30"
        order_type: "可委托可撤单"
      - start_time: "09:30"
        end_time: "11:30"
        order_type: "可委托可撤单"
      - start_time: "13:00"
        end_time: "14:57"
        order_type: "可委托可撤单"
      - start_time: "14:57"
        end_time: "15:00"
        order_type: "可委托不可撤单"

trade:
  market_type: "股票市场"
`

	// 解析YAML字符串
	var config Config
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		fmt.Println("解析YAML失败:", err)
		return
	}

	// 输出市场类型和交易时段
	for _, market := range config.Markets {
		fmt.Println("市场类型:", market.Name)
		fmt.Println("时区:", market.Timezone)
		fmt.Println("交易时段:")
		for _, period := range market.TradingPeriods {
			fmt.Println("开始时间:", period.StartTime)
			fmt.Println("结束时间:", period.EndTime)
			fmt.Println("委托类型:", period.OrderType)
		}
		fmt.Println()
	}

	// 输出当前交易配置的市场类型
	fmt.Println("交易配置中的市场类型:", config.Trade.MarketType)

	market := MarketYamlToMarket(&config.Markets[0], time.Now())

	now := time.Now()
	timestamp := time.Date(now.Year(), now.Month(), now.Day(), 9, 18, 0, 0, now.Location()).Unix()

	tradingPeriod := market.FindTradingPeriod(timestamp, "Asia/Shanghai")
	if tradingPeriod != nil {
		fmt.Println("开始时间:", tradingPeriod.StartTime)
		fmt.Println("结束时间:", tradingPeriod.EndTime)
		fmt.Println("订单类型:", tradingPeriod.OrderType)
	} else {
		fmt.Println("未找到匹配的交易时间段")
	}
}
