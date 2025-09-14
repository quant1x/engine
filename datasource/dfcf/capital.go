package dfcf

import (
	"encoding/json"
	urlpkg "net/url"
	"strings"
	"time"

	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/http"
)

type rawCapital struct {
	Xsjj []struct {
		SECUCODE              string  `json:"SECUCODE"`
		SECURITYCODE          string  `json:"SECURITY_CODE"`
		LIFTDATE              string  `json:"LIFT_DATE"`
		LIFTNUM               int     `json:"LIFT_NUM"`
		LIFTTYPE              string  `json:"LIFT_TYPE"`
		TOTALSHARESRATIO      float64 `json:"TOTAL_SHARES_RATIO"`
		UNLIMITEDASHARESRATIO float64 `json:"UNLIMITED_A_SHARES_RATIO"`
	} `json:"xsjj"` // 限售解禁
	Gbjg []struct {
		SECUCODE           string  `json:"SECUCODE"`
		SECURITYCODE       string  `json:"SECURITY_CODE"`
		NONFREESHARES      int     `json:"NON_FREE_SHARES"`
		LIMITEDSHARES      int     `json:"LIMITED_SHARES"`
		UNLIMITEDSHARES    int     `json:"UNLIMITED_SHARES"`
		TOTALSHARES        int     `json:"TOTAL_SHARES"`
		LISTEDASHARES      int     `json:"LISTED_A_SHARES"`
		BFREESHARE         int     `json:"B_FREE_SHARE"`
		HFREESHARE         int     `json:"H_FREE_SHARE"`
		OTHERFREESHARES    int     `json:"OTHER_FREE_SHARES"`
		NONFREESHARESRATIO int     `json:"NON_FREESHARES_RATIO"`
		LIMITEDSHARESRATIO float64 `json:"LIMITED_SHARES_RATIO"`
		LISTEDSHARESRATIO  float64 `json:"LISTED_SHARES_RATIO"`
		TOTALSHARESRATIO   string  `json:"TOTAL_SHARES_RATIO"`
		LISTEDARATIOPC     float64 `json:"LISTED_A_RATIOPC"`
		LISTEDBRATIOPC     float64 `json:"LISTED_B_RATIOPC"`
		LISTEDHRATIOPC     float64 `json:"LISTED_H_RATIOPC"`
		LISTEDOTHERRATIOPC float64 `json:"LISTED_OTHER_RATIOPC"`
		LISTEDSUMRATIOPC   int     `json:"LISTED_SUM_RATIOPC"`
	} `json:"gbjg"` // 股本结构
	Lngbbd []struct {
		SECUCODE               string      `json:"SECUCODE"`
		SECURITYCODE           string      `json:"SECURITY_CODE"`
		ENDDATE                string      `json:"END_DATE"`
		TOTALSHARES            int         `json:"TOTAL_SHARES"`
		LIMITEDSHARES          int         `json:"LIMITED_SHARES"`
		LIMITEDOTHARS          int         `json:"LIMITED_OTHARS"`
		LIMITEDDOMESTICNATURAL int         `json:"LIMITED_DOMESTIC_NATURAL"`
		LIMITEDSTATELEGAL      int         `json:"LIMITED_STATE_LEGAL"`
		LIMITEDOVERSEASNOSTATE interface{} `json:"LIMITED_OVERSEAS_NOSTATE"`
		LIMITEDOVERSEASNATURAL interface{} `json:"LIMITED_OVERSEAS_NATURAL"`
		UNLIMITEDSHARES        int         `json:"UNLIMITED_SHARES"`
		LISTEDASHARES          int         `json:"LISTED_A_SHARES"`
		BFREESHARE             interface{} `json:"B_FREE_SHARE"`
		HFREESHARE             interface{} `json:"H_FREE_SHARE"`
		FREESHARES             int         `json:"FREE_SHARES"`
		LIMITEDASHARES         int         `json:"LIMITED_A_SHARES"`
		NONFREESHARES          interface{} `json:"NON_FREE_SHARES"`
		LIMITEDBSHARES         interface{} `json:"LIMITED_B_SHARES"`
		OTHERFREESHARES        interface{} `json:"OTHER_FREE_SHARES"`
		LIMITEDSTATESHARES     interface{} `json:"LIMITED_STATE_SHARES"`
		LIMITEDDOMESTICNOSTATE int         `json:"LIMITED_DOMESTIC_NOSTATE"`
		LOCKSHARES             int         `json:"LOCK_SHARES"`
		LIMITEDFOREIGNSHARES   interface{} `json:"LIMITED_FOREIGN_SHARES"`
		LIMITEDHSHARES         interface{} `json:"LIMITED_H_SHARES"`
		SPONSORSHARES          interface{} `json:"SPONSOR_SHARES"`
		STATESPONSORSHARES     interface{} `json:"STATE_SPONSOR_SHARES"`
		SPONSORSOCIALSHARES    interface{} `json:"SPONSOR_SOCIAL_SHARES"`
		RAISESHARES            interface{} `json:"RAISE_SHARES"`
		RAISESTATESHARES       interface{} `json:"RAISE_STATE_SHARES"`
		RAISEDOMESTICSHARES    interface{} `json:"RAISE_DOMESTIC_SHARES"`
		RAISEOVERSEASSHARES    interface{} `json:"RAISE_OVERSEAS_SHARES"`
		CHANGEREASON           string      `json:"CHANGE_REASON"`
	} `json:"lngbbd"` // 历年股本变动
	Gbgc []struct {
		SECUCODE      string `json:"SECUCODE"`
		SECURITYCODE  string `json:"SECURITY_CODE"`
		ENDDATE       string `json:"END_DATE"`
		TOTALSHARES   int    `json:"TOTAL_SHARES"`
		LISTEDASHARES int    `json:"LISTED_A_SHARES"`
		LIMITEDSHARES int    `json:"LIMITED_SHARES"`
	} `json:"gbgc"` // 股本构成
}

const (
	urlCapitalStockStructure = "https://emweb.securities.eastmoney.com/PC_HSF10/CapitalStockStructure/PageAjax"
)

type StockCapital struct {
	Code            string // 证券代码
	Date            string // 变动日期
	TotalShares     int    // 总股本
	UnlimitedShares int    // 已流通股本
	ListedAShares   int    // 已上市流通A股
	ChangeReson     string // 变动原因
	UpdateTime      string // 更新时间
}

// CapitalChange 获取股本变动记录
//
//	deprecated: 不推荐, 太慢
func CapitalChange(securityCode string) (list []StockCapital) {
	code := exchange.CorrectSecurityCode(securityCode)
	params := urlpkg.Values{
		"code": {strings.ToUpper(code)},
	}

	url := urlCapitalStockStructure + "?" + params.Encode()
	data, lastModified, err := http.Request(url, http.MethodGet, "")
	//fmt.Println(api.Bytes2String(data))
	if err != nil {
		return
	}
	var css rawCapital
	err = json.Unmarshal(data, &css)
	if err != nil {
		return
	}
	if lastModified.UnixMilli() > 0 {
		lastModified = time.Now()
	}
	updateTime := lastModified.Format(time.DateTime)
	for _, v := range css.Lngbbd {
		sc := StockCapital{
			Code:            securityCode,
			Date:            exchange.FixTradeDate(v.ENDDATE),
			TotalShares:     v.TOTALSHARES,
			UnlimitedShares: v.UNLIMITEDSHARES,
			ListedAShares:   v.LISTEDASHARES,
			ChangeReson:     v.CHANGEREASON,
			UpdateTime:      updateTime,
		}
		list = append(list, sc)
	}
	api.SliceSort(list, func(a, b StockCapital) bool {
		return a.Date > b.Date
	})
	return
}
