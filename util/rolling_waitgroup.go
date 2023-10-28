package util

import "sync"

// RollingWaitGroup 滑动窗口n的WaitGroup
type RollingWaitGroup struct {
	window    int
	c         chan struct{}
	waitGroup sync.WaitGroup
}

// NewRollingWaitGroup initialization RollingWaitGroup struct
func NewRollingWaitGroup(n int) *RollingWaitGroup {
	return &RollingWaitGroup{
		window: n,
		c:      make(chan struct{}, n),
	}
}

func (g *RollingWaitGroup) Add(delta int) {
	g.c <- struct{}{}
	g.waitGroup.Add(delta)
}

func (g *RollingWaitGroup) Done() {
	g.waitGroup.Done()
	<-g.c
}

func (g *RollingWaitGroup) Wait() {
	g.waitGroup.Wait()
}
