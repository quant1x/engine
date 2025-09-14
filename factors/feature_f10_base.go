package factors

import (
	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/gotdx/quotes"
	"github.com/quant1x/gotdx/securities"
	"github.com/quant1x/x/api"
	"github.com/quant1x/x/concurrent"
)

func checkoutCapital(list []quotes.XdxrInfo, date string) *quotes.XdxrInfo {
	for _, v := range list {
		if v.IsCapitalChange() && date >= v.Date {
			return &v
		}
	}
	return nil
}

type f10SecurityInfo struct {
	TotalCapital float64
	Capital      float64
	VolUnit      int
	DecimalPoint int
	Name_        string
	IpoDate      string
	SubNew       bool
	UpdateDate   string
}

var (
	//__mapListingDate = map[string]string{}
	__mapListingDate = concurrent.NewHashMap[string, string]()
)

func checkoutSecurityBasicInfo(securityCode, featureDate string) f10SecurityInfo {
	list := base.GetCacheXdxrList(securityCode)
	api.SliceSort(list, func(a, b quotes.XdxrInfo) bool {
		return a.Date > b.Date
	})
	// 计算流通盘
	cover := checkoutCapital(list, featureDate)
	var f10 f10SecurityInfo
	if cover != nil {
		f10.TotalCapital = cover.HouZongGuBen * 10000
		f10.Capital = cover.HouLiuTong * 10000
	} else {
		f10.Capital, f10.TotalCapital, f10.IpoDate, f10.UpdateDate = getFinanceInfo(securityCode, featureDate)
	}
	if len(f10.IpoDate) == 0 {
		ipoDate, ok := __mapListingDate.Get(securityCode)
		if !ok {
			ipoDate = getIpoDate(securityCode, featureDate)
		}
		f10.IpoDate = ipoDate
		if len(f10.IpoDate) > 0 {
			__mapListingDate.Set(securityCode, f10.IpoDate)
		}
	}
	if len(f10.UpdateDate) == 0 || f10.UpdateDate > featureDate {
		f10.UpdateDate = featureDate
	}

	if len(f10.IpoDate) > 0 {
		f10.SubNew = IsSubNewStockByIpoDate(securityCode, f10.IpoDate, featureDate)
	}

	securityInfo, found := securities.CheckoutSecurityInfo(securityCode)
	if found {
		f10.VolUnit = int(securityInfo.VolUnit)
		f10.DecimalPoint = int(securityInfo.DecimalPoint)
		f10.Name_ = securityInfo.Name
	} else {
		f10.VolUnit = 100
		f10.DecimalPoint = 2
		f10.Name_ = securities.GetStockName(securityCode)
	}

	return f10
}
