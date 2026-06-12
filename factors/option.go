package factors

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/quant1x/gox/api"
)

// OptionFinanceBoardData 表示期权行情数据
type OptionFinanceBoardData struct {
	Date        string  `json:"日期"`
	ContractID  string  `json:"合约交易代码"`
	Price       float64 `json:"当前价"`
	ChangeRate  float64 `json:"涨跌幅"`
	PrevSettle  float64 `json:"前结价"`
	StrikePrice float64 `json:"行权价"`
	Quantity    int     `json:"数量"`
}

// RiskIndicator 上交所风险指标
//
// 合约信息:
//   - Delta值(δ): 又称对冲值，是衡量标的资产价格变动时，期权价格的变化幅度。
//   - Gamma(γ): 反映期货价格对delta值的影响程度，为delta变化量与期货价格变化量之比。
//   - vega值: 认股证对引伸波幅变动的敏感度，期权的风险指标通常用希腊字母来表示。
//   - Theta(θ): 是用来测量时间变化对期权理论价值的影响。表示时间每经过一天，期权价值会损失多少。
//   - Rho: 是指期权价格对无风险利率变化的敏感程度，是用以衡量利率转变对权证价值影响的指针。
type RiskIndicator struct {
	TradeDate       time.Time `json:"TRADE_DATE"`
	SecurityID      string    `json:"SECURITY_ID"`
	ContractID      string    `json:"CONTRACT_ID"`
	ContractSymbol  string    `json:"CONTRACT_SYMBOL"`
	Delta           float64   `json:"DELTA_VALUE"`
	Theta           float64   `json:"THETA_VALUE"`
	Gamma           float64   `json:"GAMMA_VALUE"`
	Vega            float64   `json:"VEGA_VALUE"`
	Rho             float64   `json:"RHO_VALUE"`
	ImplcVolatility float64   `json:"IMPLC_VOLATLTY"`
}

// HTTP 客户端
var client = &http.Client{Timeout: 10 * time.Second}

// ==================== 1. 期权行情数据: option_finance_board ====================

// OptionFinanceBoard 期权行情数据
func OptionFinanceBoard(symbol string, endMonth string) ([]OptionFinanceBoardData, error) {
	endMonth = endMonth[len(endMonth)-2:] // 取最后两位

	var optionUrl string
	var payload = url.Values{
		"select": {"contractid,last,chg_rate,presetpx,exepx"},
	}

	switch symbol {
	case "华夏上证50ETF期权":
		optionUrl = "http://yunhq.sse.com.cn:32041/v1/sho/list/tstyle/510050_" + endMonth
	case "华泰柏瑞沪深300ETF期权":
		optionUrl = "http://yunhq.sse.com.cn:32041/v1/sho/list/tstyle/510300_" + endMonth
	case "南方中证500ETF期权":
		optionUrl = "http://yunhq.sse.com.cn:32041/v1/sho/list/tstyle/510500_" + endMonth
	case "华夏科创50ETF期权":
		optionUrl = "http://yunhq.sse.com.cn:32041/v1/sho/list/tstyle/588000_" + endMonth
	case "易方达科创50ETF期权":
		optionUrl = "http://yunhq.sse.com.cn:32041/v1/sho/list/tstyle/588080_" + endMonth
	default:
		return nil, fmt.Errorf("不支持的 symbol: %s", symbol)
	}

	// 请求 SSE 数据
	resp, err := client.Get(optionUrl + "?" + payload.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Date  int     `json:"date"`
		Time  int     `json:"time"`
		Total int     `json:"total"`
		List  [][]any `json:"list"` // 动态类型
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var data []OptionFinanceBoardData
	timestamp := fmt.Sprintf("%d%06d", result.Date, result.Time)

	for _, item := range result.List {
		if len(item) < 5 {
			continue
		}
		price, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[1]), 64)
		chgRate, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[2]), 64)
		prevSettle, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[3]), 64)
		strike, _ := strconv.ParseFloat(fmt.Sprintf("%v", item[4]), 64)

		data = append(data, OptionFinanceBoardData{
			Date:        timestamp,
			ContractID:  fmt.Sprintf("%v", item[0]),
			Price:       price,
			ChangeRate:  chgRate,
			PrevSettle:  prevSettle,
			StrikePrice: strike,
			Quantity:    result.Total,
		})
	}
	return data, nil
}

// ==================== 2. 风险指标: option_risk_indicator_sse ====================
func OptionRiskIndicatorSSE(date string) ([]RiskIndicator, error) {
	const riskUrl = "http://query.sse.com.cn/commonQuery.do"

	params := url.Values{}
	params.Set("isPagination", "false")
	params.Set("trade_date", date)
	params.Set("sqlId", "SSE_ZQPZ_YSP_GGQQZSXT_YSHQ_QQFXZB_DATE_L")
	params.Set("contractSymbol", "")

	req, err := http.NewRequest("GET", riskUrl+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "query.sse.com.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "http://www.sse.com.cn/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result []map[string]string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var indicators []RiskIndicator
	for _, item := range result.Result {
		// 解析日期
		t, _ := api.ParseTime(item["TRADE_DATE"])

		// 转换浮点数
		delta, _ := strconv.ParseFloat(item["DELTA_VALUE"], 64)
		theta, _ := strconv.ParseFloat(item["THETA_VALUE"], 64)
		gamma, _ := strconv.ParseFloat(item["GAMMA_VALUE"], 64)
		vega, _ := strconv.ParseFloat(item["VEGA_VALUE"], 64)
		rho, _ := strconv.ParseFloat(item["RHO_VALUE"], 64)
		iv, _ := strconv.ParseFloat(item["IMPLC_VOLATLTY"], 64)

		indicators = append(indicators, RiskIndicator{
			TradeDate:       t,
			SecurityID:      item["SECURITY_ID"],
			ContractID:      item["CONTRACT_ID"],
			ContractSymbol:  item["CONTRACT_SYMBOL"],
			Delta:           delta,
			Theta:           theta,
			Gamma:           gamma,
			Vega:            vega,
			Rho:             rho,
			ImplcVolatility: iv,
		})
	}
	return indicators, nil
}

// ------------------------------- 1. 常量定义 -------------------------------
const (
	VIX_THRESHOLD_LOW        = 0.05
	VIX_THRESHOLD_HIGH       = 0.05
	HISTORICAL_QUANTILE_LOW  = 0.2
	HISTORICAL_QUANTILE_HIGH = 0.8
	RISK_FREE_RATE           = 0.02
)

// ------------------------------- 2. 数据结构定义 -------------------------------

// MergedOption
//
//	保持与您提供的 OptionFinanceBoardData, SZOptionData, RiskIndicator 定义一致
//	为了清晰，我们重新定义 MergedOption 结构
type MergedOption struct {
	ContractID      string
	Strike          float64
	Type            string
	Price           float64
	ExpireDate      time.Time
	TDays           int
	TYears          float64
	ImplcVolatility float64
	Delta           float64
}

// ------------------------------- 3. 计算“第四个星期三”函数 -------------------------------
func getFourthWednesday(year, month int) time.Time {
	// 1. 创建该月的第一天
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	// 2. 获取第一天是星期几 (0=周日, 1=周一, ..., 6=周六)
	weekdayOfFirst := int(firstDay.Weekday())

	// 3. 计算从1号到当月第一个星期三需要多少天
	// 目标星期三的编号是 2 (因为 周日=0, 周一=1, 周二=2, 周三=3, ...)
	// 这里有一个关键错误: 周三的编号是 3，不是 2！
	// 我们需要计算的是 (3 - weekdayOfFirst + 7) % 7
	daysUntilFirstWednesday := (3 - weekdayOfFirst + 7) % 7

	// 4. 第一个星期三的日期
	firstWednesday := 1 + daysUntilFirstWednesday

	// 5. 第四个星期三的日期 = 第一个星期三 + 21天
	fourthWednesdayDay := firstWednesday + 21

	return time.Date(year, time.Month(month), fourthWednesdayDay, 0, 0, 0, 0, time.Local)
}

// ------------------------------- 4. 提取并合并数据 (使用真实数据) -------------------------------

// 提取并合并数据
func extractAndMergeData(riskData []RiskIndicator, tradeDateStr string) ([]MergedOption, error) {
	// 1. 将 tradeDateStr 解析为 time.Time
	currentDate, err := api.ParseTime(tradeDateStr)
	if err != nil {
		return nil, fmt.Errorf("无效的交易日期: %s", tradeDateStr)
	}

	// 2. 从风险数据中提取所有 "华泰柏瑞沪深300ETF期权" 的价格数据
	// 我们需要先知道有哪些到期月份 (YYMM)，然后为每个月份调用 OptionFinanceBoard
	// 👇 修改 map 的 value 类型，从 float64 改为 OptionFinanceBoardData
	priceDataMap := make(map[string]map[string]OptionFinanceBoardData) // map[YYMM]map[ContractID]完整数据

	// 从风险数据中找出所有 300ETF 期权的到期月份
	seenYYMM := make(map[string]bool)
	for _, risk := range riskData {
		if strings.HasPrefix(risk.ContractID, "510300") && len(risk.ContractID) >= 13 {
			yyMM := risk.ContractID[7:11]
			seenYYMM[yyMM] = true
		}
	}

	// 为每个到期月份获取价格数据
	for yyMM := range seenYYMM {
		fmt.Printf("💰 正在获取 510300_%s 价格数据...\n", yyMM)
		priceData, err := OptionFinanceBoard("华泰柏瑞沪深300ETF期权", "20"+yyMM)
		if err != nil {
			log.Printf("⚠️ 获取 %s 价格数据失败: %v", yyMM, err)
			continue
		}

		// 👇 构建 map[ContractID]OptionFinanceBoardData
		priceMap := make(map[string]OptionFinanceBoardData)
		for _, price := range priceData {
			_, ok := priceMap[price.ContractID]
			if ok {
				continue
			}
			priceMap[price.ContractID] = price // 存储整个结构体
		}
		priceDataMap[yyMM] = priceMap
	}

	// 3. 提取并合并数据
	var merged []MergedOption
	for _, risk := range riskData {
		// 筛选 300ETF 期权
		if !strings.HasPrefix(risk.ContractID, "510300") {
			continue
		}

		// 期权合约ID
		contractID := risk.ContractID
		if len(contractID) < 13 {
			continue
		}

		// 提取类型和到期年月
		optType := string(contractID[6])
		yyMM := contractID[7:11]

		// 计算真实到期日
		year, _ := strconv.Atoi("20" + yyMM[:2])
		month, _ := strconv.Atoi(yyMM[2:4])
		expireDate := getFourthWednesday(year, month)

		// 计算剩余天数和年化时间
		tDays := int(expireDate.Sub(currentDate).Hours() / 24)
		tYears := float64(tDays) / 365.0

		// 从 priceDataMap 中获取价格
		priceMap, exists := priceDataMap[yyMM]
		if !exists {
			continue
		}
		price, exists := priceMap[contractID]
		if !exists {
			continue
		}

		// 过滤异常波动率
		if risk.ImplcVolatility <= 0.01 || risk.ImplcVolatility >= 1.0 {
			continue
		}
		//fmt.Printf("ContractID=%s, price=%f\n", contractID, price.Price)
		merged = append(merged, MergedOption{
			ContractID:      contractID,
			Strike:          price.StrikePrice,
			Type:            optType,
			Price:           price.Price,
			ExpireDate:      expireDate,
			TDays:           tDays,
			TYears:          tYears,
			ImplcVolatility: risk.ImplcVolatility,
			Delta:           risk.Delta,
		})
	}

	if len(merged) == 0 {
		return nil, fmt.Errorf("合并后数据为空")
	}

	fmt.Printf("✅ 提取 300ETF 期权: %d 条\n", len(merged))
	return merged, nil
}

// ------------------------------- 5. 计算“恐慌指数”(真实VIX) -------------------------------
func calculateRealVix(mergedData []MergedOption, tradeDateStr string, riskFreeRate float64) (float64, error) {
	currentDate, err := api.ParseTime(tradeDateStr)
	if err != nil {
		return 0, err
	}

	// 按到期日分组
	groups := make(map[time.Time][]MergedOption)
	for _, opt := range mergedData {
		groups[opt.ExpireDate] = append(groups[opt.ExpireDate], opt)
	}

	var expirations []time.Time
	for exp := range groups {
		expirations = append(expirations, exp)
	}
	sort.Slice(expirations, func(i, j int) bool {
		return expirations[i].Before(expirations[j])
	})

	if len(expirations) < 2 {
		return 0, fmt.Errorf("不足两个到期日")
	}

	// 找到 T1 < 30/365 < T2 的组合
	targetT := 30.0 / 365.0
	var t1, t2 time.Time
	var T1, T2 float64
	found := false

	for i := 0; i < len(expirations)-1; i++ {
		T1 = expirations[i].Sub(currentDate).Hours() / 24 / 365
		T2 = expirations[i+1].Sub(currentDate).Hours() / 24 / 365

		if T1 < targetT && targetT < T2 {
			t1, t2 = expirations[i], expirations[i+1]
			found = true
			break
		}
	}

	if !found {
		fmt.Println("⚠️ 无满足 T1<30<T2 的组合，使用最近两个")
		t1, t2 = expirations[0], expirations[1]
		T1 = t1.Sub(currentDate).Hours() / 24 / 365
		T2 = t2.Sub(currentDate).Hours() / 24 / 365
	}

	fmt.Printf("🎯 使用到期日: %s (%.1f天), %s (%.1f天)\n",
		t1.Format("2006-01-02"), T1*365, t2.Format("2006-01-02"), T2*365)

	term1 := groups[t1]
	term2 := groups[t2]
	fmt.Println("==>", len(term1), len(term2))
	fmt.Printf("T1=%v, T2=%v\n", T1, T2)

	var1, err := computeVariance(term1, T1, riskFreeRate)
	if err != nil {
		return 0, err
	}

	var2, err := computeVariance(term2, T2, riskFreeRate)
	if err != nil {
		return 0, err
	}

	if var1 <= 0 || var2 <= 0 {
		return 0, fmt.Errorf("方差非正")
	}
	fmt.Println(var1, var2)

	vixSquared := ((T2-targetT)*var1 + (targetT-T1)*var2) / (T2 - T1)
	vix := math.Sqrt(vixSquared) * 100
	fmt.Printf("📊 M1 方差: %.6f, M2 方差: %.6f\n", var1, var2)
	fmt.Printf("🎯 插值得到 30 天方差: %.6f → VIX = %.2f\n", vixSquared, vix)
	return math.Max(vix, 5.0), nil
}

func computeVariance(options []MergedOption, T, r float64) (float64, error) {
	if len(options) == 0 {
		return 0, fmt.Errorf("计算方差失败: 期权数据为空")
	}

	if T <= 1.0/365.0 { // 小于1天
		return 0, fmt.Errorf("T <= 1天")
	}

	discount := math.Exp(-r * T)

	// ✅ 1. 按行权价升序排序(必须！)
	sort.Slice(options, func(i, j int) bool {
		return options[i].Strike < options[j].Strike
	})
	// 在提取并合并数据后，排序时加入 TYPE 控制
	sort.Slice(options, func(i, j int) bool {
		if options[i].Strike == options[j].Strike {
			// 相同行权价时: Put 在 Call 前
			if options[i].Type == "P" && options[j].Type == "C" {
				return true
			}
			if options[i].Type == "C" && options[j].Type == "P" {
				return false
			}
			return false // 相同类型，顺序不变
		}
		return options[i].Strike < options[j].Strike
	})

	// ✅ 2. 过滤 Price > 0 的有效合约
	var validOptions []MergedOption
	for _, opt := range options {
		if opt.Price > 0 {
			validOptions = append(validOptions, opt)
		}
	}
	if len(validOptions) == 0 {
		return 0, fmt.Errorf("计算方差失败: 所有期权价格均为0")
	}
	options = validOptions

	// ✅ 3. 提取 Call 和 Put
	var calls, puts []MergedOption
	for _, opt := range options {
		if opt.Type == "C" {
			calls = append(calls, opt)
		} else if opt.Type == "P" {
			puts = append(puts, opt)
		}
	}

	if len(calls) == 0 || len(puts) == 0 {
		return 0, fmt.Errorf("缺少 Call 或 Put 合约")
	}

	// 构建 Put 行权价映射
	putMap := make(map[float64]float64)
	for _, put := range puts {
		putMap[put.Strike] = put.Price
	}

	// 对齐 C-P
	var cMinusP []float64
	var strikes []float64
	for _, call := range calls {
		if putPrice, exists := putMap[call.Strike]; exists {
			cMinusP = append(cMinusP, call.Price-putPrice)
			strikes = append(strikes, call.Strike)
		}
	}

	if len(cMinusP) == 0 {
		return 0, fmt.Errorf("没有找到有效的 Call-Put 对")
	}

	// ✅ 4. 插值找 F
	var F float64
	foundCross := false
	for i := 0; i < len(cMinusP)-1; i++ {
		if cMinusP[i]*cMinusP[i+1] <= 0 {
			k1, k2 := strikes[i], strikes[i+1]
			c1, c2 := cMinusP[i], cMinusP[i+1]
			if c2 != c1 {
				w := -c1 / (c2 - c1)
				F = k1 + w*(k2-k1)
			} else {
				F = (k1 + k2) / 2
			}
			foundCross = true
			break
		}
	}
	if !foundCross {
		minIdx := 0
		minAbs := math.Abs(cMinusP[0])
		for i, v := range cMinusP {
			if math.Abs(v) < minAbs {
				minAbs = math.Abs(v)
				minIdx = i
			}
		}
		F = strikes[minIdx]
	}

	// ✅ 5. 截断行权价范围(去噪声)
	var filtered []MergedOption
	for _, opt := range options {
		if opt.Strike >= 0.7*F && opt.Strike <= 1.3*F {
			filtered = append(filtered, opt)
		}
	}
	if len(filtered) < 2 {
		return 0, fmt.Errorf("截断后数据不足")
	}
	options = filtered

	// ✅ 6. 重新排序(确保)
	sort.Slice(options, func(i, j int) bool {
		return options[i].Strike < options[j].Strike
	})

	// ✅ 7. 找 K0: 小于等于 F 的最大行权价
	K0 := 0.0
	for _, opt := range options {
		if opt.Strike <= F {
			K0 = opt.Strike
		} else {
			break
		}
	}
	if K0 == 0.0 {
		K0 = options[0].Strike // 保底
	}

	// ✅ 8. 计算加权和(✅ 先贴现！)
	sum_ := 0.0
	fmt.Printf("\n🔍 开始计算 sum_ (T=%.4f, discount=%.6f)\n", T, discount)
	for i, opt := range options {
		Q := opt.Price * discount // ✅ 贴现价格
		if Q <= 0 {
			continue
		}
		K := opt.Strike
		var dk float64
		if i == 0 {
			dk = options[i+1].Strike - K
		} else if i == len(options)-1 {
			dk = K - options[i-1].Strike
		} else {
			dk = (options[i+1].Strike - options[i-1].Strike) / 2
		}
		weight := dk / (K * K)
		contrib := weight * Q
		fmt.Printf("  K=%5.3f | P=%6.4f | Q=%7.5f | dk=%6.4f | w=%8.6f | → %8.6f\n", K, opt.Price, Q, dk, weight, contrib)
		sum_ += contrib
	}

	// ✅ 9. 计算方差(✅ 不再乘 discount)
	variance := (2.0 / T) * sum_
	variance -= math.Pow((F/K0)-1, 2) / T
	fmt.Printf("✅ [T=%.4f] F=%.4f, K0=%.4f, sum=%.6f, variance=%.6f\n", T, F, K0, sum_, variance)
	return math.Max(variance, 1e-6), nil
}
