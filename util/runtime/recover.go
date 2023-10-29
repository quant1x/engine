package runtime

import (
	"fmt"
	"gitee.com/quant1x/gox/logger"
	"runtime/debug"
)

// CatchPanic 捕获panic
func CatchPanic() {
	if err := recover(); err != nil {
		s := string(debug.Stack())
		fmt.Printf("\nerr=%v, stack=%s\n", err, s)
		logger.Fatalf("%s 异常: %+v", ApplicationName(), err)
	}
}
