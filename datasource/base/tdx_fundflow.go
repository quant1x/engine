package base

import (
	"strconv"
	"strings"

	"github.com/quant1x/exchange"
	"github.com/quant1x/gotdx"
	"github.com/quant1x/num"
	"github.com/quant1x/pandas"
	"github.com/quant1x/x/api"
)

const (
	kCompanyInfoFundFlow = "资金动向"
	kCategoryFundFlow    = "资金流向"
)

//【5.资金流向】
//日期        主力净额  主力净额 超大单净买入 超大单净买入 大单净买入 大单净买入  主买净额  主买净额
//            金额(元)   占比(%)   金额(元)      占比(%)    金额(元)    占比(%)   金额(元)   占比(%)
//─────────────────────────────────────────────────
//2023-05-15 -8581.10万    -4.95   -3069.63万        -1.77 -5511.47万      -3.18    -1.46亿    -8.44
//2023-05-12    -3.26亿   -16.20      -3.01亿       -14.95 -2525.27万      -1.26    -1.62亿    -8.05
//2023-05-11     2.75亿    17.10       2.95亿        18.35 -2007.90万      -1.25  -878.46万    -0.55
//2023-05-10    -3.27亿   -14.45      -2.35亿       -10.37 -9216.30万      -4.07    -2.43亿   -10.75
//2023-05-09     1.96亿    25.91       1.97亿        26.02   -87.68万      -0.12    -3.76亿   -49.78
//2023-05-08 -1161.92万    -1.30   -1778.75万        -1.98   616.83万       0.69 -3290.33万    -3.67
//2023-05-05  -515.71万    -0.56   -1642.56万        -1.79  1126.85万       1.23 -3437.44万    -3.75
//2023-05-04  7685.06万     6.09    3398.48万         2.69  4286.57万       3.40     1.09亿     8.67
//2023-04-28  2158.92万     3.00    1249.07万         1.74   909.85万       1.27  4346.58万     6.05
//2023-04-27    -1.30亿   -17.61   -7470.81万       -10.13 -5517.81万      -7.48 -2907.44万    -3.94
//2023-04-26  2127.77万     2.14    1367.42万         1.37   760.34万       0.76     1.13亿    11.34
//2023-04-25   664.05万     1.08   -1458.93万        -2.36  2122.98万       3.44  -139.47万    -0.23
//2023-04-24  -126.71万    -0.21     295.24万         0.48  -421.95万      -0.69  -791.00万    -1.29
//2023-04-21 -8501.18万   -12.44   -3843.78万        -5.62 -4657.39万      -6.81 -8237.27万   -12.05
//2023-04-20  8625.04万     8.14    3603.48万         3.40  5021.56万       4.74     1.23亿    11.57
//2023-04-19  9310.97万     8.27       1.11亿         9.82 -1752.71万      -1.56     1.09亿     9.71
//2023-04-18  1312.51万     4.22    1552.84万         4.99  -240.34万      -0.77 -1989.31万    -6.39
//2023-04-17 -3970.27万   -10.76   -1600.06万        -4.34 -2370.21万      -6.42 -1318.26万    -3.57
//2023-04-14 -1514.15万    -3.66    -405.87万        -0.98 -1108.28万      -2.68 -1296.30万    -3.13
//2023-04-13    -1.09亿   -14.38   -5319.72万        -7.03 -5555.62万      -7.35 -5555.62万    -7.35
//─────────────────────────────────────────────────

var (
	TdxFieldsFundFlow = []string{"日期", "主力净额金额(元)", "主力净额占比(%)", "超大单净买入金额(元)", "超大单净买入占比(%)", "大单净买入金额(元)", "大单净买入占比(%)", "主买净额金额(元)", "主买净额占比(%)"}
)

func splitContent(content, unit string) (headers []string, lines [][]string) {
	//c := strings.ReplaceAll(content, "-\\u003e", "->")
	//arr := strings.Split(c, "\\r\\n\\r\\n")
	arr := strings.Split(content, "\r\n\r\n")
	for i, block := range arr {
		block = strings.TrimSpace(block)
		if i > 0 && strings.Index(block, unit) >= 0 {
			arr := strings.Split(block, "\r\n")
			tmpHeaders := []string{}
			numberFound := false
			for _, v := range arr {
				if strings.Index(v, unit) >= 0 {
					continue
				}
				if strings.HasPrefix(v, "──") {
					continue
				}
				v = strings.TrimSpace(v)
				// 非数字开头是表头, 数字开头为数据
				ch := v[0]
				if ch >= '0' && ch <= '9' {
					if !numberFound {
						numberFound = true
					}
					cols := []string{}
					foundDate := false
					for _, tmp := range strings.Fields(v) {
						tf := float64(0)
						tmp = strings.TrimSpace(tmp)
						//if tmp == "2023-04-17" {
						//	fmt.Println("found")
						//}
						if !foundDate {
							if _, err := api.ParseTime(tmp); err == nil {
								cols = append(cols, tmp)
								foundDate = true
								continue
							}
						}
						if fs, _, ok := strings.Cut(tmp, "万"); ok {
							tf = num.AnyToFloat64(fs) * 10000
						} else if fs, _, ok := strings.Cut(tmp, "亿"); ok {
							tf = num.AnyToFloat64(fs) * 100000000
						} else {
							tf = num.AnyToFloat64(fs)
						}
						f := strconv.FormatFloat(tf, 'f', -1, 64)
						cols = append(cols, f)
					}
					lines = append(lines, cols)

				} else {
					tmpHeaders = append(tmpHeaders, v)
				}

			}
			if numberFound {
				tmpHeadersCount := len(tmpHeaders)
				if tmpHeadersCount >= 1 {
					headers = strings.Fields(tmpHeaders[0])
				}
				numberOfHeaderFields := len(headers)
				for j := 1; j < tmpHeadersCount; j++ {
					remaining := strings.Fields(tmpHeaders[j])
					numberOfRemainingFields := len(remaining)
					for k := 0; k < numberOfRemainingFields; k++ {
						pos := 1 + k
						headers[numberOfHeaderFields-pos] += remaining[numberOfRemainingFields-pos]
					}

				}
			}
			break
		}
	}
	return
}

// FundFlow 资金流向
//
//	deprecated: 不推荐
func FundFlow(securityCode string) pandas.DataFrame {
	tdxApi := gotdx.GetTdxApi()
	securityCode = exchange.CorrectSecurityCode(securityCode)
	reply, err := tdxApi.GetCompanyInfoContent(securityCode, kCompanyInfoFundFlow)
	if err != nil {
		return pandas.DataFrame{Err: err}
	}
	//fmt.Println("code:", securityCode)
	//fmt.Printf("%+v\n", reply)
	//data, _ := json.Marshal(reply)
	//text := api.Bytes2String(data)
	//fmt.Println(text)
	//dict := reply.Map("资金流向")
	//dict.Each(func(key interface{}, value interface{}) {
	//	fmt.Println(key, value)
	//})

	headers, lines := splitContent(reply.Content, kCategoryFundFlow)
	//for _, v := range headers {
	//	fmt.Printf("%s|", v)
	//}
	//fmt.Println()
	//for _, vv := range lines {
	//	for _, v := range vv {
	//		fmt.Printf("%s|", v)
	//	}
	//	fmt.Println()
	//}
	//fmt.Println()
	if len(headers) == 0 || len(lines) == 0 || len(headers) != len(lines[0]) {
		return pandas.DataFrame{}
	}
	rows := [][]string{}
	rows = append(rows, headers)
	fieldsNum := len(headers)
	for _, v := range lines {
		if len(v) != fieldsNum {
			continue
		}
		rows = append(rows, v)
	}

	df := pandas.LoadRecords(rows)
	return df
}
