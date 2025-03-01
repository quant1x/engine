package cache

import (
	"fmt"
	"sync"
	"time"
)

// ScoreBoard 记分牌, 线程安全
type ScoreBoard struct {
	m sync.Mutex
	AdapterMetric
}

// AdapterMetric 适配器性能指标
type AdapterMetric struct {
	Name      string        `name:"name"`       // 名称
	Kind      Kind          `name:"kind"`       // 类型
	Count     int           `name:"count"`      // 总数
	Max       time.Duration `name:"max"`        // 最大值
	Min       time.Duration `name:"min"`        // 最小值
	CrossTime time.Duration `name:"cross_time"` // 总耗时
	Speed     float64       `name:"speed"`      // 速度
}

func (this *ScoreBoard) From(adapter DataAdapter) {
	this.Name = adapter.Name()
	this.Kind = adapter.Kind()
}

func (this *ScoreBoard) Add(delta int, take time.Duration) {
	this.m.Lock()
	defer this.m.Unlock()
	this.Count = this.Count + delta
	this.CrossTime += take
	if this.Min == 0 || this.Min > take {
		this.Min = take
	}
	if this.Max == 0 || this.Max < take {
		this.Max = take
	}
	this.Speed = float64(this.Count) / this.CrossTime.Seconds()
}

func (this *ScoreBoard) String() string {
	s := fmt.Sprintf("name: %s, kind: %d, total: %d, crosstime: %s, max: %d, min: %d, speed: %f", this.Name, this.Kind, this.Count, this.CrossTime, this.Max, this.Min, this.Speed)
	return s
}

func (this *ScoreBoard) Metric() AdapterMetric {
	return this.AdapterMetric
}
