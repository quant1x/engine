package tdxweb

import (
	"fmt"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/concurrent"
	"gitee.com/quant1x/gox/http"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/pkg/fastjson"
)

const (
	urlRiskAssessment            = "http://page3.tdx.com.cn:7615/site/pcwebcall_static/bxb/json/"
	defaultSafetyScore           = 100
	defaultSafetyScoreOfNotFound = 100
	defaultSafetyScoreOfIgnore   = 0
)

var (
	__mapSafetyScore = concurrent.NewTreeMap[string, int]()
)

// GetSafetyScore 获取个股安全分
//
//	http://page3.tdx.com.cn:7615/site/pcwebcall_static/bxb/bxb.html?code=600178&color=0
func GetSafetyScore(securityCode string) (score int) {
	if !exchange.AssertStockBySecurityCode(securityCode) {
		return defaultSafetyScore
	}
	if market.IsNeedIgnore(securityCode) {
		return defaultSafetyScoreOfIgnore
	}
	score = defaultSafetyScore
	_, _, code := exchange.DetectMarket(securityCode)
	if len(code) == 6 {
		url := fmt.Sprintf("%s%s.json", urlRiskAssessment, code)
		data, err := http.Get(url)
		if err != nil || err == http.NotFound {
			score = defaultSafetyScoreOfNotFound
		} else {
			obj, err := fastjson.ParseBytes(data)
			if err != nil {
				logger.Errorf("%+v\n", err)
				tmpScore, ok := __mapSafetyScore.Get(securityCode)
				if ok {
					score = tmpScore
				} else {
					score = defaultSafetyScore
				}
			} else {
				result := obj.GetArray("data")
				if result != nil && len(result) > 0 {
					tmpScore := 100
					for _, v := range result {
						rows := v.GetArray("rows")
						for _, row := range rows {
							trig := row.GetInt("trig")
							if trig == 1 {
								tmpScore = tmpScore - row.GetInt("fs")
							}
						}
					}
					score = tmpScore
					__mapSafetyScore.Put(securityCode, score)
				}
			}
		}
	}
	return score
}
