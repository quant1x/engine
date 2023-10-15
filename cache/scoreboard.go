package cache

import (
	"sync"
	"time"
)

// ScoreBoard 记分牌
type ScoreBoard struct {
	m         sync.Mutex
	Kind      Kind          // 类型
	Count     int           // 总数
	Max       time.Duration // 最大值
	Min       time.Duration // 最小值
	CrossTime time.Duration // 总耗时
	Speed     float64       // 速度
}

func (this *ScoreBoard) Add(delta int, take time.Duration) {
	this.m.Lock()
	defer this.m.Unlock()
	this.Count = this.Count + delta
	this.CrossTime += take
	if this.Min > take {
		this.Min = take
	}
	if this.Max < take {
		this.Max = take
	}
	this.Speed = float64(this.Count) / this.CrossTime.Seconds()
}
