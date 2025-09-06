package tdxweb

import (
	"fmt"

	"gitee.com/quant1x/gox/http"
)

const (
	// https://wenda.tdx.com.cn/site/wenda/stock_index.html?message=%E8%9E%8D%E8%B5%84%E8%9E%8D%E5%88%B8
	//urlMarginTrading = "https://wenda.tdx.com.cn/TQL?Entry=JNLPSE.getAllCode&RI=6C07"
	urlMarginTrading = "https://wenda.tdx.com.cn/TQL?Entry=JNLPSE.getAllCode&RI=6BFD"
)

func MarginTrading() {
	//params := urlpkg.Values{
	//	"direction": {direction.String()},
	//	"code":      {fmt.Sprintf("%s.%s", symbol, strings.ToUpper(mflag))},
	//	"price":     {fmt.Sprintf("%f", price)},
	//	"volume":    {fmt.Sprintf("%d", volume)},
	//	"strategy":  {models.QmtStrategyName(model)},
	//	"remark":    {models.QmtOrderRemark(model)},
	//}
	//body := params.Encode()
	body := `[{"nlpse_id":"7318110250698020161","op_flag":1,"sec_code":"","order_field":"sec_code","dynamic_order":"","order_flag":"1","POS":"0","COUNT":"30","timestamps":0,"RANG":"AG"}]`
	//body = `[{"op_flag":1,"sec_code":"","order_field":"sec_code","dynamic_order":"","order_flag":"1","POS":"0","COUNT":"30","timestamps":0,"RANG":"AG"}]`
	body = `[{"nlpseId":"7318094960614448284","orderField":"chg","orderFlag":"0"}]`
	//logger.Infof("trader-order: %s", body)
	header := map[string]any{
		http.ContextType: "application/x-www-form-urlencoded" + "; charset=UTF-8",
		//"Cookie":         "Hm_lvt_5c4c948b141e4d66943a8430c3d600d0=1703193725; Hm_lpvt_5c4c948b141e4d66943a8430c3d600d0=1703720553; LST=10; ASPSessionID=3755195075085152819",
		"Cookie":  "Hm_lvt_5c4c948b141e4d66943a8430c3d600d0=1703193725; Hm_lpvt_5c4c948b141e4d66943a8430c3d600d0=1703720553; LST=10; ASPSessionID=3755195113739858633",
		"Origin":  "https://wenda.tdx.com.cn",
		"Referer": "https://wenda.tdx.com.cn/",
	}

	data, _, err := http.Request(urlMarginTrading, http.MethodPost, body, header)

	fmt.Println(data, err)
}
