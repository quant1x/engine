package tdxweb

import (
	"fmt"
	"gitee.com/quant1x/gox/http"
)

const (
	// https://wenda.tdx.com.cn/site/wenda/stock_index.html?message=%E8%9E%8D%E8%B5%84%E8%9E%8D%E5%88%B8
	urlMarginTrading = "https://wenda.tdx.com.cn/TQL?Entry=NLPSE.NLPQuery&RI=6C07"
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
	body := `[{"nlpse_id":"7311044342301379250","op_flag":1,"sec_code":"","order_field":"sec_code","dynamic_order":"","order_flag":"1","POS":"0","COUNT":"30","timestamps":0,"RANG":"AG"}]`
	//logger.Infof("trader-order: %s", body)
	data, err := http.Post(urlMarginTrading, body)

	fmt.Println(data, err)
}
