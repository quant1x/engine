package command

import (
	"fmt"

	"github.com/quant1x/engine/config"
	"github.com/quant1x/pkg/yaml"
	"github.com/quant1x/x/api"
	cmder "github.com/spf13/cobra"
)

// CmdConfig 显示配置信息
var CmdConfig = &cmder.Command{
	Use:   "config",
	Short: "显示配置信息",
	Run: func(cmd *cmder.Command, args []string) {
		data, err := yaml.Marshal(config.GlobalConfig)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(api.Bytes2String(data))
		}
	},
}
