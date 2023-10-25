package command

import (
	"strings"
)

func Init() {
	initPrint()
	initRepair()
	initUpdate()
	initRule()
	initBackTesting()
}

func parseFlagError(err error) (flag, value string) {
	before, _, ok := strings.Cut(err.Error(), "flag:")
	if !ok {
		return
	}
	before = strings.TrimSpace(before)
	//_, err1 := fmt.Sscanf(before, "invalid argument \"%s\" for \"--%s\"", &value, &flag)
	//if err1 != nil {
	//	return
	//}
	arr := strings.Split(before, "\"")
	if len(arr) != 5 {
		return
	}
	value = strings.TrimSpace(arr[1])
	flag = strings.TrimSpace(arr[3])
	arr = strings.Split(flag, "-")
	flag = arr[len(arr)-1]
	return
}
