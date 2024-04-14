package tracker

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gotdx/securities"
)

// SectorInfo 板块信息
type SectorInfo struct {
	Code           string   `name:"代码"`     // 板块代码
	Name           string   `name:"名称"`     // 板块名称
	Type           string   `name:"类型"`     // 板块类型
	OpenAmount     float64  `name:"开盘金额"`   // 开盘金额
	OpenChangeRate float64  `name:"开盘涨幅"`   // 开盘涨幅
	ChangeRate     float64  `name:"板块涨幅"`   // 板块涨幅
	Rank           int      `name:"板块排名"`   // 板块排名
	TopCode        string   `name:"领涨个股"`   // 领涨个股
	TopName        string   `name:"领涨个股名称"` // 领涨个股名称
	TopRate        float64  `name:"领涨个股涨幅"` // 领涨个股涨幅
	Count          int      `name:"总数"`     // 总数
	LimitUpNum     int      `name:"涨停数"`    // 涨停数
	NoChangeNum    int      `name:"平盘数"`    // 平盘数
	UpCount        int      `name:"上涨家数"`   // 上涨家数
	DownCount      int      `name:"下跌家数"`   // 下跌家数
	Capital        float64  // 流通盘
	FreeCapital    float64  // 自由流通股本
	OpenTurnZ      float64  `name:"开盘换手"`   // 开盘换手
	StockCodes     []string `dataframe:"-"` // 股票代码
}

// 通过板块类型扫描板块
//
//	板块排名从1开始
func scanBlockByType(pbarIndex *int, blockType securities.BlockType, rule *config.StrategyParameter) []SectorInfo {
	bs := []SectorInfo{}
	isHead := rule.Flag == models.OrderFlagHead
	blocks := scanSectorSnapshots(pbarIndex, blockType, isHead)
	rank := 0
	for i := 0; i < len(blocks); i++ {
		v := blocks[i]
		// 获取板块内个股列表
		blockInfo := securities.GetBlockInfo(v.SecurityCode)
		stockCodes := blockInfo.ConstituentStocks
		stockCount := len(stockCodes)
		if stockCount == 0 {
			continue
		}
		rank++
		bi := SectorInfo{
			Code:           v.SecurityCode,
			Name:           v.Name,
			Type:           BlockTypeName(v.SecurityCode),
			OpenAmount:     float64(v.IndexOpenAmount),
			OpenChangeRate: v.OpeningChangeRate,
			ChangeRate:     v.ChangeRate,
			Rank:           rank,
			OpenTurnZ:      v.OpenTurnZ,
			StockCodes:     stockCodes,
		}
		bs = append(bs, bi)
	}
	return bs
}

func scanBlockByTypeForTick(pbarIndex *int, blockType securities.BlockType) []SectorInfo {
	bs := []SectorInfo{}
	blocks := scanSectorSnapshots(pbarIndex, blockType, false)
	rank := 0
	for i := 0; i < len(blocks); i++ {
		v := blocks[i]
		// 获取板块内个股列表
		blockInfo := securities.GetBlockInfo(v.SecurityCode)
		stockCodes := blockInfo.ConstituentStocks
		stockCount := len(stockCodes)
		if stockCount == 0 {
			continue
		}
		rank++
		bi := SectorInfo{
			Code:           v.SecurityCode,
			Name:           v.Name,
			Type:           BlockTypeName(v.SecurityCode),
			OpenAmount:     float64(v.IndexOpenAmount),
			OpenChangeRate: v.OpeningChangeRate,
			ChangeRate:     v.ChangeRate,
			Rank:           rank,
			OpenTurnZ:      v.OpenTurnZ,
			StockCodes:     stockCodes,
		}
		bs = append(bs, bi)
	}
	return bs
}

// TopBlockWithType 板块排行
func TopBlockWithType(pbarIndex *int, rule *config.StrategyParameter) map[securities.BlockType][]SectorInfo {
	tmpMap := map[securities.BlockType][]SectorInfo{}
	blockTypes := []securities.BlockType{securities.BK_GAINIAN}
	for _, blockType := range blockTypes {
		blocks := scanBlockByType(pbarIndex, blockType, rule)
		tmpMap[blockType] = blocks
	}
	return tmpMap
}

var (
	// 缓存板块类型名称
	__mapBlockTypeName = map[string]string{}
)

func init() {
	_ = GetBlockList()
}

func GetBlockList() []string {
	// 执行板块指数的检测
	blockInfos := securities.BlockList()
	var blockCodes []string
	for _, v := range blockInfos {
		// 只保留行业和概念
		if v.Type != securities.BK_HANGYE && v.Type != securities.BK_GAINIAN {
			continue
		}
		blockCode := v.Code
		blockCodes = append(blockCodes, blockCode)
		blockTypeName, _ := securities.BlockTypeNameByTypeCode(v.Type)
		__mapBlockTypeName[blockCode] = blockTypeName
	}

	return blockCodes
}

func BlockTypeName(blockCode string) string {
	name, _ := __mapBlockTypeName[blockCode]
	return name
}
