package factors

import (
	"github.com/quant1x/gotdx/quotes"
)

// BaseData 数据层, 数据集接口
//
//	数据集是基础数据, 应当遵循结构简单, 尽量减小缓存的文件数量, 加载迅速
//	检索的规则是按日期和代码进行查询
type BaseData interface {
	// Clone 克隆一个BaseData, 是所有写操作的基础
	Clone(featureDate, code string) BaseData
	// Update 更新数据
	Update(featureDate, code string) error
	// Repair 回补数据
	Repair(featureDate, code string) error
	// Increase 增量计算, 用快照增量计算特征
	Increase(snapshot quotes.Snapshot) error
}
