package cache

import (
	"gitee.com/quant1x/gox/errors"
	"sync"
)

// DataPlugin 数据插件
type DataPlugin interface {
	Name() string
	Setup(config map[string]string) error
}

type Depend interface {
	DependOn() []string
}

var (
	ErrIsExists = errors.New("plugin is already exists")
)

var (
	pluginMutex    sync.Mutex
	mapDataPlugins = map[string]DataPlugin{}
	//setupStatus map[string]bool
)

// Register 注册插件
func Register(plugin DataPlugin) error {
	pluginMutex.Lock()
	defer pluginMutex.Unlock()
	_, ok := mapDataPlugins[plugin.Name()]
	if ok {
		return ErrIsExists
	}
	mapDataPlugins[plugin.Name()] = plugin
	return nil
}

// 获取所有注册插件
func loadPlugins() (plugin chan DataPlugin, setupStatus map[string]bool) {
	// 这里定义一个长度为10的队列
	var sortPlugin = make(chan DataPlugin, 10)
	setupStatus = map[string]bool{}

	// 所有的插件
	for name, plugin := range mapDataPlugins {
		sortPlugin <- plugin
		setupStatus[name] = false
	}

	return sortPlugin, setupStatus
}

// SetupPlugins 加载所有插件
func SetupPlugins(pluginChan chan DataPlugin, setupStatus map[string]bool) error {
	num := len(pluginChan)
	for num > 0 {
		plugin := <-pluginChan
		canSetup := true
		if deps, ok := plugin.(Depend); ok {
			depends := deps.DependOn()
			for _, dependName := range depends {
				if _, setuped := setupStatus[dependName]; !setuped {
					// 有未加载的插件
					canSetup = false
					break
				}
			}
		}

		// 如果这个插件能被setup
		if canSetup {
			_ = plugin.Setup(nil)
			setupStatus[plugin.Name()] = true
		} else {
			// 如果插件不能被setup, 这个plugin就塞入到最后一个队列
			pluginChan <- plugin
		}
	}
	return nil
}
