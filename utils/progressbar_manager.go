package utils

import "gitee.com/quant1x/gox/progressbar"

type ProgressBarManager struct {
	bar      *progressbar.Bar
	index    int
	modName  string
	totalNum int
}

func NewProgressBarManager(modName string, totalNum int) *ProgressBarManager {
	return &ProgressBarManager{
		modName:  modName,
		totalNum: totalNum,
	}
}

func (m *ProgressBarManager) Start() {
	m.index = 0
	if m.totalNum > 0 {
		m.bar = progressbar.NewBar(m.index, "执行["+m.modName+"]", m.totalNum)
	}
}

func (m *ProgressBarManager) Update(delta int) {
	if m.bar != nil {
		m.bar.Add(delta)
		m.index += delta
	}
}

func (m *ProgressBarManager) Wait() {
	if m.bar != nil {
		m.bar.Wait()
	}
}
