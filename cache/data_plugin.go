package cache

import (
	"gitee.com/quant1x/gox/errors"
	"slices"
	"sync"
)

type Kind = uint64

const (
	PluginMaskBaseData Kind = 0x1000000000000000
	PluginMaskFeature  Kind = 0x2000000000000000
)

// DataPlugin 数据插件
type DataPlugin interface {
	// Kind 优先级排序字段, 潜在的依赖关系
	Kind() Kind
	// Key 字符串关键词
	Key() string
	// Usage 控制台参数提示信息
	Usage() string
	// Init 初始化, 加载配置信息
	Init(barIndex *int, date string) error
	// Print 控制台输出指定日期的数据
	Print(code string, date ...string)

	//Setup(config map[string]string) error
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

// DataItem 单行数据
type DataItem interface {
	GetDate() string         // 日期
	GetSecurityCode() string // 证券代码
}

var (
	ErrAlreadyExists = errors.New("plugin is already exists")
)

var (
	pluginMutex    sync.Mutex
	mapDataPlugins = map[Kind]DataPlugin{}
	//setupStatus map[string]bool
)

// Register 注册插件
func Register(plugin DataPlugin) error {
	pluginMutex.Lock()
	defer pluginMutex.Unlock()
	_, ok := mapDataPlugins[plugin.Kind()]
	if ok {
		return ErrAlreadyExists
	}
	mapDataPlugins[plugin.Kind()] = plugin
	return nil
}

//// 获取所有注册插件
//func loadPlugins() (plugin chan DataPlugin, setupStatus map[Type]bool) {
//	// 这里定义一个长度为10的队列
//	var sortPlugin = make(chan DataPlugin, 10)
//	setupStatus = map[Type]bool{}
//
//	// 所有的插件
//	for kind, plugin := range mapDataPlugins {
//		sortPlugin <- plugin
//		setupStatus[kind] = false
//	}
//
//	return sortPlugin, setupStatus
//}
//
//// SetupPlugins 加载所有插件
//func SetupPlugins(pluginChan chan DataPlugin, setupStatus map[Type]bool) error {
//	num := len(pluginChan)
//	for num > 0 {
//		plugin := <-pluginChan
//		canSetup := true
//		if deps, ok := plugin.(Depend); ok {
//			depends := deps.DependOn()
//			for _, dependName := range depends {
//				if _, setuped := setupStatus[dependName]; !setuped {
//					// 有未加载的插件
//					canSetup = false
//					break
//				}
//			}
//		}
//
//		// 如果这个插件能被setup
//		if canSetup {
//			_ = plugin.Setup(nil)
//			setupStatus[plugin.Type()] = true
//		} else {
//			// 如果插件不能被setup, 这个plugin就塞入到最后一个队列
//			pluginChan <- plugin
//		}
//	}
//	return nil
//}

// Plugins 按照类型标志位捡出数据插件
func Plugins(mask ...Kind) (list []DataPlugin) {
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

func PluginsWithName(pluginType Kind, keywords ...string) (list []DataPlugin) {
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

//// Get 从注册的数据插件中获取数据
//func Get(kind Type, securityCode string, date ...string) any {
//	data, ok := mapDataPlugins[kind]
//	if ok {
//		ptr := data.Get(securityCode, date...)
//		return ptr
//	}
//	return nil
//}

//// Get 从注册的数据插件中获取数据
//func Get(kind Type, securityCode string, date ...string) any {
//	data, ok := mapDataPlugins[kind]
//	if ok {
//		ptr := data.Get(securityCode, date...)
//		return ptr
//	}
//	return nil
//}
