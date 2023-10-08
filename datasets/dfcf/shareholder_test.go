package dfcf

import (
	"fmt"
	"testing"
)

func TestFreeHoldingDetail(t *testing.T) {
	data := freeHoldingDetail()
	fmt.Println(data)
	//df := pandas.LoadStructs(data)
	//_ = df.WriteCSV("capitals-0.csv")
}

func TestFreeHoldingAnalyse(t *testing.T) {
	data, num, err := getFreeHoldingAnalyse(1)
	fmt.Println(num, err)
	fmt.Printf("%+v\n", data)
	//df := pandas.LoadStructs(data)
	//_ = df.WriteCSV("capitals-1.csv")
}
