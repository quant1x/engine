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

const (
	// DefaultDataProvider 默认数据提供者
	DefaultDataProvider = "engine"
)

// DataAdapter 数据插件
type DataAdapter interface {
	Trait // 继承特性接口

	//// Kind 数据类型
	//Kind() Kind
	//// Key 数据关键词, key与cache落地强关联
	//Key() string
	//// Name 特性名称
	//Name() string
	//// Usage 控制台参数提示信息, 数据描述(data description)
	//Usage() string
	//// Owner 提供者
	//Owner() string
	//// Init 初始化, 接受context, 日期和证券代码作为入参
	//Init(ctx context.Context, date, securityCode string) error

	// DataCommand string
	//DataCommand // 控制台命令字接口
	//// Init 初始化, 加载配置信息
	//Init(ctx context.Context, date, securityCode string) error
	//// Check 数据校验
	//Check(cacheDate, featureDate string)

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

var (
	ErrAlreadyExists = errors.New("plugin is already exists")
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

//// 获取所有注册插件
//func loadPlugins() (plugin chan DataAdapter, setupStatus map[Type]bool) {
//	// 这里定义一个长度为10的队列
//	var sortPlugin = make(chan DataAdapter, 10)
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
//func SetupPlugins(pluginChan chan DataAdapter, setupStatus map[Type]bool) error {
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