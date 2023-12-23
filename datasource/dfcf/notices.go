package dfcf

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/exception"
	"gitee.com/quant1x/gox/http"
	"gitee.com/quant1x/pandas/stat"
	"math"
	urlpkg "net/url"
	"strings"
)

const (
	CacheL5KeyNotices        = "cache/notices"
	kUrlEastmoneyNotices     = "https://np-anotice-stock.eastmoney.com/api/security/ann"
	EastmoneyNoticesPageSize = 100

	kErrorBaseNotice = 91000
)

var (
	ErrNoticeBadApi   = exception.New(kErrorBaseNotice, "接口异常")
	ErrNoticeNotFound = exception.New(kErrorBaseNotice+1, "没有数据")
)

var (
	// 风险检测的关键词
	riskKeywords = []string{"处罚", "冻结", "诉讼", "质押", "仲裁", "持股5%以上股东权益变动", "信用减值", "商誉减值", "重大风险", "退市风险"}
)

type EMNoticeType = int

const (
	NoticeAll          EMNoticeType = iota // 全部
	NoticeUnused1                          // 财务报告
	NoticeUnused2                          // 融资公告
	NoticeUnused3                          // 风险提示
	NoticeUnused4                          // 信息变更
	NoticeWarning                          // 重大事项
	NoticeUnused6                          // 资产重组
	NoticeHolderChange                     // 持股变动
)

func GetNoticeType(noticeType EMNoticeType) string {
	switch noticeType {
	case NoticeAll:
		return "全部"
	case NoticeUnused1:
		return "财务报告"
	case NoticeUnused2:
		return "融资公告"
	case NoticeUnused3:
		return "风险提示"
	case NoticeUnused4:
		return "信息变更"
	case NoticeWarning:
		return "重大事项"
	case NoticeUnused6:
		return "资产重组"
	case NoticeHolderChange:
		return "持股变动"
	default:
		return "其它"
	}
}

// 公告原始的消息结构
type rawNoticePackage struct {
	Data struct {
		List []struct {
			ArtCode string `json:"art_code"`
			Codes   []struct {
				AnnType    string `json:"ann_type"`
				InnerCode  string `json:"inner_code"`
				MarketCode string `json:"market_code"`
				ShortName  string `json:"short_name"`
				StockCode  string `json:"stock_code"`
			} `json:"codes"`
			Columns []struct {
				ColumnCode string `json:"column_code"`
				ColumnName string `json:"column_name"`
			} `json:"columns"`
			DisplayTime string `json:"display_time"`
			EiTime      string `json:"eiTime"`
			Language    string `json:"language"`
			NoticeDate  string `json:"notice_date"`
			ProductCode string `json:"product_code"`
			SortDate    string `json:"sort_date"`
			SourceType  string `json:"source_type"`
			Title       string `json:"title"`
			TitleCh     string `json:"title_ch"`
			TitleEn     string `json:"title_en"`
		} `json:"list"`
		PageIndex int `json:"page_index"`
		PageSize  int `json:"page_size"`
		TotalHits int `json:"total_hits"`
	} `json:"data"`
	Error   string `json:"error"`
	Success int    `json:"success"`
}

// NoticeDetail 公告详情
type NoticeDetail struct {
	Code         string `csv:"证券代码" dataframe:"证券代码"`   // 证券代码
	Name         string `csv:"证券名称" dataframe:"证券名称"`   // 证券名称
	DisplayTime  string `csv:"显示时间" dataframe:"显示时间"`   // 显示时间
	NoticeDate   string `csv:"公告时间" dataframe:"公告时间"`   // 公告时间
	Title        string `csv:"内容提要" dataframe:"公告标题"`   // 公告标题
	Keywords     string `csv:"关键词" dataframe:"关键词"`     // 公告关键词
	Increase     int    `csv:"增持" dataframe:"增持"`       // 增持
	Reduce       int    `csv:"减持" dataframe:"减持"`       // 减持
	HolderChange int    `csv:"控制人变更" dataframe:"控制人变更"` // 实际控制人变更
	Risk         int    `csv:"风险数" dataframe:"监管"`      // 风险数
}

// AllNotices 东方财富网-数据中心-公告大全-沪深京 A 股公告
//
//	http://data.eastmoney.com/notices/hsa/5.html
//	:param symbol: 报告类型; choice of {"全部", "重大事项", "财务报告", "融资公告", "风险提示", "资产重组", "信息变更", "持股变动"}
//	:type symbol: str
//	:param date: 制定日期
//	:type date: str
//	:return: 沪深京 A 股公告
//	Deprecated: 弃用
func AllNotices(noticeType EMNoticeType, date string, pageNumber ...int) (notices []NoticeDetail, pages int, err error) {
	pageNo := 1
	if len(pageNumber) > 0 {
		pageNo = pageNumber[0]
	}
	beginDate := trading.FixTradeDate(date)
	endDate := trading.Today()
	pageSize := EastmoneyNoticesPageSize
	params := urlpkg.Values{
		"sr":         {"-1"},
		"page_size":  {fmt.Sprintf("%d", pageSize)},
		"page_index": {fmt.Sprintf("%d", pageNo)},
		"ann_type":   {"SHA,CYB,SZA,BJA"},
		//"ann_type":      {"A"},
		//"ann_type":      {"SHA,SZA"},
		"client_source": {"web"},
		"f_node":        {fmt.Sprintf("%d", noticeType)},
		"s_node":        {"0"},
		"begin_time":    {beginDate},
		"end_time":      {endDate},
		//"cb": {"jQuery112305241416374967685_1683838825141"},
	}
	// Host: np-anotice-stock.eastmoney.com
	header := map[string]any{
		//"User-Agent": config.HTTP_REQUEST_HEADER_USER_AGENT,
		//"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	}
	url := kUrlEastmoneyNotices + "?" + params.Encode()
	//url = "https://np-anotice-stock.eastmoney.com/api/security/ann?cb=jQuery112305241416374967685_1683838825141&sr=-1&page_size=50&page_index=1&ann_type=SHA%2CCYB%2CSZA%2CBJA&client_source=web&f_node=0&s_node=0"
	data, _, err := http.Request(url, http.MethodGet, "", header)
	if err != nil {
		return
	}
	//fmt.Println(api.Bytes2String(data))
	var raw rawNoticePackage
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return
	}
	if raw.Success != 1 || len(raw.Data.List) == 0 {
		err = ErrNoticeNotFound
		return
	}
	pages = int(math.Ceil(float64(raw.Data.TotalHits) / float64(EastmoneyNoticesPageSize)))

	for _, v := range raw.Data.List {
		marketCode := proto.MarketIdShenZhen
		if len(v.Codes) == 0 || len(v.Columns) == 0 {
			continue
		}
		code := v.Codes[0]
		mc := strings.TrimSpace(code.MarketCode)
		marketCode = proto.MarketType(stat.AnyToInt64(mc))
		securityCode := proto.GetSecurityCode(marketCode, strings.TrimSpace(code.StockCode))
		securityName := strings.TrimSpace(code.ShortName)
		//if securityCode == "sz300027" {
		//	fmt.Printf("\n%+v\n", v)
		//}
		notice := NoticeDetail{
			//Code         string `dataframe:"证券代码"`  // 证券代码
			Code: securityCode,
			//Name         string `dataframe:"证券名称"`  // 证券名称
			Name: securityName,
			//DisplayTime  string `dataframe:"显示时间"`  // 显示时间
			DisplayTime: strings.TrimSpace(v.EiTime),
			//DisplayTime: strings.TrimSpace(v.DisplayTime),
			//NoticeDate   string `dataframe:"公告时间"`  // 公告时间
			NoticeDate: strings.TrimSpace(v.NoticeDate),
			//Title        string `dataframe:"内容提要"`  // 公告标题
			Title: strings.TrimSpace(v.TitleCh),
			//Keywords     string `dataframe:"关键词"`   // 公告关键词
			//Increase     int    `dataframe:"增持"`    // 增持
			//Reduces       int    `dataframe:"减持"`    // 减持
			//HolderChange int    `dataframe:"控制人变更"` // 实际控制人变更
		}
		noticeKeywords := []string{}

		checkRisk := func(content string) {
			key := "减持"
			if strings.Contains(content, key) {
				noticeKeywords = append(noticeKeywords, key)
				notice.Reduce += 1
			}
			key = "增持"
			if strings.Contains(content, key) {
				noticeKeywords = append(noticeKeywords, key)
				notice.Increase += 1
			}
			key = "控制人变更"
			if strings.Contains(content, key) {
				noticeKeywords = append(noticeKeywords, key)
				notice.HolderChange += 1
			}
			for _, key := range riskKeywords {
				if strings.Contains(content, key) {
					noticeKeywords = append(noticeKeywords, key)
					notice.Risk += 1
				}
			}
		}

		for _, words := range v.Columns {
			//if securityCode == "sh600730" {
			//	fmt.Println(securityCode, words.ColumnName)
			//}
			checkRisk(words.ColumnName)
		}
		checkRisk(notice.Title)
		if len(noticeKeywords) > 0 {
			notice.Keywords = strings.Join(noticeKeywords, ",")
		}

		notices = append(notices, notice)
	}
	return notices, pages, nil
}

// StockNotices 个股公告
func StockNotices(securityCode, beginDate, endDate string, pageNumber ...int) (notices []NoticeDetail, pages int, err error) {
	pageNo := 1
	if len(pageNumber) > 0 {
		pageNo = pageNumber[0]
	}
	beginDate = trading.FixTradeDate(beginDate)
	if len(endDate) > 0 {
		endDate = trading.FixTradeDate(endDate)
	} else {
		endDate = trading.Today()
	}
	pageSize := EastmoneyNoticesPageSize
	params := urlpkg.Values{
		"sr":         {"-1"},
		"page_size":  {fmt.Sprintf("%d", pageSize)},
		"page_index": {fmt.Sprintf("%d", pageNo)},
		//"ann_type":   {"SHA,CYB,SZA,BJA"},
		"ann_type": {"A"},
		//"ann_type":      {"SHA,SZA"},
		"client_source": {"web"},
		"f_node":        {"0"},
		//"f_node":     {fmt.Sprintf("%d", NoticeWarning)},
		"s_node":     {"0"},
		"begin_time": {beginDate},
		"end_time":   {endDate},
		//"cb": {"jQuery112305241416374967685_1683838825141"},
	}
	_, _, symbol := proto.DetectMarket(securityCode)
	params.Add("stock_list", symbol)
	// Host: np-anotice-stock.eastmoney.com
	header := map[string]any{
		//"User-Agent": config.HTTP_REQUEST_HEADER_USER_AGENT,
		//"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	}
	url := kUrlEastmoneyNotices + "?" + params.Encode()
	//url = "https://np-anotice-stock.eastmoney.com/api/security/ann?cb=jQuery112305241416374967685_1683838825141&sr=-1&page_size=50&page_index=1&ann_type=SHA%2CCYB%2CSZA%2CBJA&client_source=web&f_node=0&s_node=0"
	data, err := http.Get(url, header)
	if err != nil {
		return
	}
	//fmt.Println(api.Bytes2String(data))
	var raw rawNoticePackage
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return
	}
	if raw.Success != 1 || len(raw.Data.List) == 0 {
		err = ErrNoticeNotFound
		return
	}
	//pages = int(math.Ceil(float64(raw.Data.TotalHits) / float64(EastmoneyNoticesPageSize)))
	pages = utils.GetPages(pageSize, raw.Data.TotalHits)

	for _, v := range raw.Data.List {
		marketCode := proto.MarketIdShenZhen
		if len(v.Codes) == 0 || len(v.Columns) == 0 {
			continue
		}
		code := v.Codes[0]
		mc := strings.TrimSpace(code.MarketCode)
		marketCode = proto.MarketType(stat.AnyToInt64(mc))
		securityCode := proto.GetSecurityCode(marketCode, strings.TrimSpace(code.StockCode))
		securityName := strings.TrimSpace(code.ShortName)
		//if securityCode == "sz300027" {
		//	fmt.Printf("\n%+v\n", v)
		//}
		notice := NoticeDetail{
			//Code         string `dataframe:"证券代码"`  // 证券代码
			Code: securityCode,
			//Name         string `dataframe:"证券名称"`  // 证券名称
			Name: securityName,
			//DisplayTime  string `dataframe:"显示时间"`  // 显示时间
			DisplayTime: strings.TrimSpace(v.EiTime),
			//DisplayTime: strings.TrimSpace(v.DisplayTime),
			//NoticeDate   string `dataframe:"公告时间"`  // 公告时间
			NoticeDate: strings.TrimSpace(v.NoticeDate),
			//Title        string `dataframe:"内容提要"`  // 公告标题
			Title: strings.TrimSpace(v.TitleCh),
			//Keywords     string `dataframe:"关键词"`   // 公告关键词
			//Increase     int    `dataframe:"增持"`    // 增持
			//Reduces       int    `dataframe:"减持"`    // 减持
			//HolderChange int    `dataframe:"控制人变更"` // 实际控制人变更
		}
		noticeKeywords := []string{}
		// 评估风险
		checkRisk := func(content string) {
			key := "减持"
			if strings.Contains(content, key) {
				noticeKeywords = append(noticeKeywords, key)
				notice.Reduce += 1
			}
			key = "增持"
			if strings.Contains(content, key) {
				noticeKeywords = append(noticeKeywords, key)
				notice.Increase += 1
			}
			key = "控制人变更"
			if strings.Contains(content, key) {
				noticeKeywords = append(noticeKeywords, key)
				notice.HolderChange += 1
			}
			for _, key := range riskKeywords {
				if strings.Contains(content, key) {
					noticeKeywords = append(noticeKeywords, key)
					notice.Risk += 1
				}
			}
		}

		for _, words := range v.Columns {
			//if securityCode == "sh600730" {
			//	fmt.Println(securityCode, words.ColumnName)
			//}
			checkRisk(words.ColumnName)
		}
		checkRisk(notice.Title)
		if len(noticeKeywords) > 0 {
			notice.Keywords = strings.Join(noticeKeywords, ",")
		}

		notices = append(notices, notice)
	}
	return notices, pages, nil
}
