package command

import (
	"fmt"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/permissions"
	"gitee.com/quant1x/engine/tracker"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	cmder "github.com/spf13/cobra"
	"strings"
)

const (
	trackerCommand     = "tracker"
	trackerDescription = "实时跟踪"
)

var (
	trackerStrategyCodes                = "1" // 策略编号
	CmdTracker           *cmder.Command = nil // 实时跟踪
)

func initTracker() {
	CmdTracker = &cmder.Command{
		Use:     trackerCommand,
		Example: Application + " " + trackerCommand + " --no=1",
		//Args:    cobra.MinimumNArgs(0),
		Args: func(cmd *cmder.Command, args []string) error {
			return nil
		},
		Short: trackerDescription,
		Long:  trackerDescription,
		Run: func(cmd *cmder.Command, args []string) {
			var strategyCodes []uint64
			array := strings.Split(trackerStrategyCodes, ",")
			for _, strategyNumber := range array {
				strategyNumber := strings.TrimSpace(strategyNumber)
				code := api.ParseUint(strategyNumber)
				// 1. 确定策略是否存在
				medel, err := models.CheckoutStrategy(code)
				if err != nil {
					fmt.Printf("策略编号%d, 不存在\n", code)
					logger.Errorf("策略编号%d, 不存在", code)
					continue
				}
				// 2. 确定策略是否有权限
				err = permissions.CheckPermission(medel)
				if err != nil {
					fmt.Printf("策略编号%d, 权限验证失败: %+v\n", code, err)
					logger.Errorf("策略编号%d, 权限验证失败: %+v", code, err)
					continue
				}
				strategyCodes = append(strategyCodes, code)
			}
			if len(strategyCodes) == 0 {
				fmt.Println("没有有效的策略编号, 实时扫描结束")
				logger.Info("没有有效的策略编号, 实时扫描结束")
				return
			}
			tracker.Tracker(strategyCodes...)
		},
	}

	CmdTracker.Flags().StringVar(&trackerStrategyCodes, "no", trackerStrategyCodes, "策略编号, 多个用逗号分隔")
}
