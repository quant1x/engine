package config

import (
	"testing"

	"github.com/quant1x/gox/encoding/yaml"
)

// Test_unmarshalConfig 内置配置文件的解析
func Test_unmarshalConfig(t *testing.T) {
	dataBytes, err := resources.ReadFile(ResourcesPath + "/" + configFilename)
	if err != nil {
		t.Fatal(err)
	}
	var conf Quant1XConfig
	if err = yaml.Unmarshal(dataBytes, &conf); err != nil {
		t.Fatal(err)
	}
	if len(conf.Trader.Strategies) != 2 {
		t.Fatalf("trader.strategies, 期望2项, 实际%d", len(conf.Trader.Strategies))
	}
	// yaml中显式配置的值, 以yaml为准
	if conf.Trader.Strategies[0].Total != 6 {
		t.Errorf("strategies[0].total, 期望6, 实际%d", conf.Trader.Strategies[0].Total)
	}
	// yaml中显式配置的零值, 不被default标签覆盖
	if conf.Trader.Strategies[1].Total != 0 {
		t.Errorf("strategies[1].total, 期望0, 实际%d", conf.Trader.Strategies[1].Total)
	}
	if conf.Trader.Strategies[0].Session.Size() != 1 {
		t.Errorf("strategies[0].time, 期望1个时段, 实际%d", conf.Trader.Strategies[0].Session.Size())
	}
	// 没有出现的字段, 取default标签的值
	if conf.Trader.TopN != 3 {
		t.Errorf("trader.top_n, 期望3, 实际%d", conf.Trader.TopN)
	}
	if conf.Data.Feature.Wave.Periods != 89 {
		t.Errorf("data.feature.wave.periods, 期望89, 实际%d", conf.Data.Feature.Wave.Periods)
	}
	if conf.Runtime.Crontab["realtime_kline"].Enable {
		t.Errorf("runtime.crontab.realtime_kline.enable, 期望false")
	}
	if !conf.Runtime.Crontab["sell_117"].Enable {
		t.Errorf("runtime.crontab.sell_117.enable, 期望true")
	}
}
