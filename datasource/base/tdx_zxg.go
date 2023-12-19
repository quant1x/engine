package base

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/pandas"
)

const (
	BlockPath = "/T0002/blocknew"
	ZxgBlk    = "zxg.blk"
	BkltBlk   = "BKLT.blk"
	ZdBk      = "ZDBK.blk"
)

func GetZxgList() []string {
	filename := cache.GetZxgFile()
	df := pandas.ReadCSV(filename)
	if df.Nrow() == 0 {
		return []string{}
	}
	rows := df.Col("code")
	if rows.Len() == 0 {
		return []string{}
	}
	// 校验证券代码, 统一格式前缀加代码
	cs := rows.Strings()
	for i, v := range cs {
		cs[i] = proto.CorrectSecurityCode(v)
	}
	return cs
}
