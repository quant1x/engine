package command

var (
	Application = "stock"
	MinVersion  = "0.0.1" // 版本号
	barIndex    = 1
)

// UpdateApplicationVersion 更新版本号
func UpdateApplicationVersion(v string) {
	MinVersion = v
}
