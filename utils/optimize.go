package utils

import (
	"errors"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/pandas/stat"
	"golang.org/x/sys/cpu"
)

var (
	ErrAccelerationNotSupported = errors.New("acceleration not supported on this platform")
)

// Optimize 系统优化系列
func Optimize() {
	// 如果支持AVX2就打开
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
		stat.SetAvx2Enabled(true)
	} else {
		logger.Warn(ErrAccelerationNotSupported)
	}
}
