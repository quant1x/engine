package factors

import (
	"gitee.com/quant1x/engine/datasource/dfcf"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"strings"
	"time"
)

type companyNotice struct {
	Increase     int
	Reduce       int
	Risk         int
	RiskKeywords string
}

// 上市公司公告
func getOneNotice(securityCode, currentDate string) (notice companyNotice) {
	if !exchange.AssertStockBySecurityCode(securityCode) {
		return
	}
	now, _ := api.ParseTime(currentDate)
	now = now.AddDate(0, -1, 0)
	beginDate := now.Format(time.DateOnly)
	endDate := currentDate
	//list, pages, err := dfcf.StockNotices(securityCode, beginDate, endDate, 1)
	//if pages < 1 {
	//	return
	//}
	pagesCount := 1
	var tmpNotice *dfcf.NoticeDetail = nil
	for pageNo := 1; pageNo < pagesCount+1; pageNo++ {
		list, pages, err := dfcf.StockNotices(securityCode, beginDate, endDate, pageNo)
		if err != nil || pages < 1 {
			logger.Errorf("notice: code=%s, %s=>%s, %s", securityCode, beginDate, endDate, err)
			break
		}
		if pagesCount < pages {
			pagesCount = pages
		}

		count := len(list)
		if count == 0 {
			break
		}
		for _, v := range list {
			if tmpNotice != nil {
				tmpNotice.Name = v.Name
				if tmpNotice.NoticeDate < v.NoticeDate {
					tmpNotice.DisplayTime = v.DisplayTime
					tmpNotice.NoticeDate = v.NoticeDate
				}
				// 使用最近的标题
				tmpNotice.Title = v.Title
				keywords := tmpNotice.Keywords
				if len(v.Keywords) > 0 {
					if len(keywords) == 0 {
						keywords += v.Keywords
					} else {
						keywords += "," + v.Keywords
					}
				}
				tmpArr := strings.Split(keywords, ",")
				//api.Unique(api.StringSlice{P: &tmpArr})
				tmpArr = api.Unique(tmpArr)
				tmpNotice.Keywords = strings.Join(tmpArr, ",")
				tmpNotice.Increase += v.Increase
				tmpNotice.Reduce += v.Reduce
				tmpNotice.HolderChange += v.HolderChange
				tmpNotice.Risk += v.Risk
			} else {
				tmpNotice = &v
			}
		}
		if count < dfcf.EastmoneyNoticesPageSize {
			break
		}
	}
	if tmpNotice != nil {
		notice.Increase = tmpNotice.Increase
		notice.Reduce = tmpNotice.Reduce
		notice.Risk = tmpNotice.Risk
		notice.RiskKeywords = tmpNotice.Keywords
	}
	return
}
