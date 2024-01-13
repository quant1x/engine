package command

import (
	"gitee.com/quant1x/engine/tracker"
	"gitee.com/quant1x/gox/api"
	cmder "github.com/spf13/cobra"
	"strings"
)

const (
	commandTracker = "tracker"
)

var (
	trackerStrategyCodes = "1" // 策略编号
)

// CmdTracker 实时跟踪
var CmdTracker = &cmder.Command{
	Use:     commandTracker,
	Example: Application + " " + commandTracker + " --no=1",
	//Args:    cobra.MinimumNArgs(0),
	Args: func(cmd *cmder.Command, args []string) error {
		return nil
	},
	Short: "实时跟踪",
	Long:  `实时跟踪`,
	Run: func(cmd *cmder.Command, args []string) {
		//if !CheckPermission(licenses.QuantStrategyNo81 | licenses.QuantStrategyNo82) {
		//	fmt.Println("没有策略权限")
		//	return
		//}
		var strategyCodes []int
		array := strings.Split(trackerStrategyCodes, ",")
		for _, strategyNumber := range array {
			strategyNumber := strings.TrimSpace(strategyNumber)
			code := api.ParseInt(strategyNumber)
			strategyCodes = append(strategyCodes, int(code))
		}
		tracker.Tracker(strategyCodes...)
	},
}

func initTracker() {
	CmdTracker.Flags().StringVar(&trackerStrategyCodes, "no", trackerStrategyCodes, "策略编号, 多个用逗号分隔")
}
