package trader

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/concurrent"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
	"time"
)

const (
	blacklistFilename = "safes.csv"
)

// SecureType 安全类型
//
//	白名单: Pure(纯净), Safe(安全), Trustworthy(值得信赖), Approved(已批准), Acceptable(可接受), Permitted(允许)
//	黑名单: Forbidden(禁止), Unsafe(不安全), Untrustworthy(不值得信赖), Rejected(被拒绝), Unacceptable(不可接受), Prohibited(禁止)
//	备选: SecureVariety, SecureType, SecureKind
type SecureType uint8

const (
	FreeTrading        SecureType = iota // 自由交易
	ProhibitForBuying                    // 禁止买入
	ProhibitForSelling                   // 禁止卖出
	ProhibitTrading    SecureType = 0xff // 禁止交易, 买入和卖出

	//NotForSale         SecureType = 16   // 非卖品
	//NotForBuy          SecureType = 86   // 非买品
)

var (
	mapSecureTypes = map[SecureType]string{
		FreeTrading:        "自由交易",
		ProhibitForBuying:  "禁止买入",
		ProhibitForSelling: "禁止卖出",
		ProhibitTrading:    "禁止交易",
	}
)

func UsageOfSecureType() string {
	keys := api.Keys(mapSecureTypes)
	slices.Sort(keys)
	var builder strings.Builder
	for _, typ := range keys {
		desc, _ := mapSecureTypes[typ]
		builder.WriteString(fmt.Sprintf("%d: %s\n", typ, desc))
	}
	return builder.String()
}

var (
	onceSafes  coroutine.PeriodicOnce
	mapSafes   = concurrent.NewTreeMap[string, SecureType]()
	mutexSafes sync.RWMutex
	//ctxSafes, _ = coroutine.GetContextWithCancel()
)

// BlackAndWhite 黑白名单
type BlackAndWhite struct {
	Code string     `name:"证券代码" dataframe:"code"`
	Type SecureType `name:"类型" dataframe:"type"`
}

//func notify() {
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		logger.Error(err)
//		return
//	}
//	defer watcher.Close()
//}

func getFileModTime(filename string) time.Time {
	fileStat, err := os.Lstat(filename)
	if err != nil {
		return time.Now()
	}
	return fileStat.ModTime()
}

// 监控黑白名单文件的变化
func notifyBlackAndWhiteList() {
	filename := path.Join(cache.GetRootPath(), blacklistFilename)
	if !api.FileExist(filename) {
		_ = api.CheckFilepath(filename, true)
		_ = os.WriteFile(filename, nil, 0644)
	}

	lastModTime := getFileModTime(filename)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			modTime := getFileModTime(filename)
			if modTime.After(lastModTime) {
				logger.Warnf("黑白名单文件有变化, 重新加载...")
				lazyLoadListOfBlackAndWhite()
				logger.Warnf("黑白名单文件有变化, 重新加载...OK")
				lastModTime = modTime
			}
		}
	}
}

func init() {
	go notifyBlackAndWhiteList()
}

// 加载黑白名单
func lazyLoadListOfBlackAndWhite() {
	mutexSafes.Lock()
	defer mutexSafes.Unlock()
	filename := path.Join(cache.GetRootPath(), blacklistFilename)
	var list []BlackAndWhite
	err := api.CsvToSlices(filename, &list)
	if err != nil || len(list) == 0 {
		return
	}
	mapSafes.Clear()
	for _, v := range list {
		securityCode := exchange.CorrectSecurityCode(v.Code)
		mapSafes.Put(securityCode, v.Type)
	}
}

// SyncLoadListOfBlackAndWhite 同步黑白名单
func SyncLoadListOfBlackAndWhite() {
	mutexSafes.Lock()
	defer mutexSafes.Unlock()
	filename := path.Join(cache.GetRootPath(), blacklistFilename)
	var list []BlackAndWhite
	mapSafes.Each(func(key string, value SecureType) {
		if value == FreeTrading {
			return
		}
		list = append(list, BlackAndWhite{Code: key, Type: value})
	})
	if len(list) > 0 {
		_ = api.SlicesToCsv(filename, list)
	}
}

// 校验和修正证券代码
func verify_and_correct_for_security_code(code string) (securityCode string) {
	onceSafes.Do(lazyLoadListOfBlackAndWhite)
	return exchange.CorrectSecurityCode(code)
}

// 检查是否禁止的类型
func checkNotVarietyType(code string, varietyType SecureType) bool {
	securityCode := verify_and_correct_for_security_code(code)
	v, ok := mapSafes.Get(securityCode)
	if !ok {
		return true
	}
	if v != varietyType && v != ProhibitTrading {
		return true
	}
	return false
}

// AddCodeToBlackList 新增黑白名单成分股
func AddCodeToBlackList(code string, secureType SecureType) {
	_, ok := mapSecureTypes[secureType]
	if !ok {
		fmt.Println()
	}
	securityCode := verify_and_correct_for_security_code(code)
	mapSafes.Put(securityCode, secureType)
	SyncLoadListOfBlackAndWhite()
}

// ProhibitTradingToBlackList 禁止交易 - 双向
func ProhibitTradingToBlackList(code string) {
	AddCodeToBlackList(code, ProhibitTrading)
}

func ProhibitBuyingToBlackList(code string) {
	AddCodeToBlackList(code, ProhibitForBuying)
}

func ProhibitSellingToBlackList(code string) {
	AddCodeToBlackList(code, ProhibitForSelling)
}

// CheckForBuy 检查是否可以买
func CheckForBuy(code string) bool {
	return checkNotVarietyType(code, ProhibitForBuying)
}

// CheckForSell 检查是否可卖
func CheckForSell(code string) bool {
	return checkNotVarietyType(code, ProhibitForSelling)
}
