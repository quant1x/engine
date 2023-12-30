package cache

import (
	"errors"
	"slices"
	"sync"
)

type Kind = uint64

const (
	PluginMaskBaseData Kind = 0x1000000000000000
	PluginMaskFeature  Kind = 0x2000000000000000
)

const (
	// DefaultDataProvider 默认数据提供者
	DefaultDataProvider = "engine"
)

// DataAdapter 数据插件
type DataAdapter interface {
	// Schema 继承基础特性接口
	Schema
	// Print 控制台输出指定日期的数据
	Print(code string, date ...string)
}

// Handover 缓存切换接口
type Handover interface {
	// ChangingOverDate 缓存数据转换日期
	//	数据集等基础数据不需要切换日期
	ChangingOverDate(date string)
}

type Depend interface {
	DependOn() []Kind
}

var (
	ErrAlreadyExists = errors.New("the plugin already exists")
)

var (
	pluginMutex    sync.Mutex
	mapDataPlugins = map[Kind]DataAdapter{}
	//setupStatus map[string]bool
)

// Register 注册插件
func Register(plugin DataAdapter) error {
	pluginMutex.Lock()
	defer pluginMutex.Unlock()
	_, ok := mapDataPlugins[plugin.Kind()]
	if ok {
		return ErrAlreadyExists
	}
	mapDataPlugins[plugin.Kind()] = plugin
	return nil
}

// Plugins 按照类型标志位捡出数据插件
func Plugins(mask ...Kind) (list []DataAdapter) {
	pluginMutex.Lock()
	defer pluginMutex.Unlock()
	pluginType := Kind(0)
	if len(mask) > 0 {
		if mask[0] == PluginMaskBaseData || mask[0] == PluginMaskFeature {
			pluginType = mask[0]
		}
	}
	// TODO: 这个地方的内存申请方面需要优化
	var kinds []Kind
	for kind, _ := range mapDataPlugins {
		if pluginType == 0 || kind&pluginType == pluginType {
			kinds = append(kinds, kind)
		}
	}
	slices.Sort(kinds)
	for _, kind := range kinds {
		plugin, ok := mapDataPlugins[kind]
		if ok {
			list = append(list, plugin)
		}
	}
	return
}

func PluginsWithName(pluginType Kind, keywords ...string) (list []DataAdapter) {
	pluginMutex.Lock()
	defer pluginMutex.Unlock()
	if len(keywords) == 0 {
		return
	}
	var kinds []Kind
	for kind, plugin := range mapDataPlugins {
		if kind&pluginType == pluginType && slices.Contains(keywords, plugin.Key()) {
			kinds = append(kinds, kind)
		}
	}
	if len(kinds) == 0 {
		return
	}
	slices.Sort(kinds)
	for _, kind := range kinds {
		plugin, ok := mapDataPlugins[kind]
		if ok {
			list = append(list, plugin)
		}
	}
	return
}
