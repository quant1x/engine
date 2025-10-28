package dfcf

import (
	"encoding/json"
	"fmt"
	urlpkg "net/url"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/gox/http"
)

//龙虎榜分为两个接口
//https://datacenter-web.eastmoney.com/api/data/v1/get?reportName=RPT_BILLBOARD_DAILYDETAILSBUY&columns=ALL&filter=(TRADE_DATE%3D%272024-07-16%27)(SECURITY_CODE%3D%22600611%22)&pageNumber=1&pageSize=50&sortTypes=-1&sortColumns=BUY&source=WEB&client=WEB&_=1721172995452
//https://datacenter-web.eastmoney.com/api/data/v1/get?reportName=RPT_BILLBOARD_DAILYDETAILSSELL&columns=ALL&filter=(TRADE_DATE%3D%272024-07-16%27)(SECURITY_CODE%3D%22600611%22)&pageNumber=1&pageSize=50&sortTypes=-1&sortColumns=SELL&source=WEB&client=WEB&_=1721172995453
//{
//    "version": "60e1d52ca657ea6eb5fc5076817e6020",
//    "result": {
//        "pages": 1,
//        "data": [
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10023614",
//                "OPERATEDEPT_NAME": "国泰君安证券股份有限公司总部",
//                "EXPLANATION": "有价格涨跌幅限制的日收盘价格涨幅偏离值达到7%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 37521554.1,
//                "SELL": null,
//                "NET": 37521554.1,
//                "RISE_PROBABILITY_3DAY": 39.393939393939,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 231,
//                "CHANGE_TYPE": "137001002001001",
//                "OPERATEDEPT_CODE_OLD": "80033895",
//                "TOTAL_BUYRIO": 0.019000938125,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016460"
//            },
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10023614",
//                "OPERATEDEPT_NAME": "国泰君安证券股份有限公司总部",
//                "EXPLANATION": "有价格涨跌幅限制的日换手率达到20%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 37521554.1,
//                "SELL": null,
//                "NET": 37521554.1,
//                "RISE_PROBABILITY_3DAY": 39.393939393939,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 231,
//                "CHANGE_TYPE": "137001004001",
//                "OPERATEDEPT_CODE_OLD": "80033895",
//                "TOTAL_BUYRIO": 0.019000938125,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016469"
//            },
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10023739",
//                "OPERATEDEPT_NAME": "东方证券股份有限公司上海浦东新区源深路证券营业部",
//                "EXPLANATION": "有价格涨跌幅限制的日收盘价格涨幅偏离值达到7%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 37045933,
//                "SELL": null,
//                "NET": 37045933,
//                "RISE_PROBABILITY_3DAY": 66.666666666667,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 15,
//                "CHANGE_TYPE": "137001002001001",
//                "OPERATEDEPT_CODE_OLD": "80034027",
//                "TOTAL_BUYRIO": 0.018760083307,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016460"
//            },
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10023739",
//                "OPERATEDEPT_NAME": "东方证券股份有限公司上海浦东新区源深路证券营业部",
//                "EXPLANATION": "有价格涨跌幅限制的日换手率达到20%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 37045933,
//                "SELL": null,
//                "NET": 37045933,
//                "RISE_PROBABILITY_3DAY": 66.666666666667,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 15,
//                "CHANGE_TYPE": "137001004001",
//                "OPERATEDEPT_CODE_OLD": "80034027",
//                "TOTAL_BUYRIO": 0.018760083307,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016469"
//            },
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10456710",
//                "OPERATEDEPT_NAME": "国盛证券有限责任公司宁波桑田路证券营业部",
//                "EXPLANATION": "有价格涨跌幅限制的日收盘价格涨幅偏离值达到7%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 32388929.49,
//                "SELL": null,
//                "NET": 32388929.49,
//                "RISE_PROBABILITY_3DAY": 37.5,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 128,
//                "CHANGE_TYPE": "137001002001001",
//                "OPERATEDEPT_CODE_OLD": "80425702",
//                "TOTAL_BUYRIO": 0.016401773859,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016460"
//            },
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10456710",
//                "OPERATEDEPT_NAME": "国盛证券有限责任公司宁波桑田路证券营业部",
//                "EXPLANATION": "有价格涨跌幅限制的日换手率达到20%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 32388929.49,
//                "SELL": null,
//                "NET": 32388929.49,
//                "RISE_PROBABILITY_3DAY": 37.5,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 128,
//                "CHANGE_TYPE": "137001004001",
//                "OPERATEDEPT_CODE_OLD": "80425702",
//                "TOTAL_BUYRIO": 0.016401773859,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016469"
//            },
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10453418",
//                "OPERATEDEPT_NAME": "中信证券股份有限公司西安朱雀大街证券营业部",
//                "EXPLANATION": "有价格涨跌幅限制的日收盘价格涨幅偏离值达到7%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 27223118,
//                "SELL": null,
//                "NET": 27223118,
//                "RISE_PROBABILITY_3DAY": 35.555555555556,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 90,
//                "CHANGE_TYPE": "137001002001001",
//                "OPERATEDEPT_CODE_OLD": "80422556",
//                "TOTAL_BUYRIO": 0.013785803736,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016460"
//            },
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10453418",
//                "OPERATEDEPT_NAME": "中信证券股份有限公司西安朱雀大街证券营业部",
//                "EXPLANATION": "有价格涨跌幅限制的日换手率达到20%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 27223118,
//                "SELL": null,
//                "NET": 27223118,
//                "RISE_PROBABILITY_3DAY": 35.555555555556,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 90,
//                "CHANGE_TYPE": "137001004001",
//                "OPERATEDEPT_CODE_OLD": "80422556",
//                "TOTAL_BUYRIO": 0.013785803736,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016469"
//            },
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10025200",
//                "OPERATEDEPT_NAME": "国泰君安证券股份有限公司昆明人民中路证券营业部",
//                "EXPLANATION": "有价格涨跌幅限制的日收盘价格涨幅偏离值达到7%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 25337327,
//                "SELL": null,
//                "NET": 25337327,
//                "RISE_PROBABILITY_3DAY": 37.931034482759,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 29,
//                "CHANGE_TYPE": "137001002001001",
//                "OPERATEDEPT_CODE_OLD": "80034252",
//                "TOTAL_BUYRIO": 0.012830838011,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016460"
//            },
//            {
//                "SECURITY_CODE": "600611",
//                "SECUCODE": "600611.SH",
//                "TRADE_DATE": "2024-07-16 00:00:00",
//                "OPERATEDEPT_CODE": "10025200",
//                "OPERATEDEPT_NAME": "国泰君安证券股份有限公司昆明人民中路证券营业部",
//                "EXPLANATION": "有价格涨跌幅限制的日换手率达到20%的前五只证券",
//                "CHANGE_RATE": 10.1124,
//                "CLOSE_PRICE": 4.9,
//                "ACCUM_AMOUNT": 1974721135,
//                "ACCUM_VOLUME": 423841126,
//                "BUY": 25337327,
//                "SELL": null,
//                "NET": 25337327,
//                "RISE_PROBABILITY_3DAY": 37.931034482759,
//                "TOTAL_BUYER_SALESTIMES_3DAY": 29,
//                "CHANGE_TYPE": "137001004001",
//                "OPERATEDEPT_CODE_OLD": "80034252",
//                "TOTAL_BUYRIO": 0.012830838011,
//                "TOTAL_SELLRIO": null,
//                "TRADE_ID": "5016469"
//            }
//        ],
//        "count": 10
//    },
//    "success": true,
//    "message": "ok",
//    "code": 0
//}

const (
	urlLongHuBang = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	lhbBuy        = "BUY"  // 龙虎榜买入
	lhbSell       = "SELL" // 龙虎榜卖出
)

var (
	mapTypeLonghuBang = map[string]string{
		lhbBuy:  "RPT_BILLBOARD_DAILYDETAILSBUY",
		lhbSell: "RPT_BILLBOARD_DAILYDETAILSSELL",
	}
)

// BillBoard 龙虎榜数据结构
type BillBoard struct {
	SECURITY_CODE               string  `json:"SECURITY_CODE"`               // 证券代码
	SECUCODE                    string  `json:"SECUCODE"`                    // 证券代码
	TRADE_DATE                  string  `json:"TRADE_DATE"`                  // 交易日期
	OPERATEDEPT_CODE            string  `json:"OPERATEDEPT_CODE"`            // 交易营业部代码
	OPERATEDEPT_NAME            string  `json:"OPERATEDEPT_NAME"`            // 交易营业部名称
	EXPLANATION                 string  `json:"EXPLANATION"`                 // 龙虎榜类型, 说明
	CHANGE_RATE                 float64 `json:"CHANGE_RATE"`                 // 涨跌幅
	CLOSE_PRICE                 float64 `json:"CLOSE_PRICE"`                 // 收盘价
	ACCUM_AMOUNT                int     `json:"ACCUM_AMOUNT"`                // 成交金额
	ACCUM_VOLUME                int     `json:"ACCUM_VOLUME"`                // 成交量
	BUY                         float64 `json:"BUY"`                         // 买入
	SELL                        float64 `json:"SELL"`                        // 卖出
	NET                         float64 `json:"NET"`                         // 净额
	RISE_PROBABILITY_3DAY       float64 `json:"RISE_PROBABILITY_3DAY"`       // 上榜3日
	TOTAL_BUYER_SALESTIMES_3DAY int     `json:"TOTAL_BUYER_SALESTIMES_3DAY"` // 买方3日
	CHANGE_TYPE                 string  `json:"CHANGE_TYPE"`                 // 变动类型
	OPERATEDEPT_CODE_OLD        string  `json:"OPERATEDEPT_CODE_OLD"`        // 交易营业部旧代码
	TOTAL_BUYRIO                float64 `json:"TOTAL_BUYRIO"`                // 买入占比
	TOTAL_SELLRIO               float64 `json:"TOTAL_SELLRIO"`               // 卖出占比
	TRADE_ID                    string  `json:"TRADE_ID"`                    // 交易ID
}

// 原始的龙虎榜列表
func rawBillBoardList(date string, pageNumber int, direction string) ([]BillBoard, int, error) {
	tradeDate := exchange.FixTradeDate(date)
	params := urlpkg.Values{
		"reportName":  {mapTypeLonghuBang[direction]},
		"columns":     {"ALL"},
		"source":      {"WEB"},
		"client":      {"WEB"},
		"sortColumns": {direction},
		"sortTypes":   {"-1"},
		"pageSize":    {fmt.Sprintf("%d", rzrqPageSize)},
		"pageNumber":  {fmt.Sprintf("%d", pageNumber)},
		"filter":      {fmt.Sprintf(`(TRADE_DATE='%s')`, tradeDate)},
	}

	url := urlLongHuBang + "?" + params.Encode()
	data, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	var raw rawResult[BillBoard]
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return nil, 0, err
	}
	return raw.Data, raw.Pages, nil
}
