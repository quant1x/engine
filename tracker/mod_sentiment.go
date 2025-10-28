package tracker

import (
	"fmt"

	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/num"
	"github.com/fatih/color"
)

// MarketSentiment 市场情绪
func MarketSentiment() {
	//sh000001 := GetStrategySnapshot("sh000001")
	//sz399107 := GetStrategySnapshot("sz399107")
	//fmt.Println("上证上涨: ", sh000001.IndexUp)
	//fmt.Println("上证下跌: ", sh000001.IndexDown)
	//fmt.Println("深证上涨: ", sz399107.IndexUp)
	//fmt.Println("深证下跌: ", sz399107.IndexDown)
	//
	//fmt.Println("上涨: ", sh000001.IndexUp+sz399107.IndexUp)
	//fmt.Println("下跌: ", sh000001.IndexDown+sz399107.IndexDown)
	// 涨跌家数
	zdjs := "sh880005"
	////sh880005 := models.GetTickFromMemory("sh880005")
	//tdxApi := level1.GetApi()
	//defer tdxApi.Close()
	//stockShots, _ := tdxApi.GetSnapshot([]string{zdjs})
	//sh880005 := stockShots[0]
	sh880005 := models.GetTickFromMemory(zdjs)
	//fmt.Printf("%+v\n", sh880005)
	up := sh880005.BidVol1 + sh880005.BidVol2 + sh880005.BidVol3 + sh880005.BidVol4 + sh880005.BidVol5
	down := sh880005.AskVol1 + sh880005.AskVol2 + sh880005.AskVol3 + sh880005.AskVol4 + sh880005.AskVol5
	//fmt.Printf("市场情绪：%.2f\n", 100*num.ChangeRate(up+down, up))
	_, _ = fmt.Fprintf(color.Output, "\n市场情绪：%s\n", color.RedString("%.2f", 100*num.ChangeRate(up+down, up)))
}
