package factors

import (
	"fmt"
	"gitee.com/quant1x/engine/factors/pb"
	"gitee.com/quant1x/pkg/yaml"
	"google.golang.org/protobuf/proto"
	"os"
	"testing"
)

func Test_updateChipDistribution(t *testing.T) {
	code := "000701"
	date := "2025-03-11"
	err := updateChipDistribution(date, code)
	if err != nil {
		fmt.Println(err)
	}
	filename := "t1.bin"
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	cd := pb.ChipDistribution{}
	err = proto.Unmarshal(dataBytes, &cd)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cd.String())
}

func Test_v1updateChipDistribution(t *testing.T) {
	code := "000701"
	date := "2025-03-11"
	err := updateChipDistribution(date, code)
	if err != nil {
		fmt.Println(err)
	}
	filename := "t1.yaml"
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	tmp := map[float64]float64{}
	err = yaml.Unmarshal(dataBytes, tmp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tmp)
}
