package runtime

import (
	"fmt"
	"gitee.com/quant1x/engine/command"
	"gitee.com/quant1x/gox/logger"
	"runtime/debug"
)

func Recover() {
	if err := recover(); err != nil {
		s := string(debug.Stack())
		fmt.Printf("\nerr=%v, stack=%s\n", err, s)
		logger.Fatalf("%s 异常: %+v", command.Application, err)
	}
}
