package tools

import (
	"fmt"
	"gitee.com/quant1x/gox/util/homedir"
	"gitee.com/quant1x/pkg/tools/tail"
	"strings"
)

// TailFile 跟踪文件更新 tail -f
func TailFile(filename string, config tail.Config, done chan bool) {
	defer func() { done <- true }()
	filename, _ = homedir.Expand(filename)
	t, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
	err = t.Wait()
	if err != nil {
		fmt.Println(err)
	}
}

// TailFileWithNumber 查看最后n行数据
func TailFileWithNumber(filename string, config tail.Config, n int) {
	filename, _ = homedir.Expand(filename)
	t, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	builder := strings.Builder{}
	for line := range t.Lines {
		builder.WriteString(line.Text + "\n")
	}
	arr := strings.Split(builder.String(), "\n")
	total := len(arr)
	pos := n + 1
	if n >= total {
		n = total
	}
	//var lines []string
	for i, v := range arr {
		if i < total-pos {
			continue
		}
		fmt.Printf(v + "\n")
	}
}
