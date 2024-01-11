package factors

import (
	"fmt"
	"testing"
)

func TestKLineShape(t *testing.T) {
	code := "sz002992"
	code = "sz300569"
	code = "sz300748"
	code = "sz002284"
	code = "sh603551"
	code = "sh603818"
	code = "sz300008"
	code = "sh600663"
	code = "sh600522"
	code = "sh603099"
	//code = "sz002043"
	//code = "sz300591"
	df := BasicKLine(code)
	df = df.Sub(0, -2)
	v := KLineShape(df, code)
	fmt.Println(v)
	fmt.Println(v & (KLineShapeDoji | KLineShapeShrinkToHalf))
}
