package utils

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"strings"
	"testing"
)

var propertyList = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>KeepAlive</key>
	<true/>
	<key>Label</key>
	<string>{{.Name}}</string>
	<key>ProgramArguments</key>
	<array>
	    <string>{{.Path}}</string>
		<string>daemon</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
    <key>WorkingDirectory</key>
    <string>${ROOT_PATH}</string>
    <key>StandardErrorPath</key>
    <string>${LOG_PATH}/{{.Name}}.err</string>
    <key>StandardOutPath</key>
    <string>${LOG_PATH}/{{.Name}}.log</string>
</dict>
</plist>
`

func TestTemplate(t *testing.T) {
	replacer := strings.NewReplacer("${ROOT_PATH}", cache.GetRootPath(), "${LOG_PATH}", cache.GetLoggerPath())
	//v := strings.ReplaceAll(propertyList, "${PATH}", cache.GetRootPath())
	v := replacer.Replace(propertyList)
	fmt.Println(v)
}
