package cache

import (
	"fmt"
	"sync"
	"time"
)

// ScoreBoard 记分牌, 线程安全
type ScoreBoard struct {
	m sync.Mutex
	FactorMetrics
}

// FactorMetrics 适配器性能/绩效指标
type FactorMetrics struct {
	Name           string        `name:"name"`            // 名称
	Kind           Kind          `name:"kind"`            // 类型
	Count          int           `name:"count"`           // 总数(处理样本数)
	Signals        int           `name:"signals"`         // 信号数
	Passed         int           `name:"passed"`          // 通过检测数
	Max            time.Duration `name:"max"`             // 最大值
	Min            time.Duration `name:"min"`             // 最小值
	TotalDuration  time.Duration `name:"total_duration"`  // 总耗时
	Speed          float64       `name:"speed"`           // 速度 = Count / seconds
	SignalCoverage float64       `name:"signal_coverage"` // 信号覆盖率 = Signals / Count
	WinRate        float64       `name:"win_rate"`        // 胜率 = Passed / Signals
}

func (s *ScoreBoard) From(adapter DataAdapter) {
	s.Name = adapter.Name()
	s.Kind = adapter.Kind()
}

// recompute 派生指标
func (s *ScoreBoard) recompute() {
	if s.TotalDuration > 0 {
		secs := s.TotalDuration.Seconds()
		if secs > 0 {
			s.Speed = float64(s.Count) / secs
		}
	}
	if s.Count > 0 {
		s.SignalCoverage = float64(s.Signals) / float64(s.Count)
	} else {
		s.SignalCoverage = 0
	}
	if s.Signals > 0 {
		s.WinRate = float64(s.Passed) / float64(s.Signals)
	} else {
		s.WinRate = 0
	}
}

// Add 记录一次处理的性能及结果
// 参数说明:
//
//	delta  : 本次增加的样本数量(通常=1)
//	take   : 本次处理耗时
//	signal : 是否产生信号(hasSignal && err==nil)
//	pass   : 是否校验通过(err==nil)
//
// 统计逻辑:
//
//	Count   += delta
//	Signals += signal?1:0
//	Passed  += pass?1:0
//	SignalCoverage = Signals / Count
//	WinRate        = Passed / Signals (Signals==0时为0)
func (s *ScoreBoard) Add(delta int, take time.Duration, signal, pass bool) {
	s.m.Lock()
	defer s.m.Unlock()
	s.Count += delta
	s.TotalDuration += take
	if s.Min == 0 || s.Min > take {
		s.Min = take
	}
	if s.Max == 0 || s.Max < take {
		s.Max = take
	}
	if signal {
		s.Signals++
	}
	if pass {
		s.Passed++
	}

	s.recompute()
}

func (s *ScoreBoard) String() string {
	return fmt.Sprintf("name: %s, kind: %d, total: %d, signals: %d, passed: %d, coverage: %.4f, win: %.4f, crosstime: %s, max: %d, min: %d, speed: %.4f", s.Name, s.Kind, s.Count, s.Signals, s.Passed, s.SignalCoverage, s.WinRate, s.TotalDuration, s.Max, s.Min, s.Speed)
}

func (s *ScoreBoard) Metric() FactorMetrics {
	return s.FactorMetrics
}
