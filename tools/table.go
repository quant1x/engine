package tools

import (
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// tableRendition 兼容 tablewriter v0 默认风格的渲染样式
//
//	ASCII 边框(+---+ / |), 表头下有分隔线, 数据行之间无分隔线
var tableRendition = tw.Rendition{
	Borders: tw.Border{Left: tw.On, Right: tw.On, Top: tw.On, Bottom: tw.On},
	Symbols: tw.NewSymbols(tw.StyleASCII),
	Settings: tw.Settings{
		Separators: tw.Separators{
			ShowHeader:     tw.On,
			ShowFooter:     tw.On,
			BetweenColumns: tw.On,
			BetweenRows:    tw.Off,
		},
		Lines: tw.Lines{
			ShowTop:        tw.On,
			ShowBottom:     tw.On,
			ShowHeaderLine: tw.On,
			ShowFooterLine: tw.On,
		},
	},
}

// NewTable 创建一个表格, 样式与 tablewriter v0 的默认风格保持一致
func NewTable(w io.Writer) *tablewriter.Table {
	return tablewriter.NewTable(w, tablewriter.WithRendition(tableRendition))
}
