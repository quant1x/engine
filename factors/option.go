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

	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/num"
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

// ==================== 1. 期权行情数据：option_finance_board ====================

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

// ==================== 2. 风险指标：option_risk_indicator_sse ====================
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
	// 这里有一个关键错误：周三的编号是 3，不是 2！
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
		fmt.Printf("ContractID=%s, price=%f\n", contractID, price.Price)
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

// ------------------------------- 5. 计算“恐慌指数”（真实VIX） -------------------------------
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
	fmt.Println("==>", T1, T2)

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

	return math.Max(vix, 5.0), nil
}

func computeVariance(options []MergedOption, T, r float64) (float64, error) {
	if len(options) == 0 {
		return 0, fmt.Errorf("计算方差失败：期权数据为空")
	}

	if T <= 0 {
		return 0, fmt.Errorf("T <= 0")
	}

	discount := math.Exp(-r * T)
	sort.Slice(options, func(i, j int) bool {
		return options[i].Strike < options[j].Strike
	})

	// 👉 1. 创建新的切片，只包含 Price > 0 的合约
	var validOptions []MergedOption
	for _, opt := range options {
		if opt.Price > 0 {
			validOptions = append(validOptions, opt)
		}
	}

	if len(validOptions) == 0 {
		return 0, fmt.Errorf("计算方差失败：所有期权价格均为0")
	}

	// 👉 2. 使用过滤后的 validOptions 进行后续计算
	var calls, puts []MergedOption
	for _, opt := range validOptions {
		if opt.Type == "C" {
			calls = append(calls, opt)
		} else if opt.Type == "P" {
			puts = append(puts, opt)
		}
	}

	if len(calls) == 0 || len(puts) == 0 {
		return 0, fmt.Errorf("计算方差失败：缺少 Call 或 Put 合约")
	}

	putMap := make(map[float64]float64)
	for _, put := range puts {
		putMap[put.Strike] = put.Price
	}
	fmt.Printf("Debug: Total options: %d, Calls: %d, Puts: %d\n", len(options), len(calls), len(puts))
	fmt.Println("Debug: Call-Put Pairs:")
	for _, call := range calls {
		putPrice, exists := putMap[call.Strike]
		if exists {
			fmt.Printf("  Strike: %.3f, C: %.4f, P: %.4f, C-P: %.4f\n",
				call.Strike, call.Price, putPrice, call.Price-putPrice)
		}
	}

	var cMinusP []float64
	var strikes []float64
	for _, call := range calls {
		putPrice, exists := putMap[call.Strike]
		if !exists {
			continue
		}
		cMinusP = append(cMinusP, call.Price-putPrice)
		strikes = append(strikes, call.Strike)
	}

	if len(cMinusP) == 0 {
		return 0, fmt.Errorf("计算方差失败：没有找到有效的 Call-Put 对")
	}

	// 插值找到 C-P=0 的点 (F)
	var F float64
	found := false
	for i := 0; i < len(cMinusP)-1; i++ {
		if cMinusP[i]*cMinusP[i+1] <= 0 {
			// 找到交叉点，进行线性插值
			k1, k2 := strikes[i], strikes[i+1]
			c1, c2 := cMinusP[i], cMinusP[i+1]
			if c2 != c1 {
				w := -c1 / (c2 - c1)
				F = k1 + w*(k2-k1)
			} else {
				F = (k1 + k2) / 2
			}
			found = true
			break
		}
	}
	if !found {
		// 如果没有交叉点，取绝对值最小的
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

	fmt.Println("远期价格 F ≈ ", F)

	// 👉 4. 找到最接近 F 的行权价 K0
	var K0 float64
	//minDiff := math.Abs(options[0].Strike - F)
	//K0 = options[0].Strike
	//for _, opt := range options {
	//	diff := math.Abs(opt.Strike - F)
	//	if diff < minDiff {
	//		minDiff = diff
	//		K0 = opt.Strike
	//	}
	//}
	for _, opt := range options {
		if F >= opt.Strike {
			K0 = opt.Strike
		} else {
			break
		}
	}

	// 👉 5. 计算主项的加权和
	var sum_ float64
	for i, opt := range options {
		var K float64
		var dk float64
		Q := opt.Price
		if num.IsNaN(Q) || Q <= 0 {
			continue
		}
		K = opt.Strike
		if i == 0 {
			dk = options[i+1].Strike - opt.Strike
		} else if i == len(options)-1 {
			dk = opt.Strike - options[i-1].Strike
		} else {
			dk = (options[i+1].Strike - options[i-1].Strike) / 2
		}
		fmt.Printf("%d: dk=%f, K=%f, Q=%f\n", i, dk, K, Q)
		weight := dk / (K * K)
		sum_ += weight * Q
		fmt.Printf("sum_: %f\n", sum_)
	}
	fmt.Println("        T =", T)
	fmt.Println("      sum =", sum_)
	fmt.Println("        F =", F)
	fmt.Println("       K0 =", K0)
	fmt.Println(" discount =", discount)
	// 👉 6. 计算完整的方差 (包含修正项)
	variance := (2.0 / T) * sum_
	variance -= math.Pow((F/K0)-1, 2) / T
	variance *= discount

	return variance, nil
}
