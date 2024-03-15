package cache

import (
	"fmt"
	"testing"
)

func TestFilename(t *testing.T) {
	date := "2023-09-28"
	code := "sh600105"
	filename := QuarterlyReportFilename(code, date)
	fmt.Println(filename)
}

//// import "github.com/syndtr/goleveldb/leveldb"
//func TestLevelDB(t *testing.T) {
//	db, err := leveldb.OpenFile("t1.db", nil)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer db.Close()
//	db.Put([]byte("a"), []byte("1"), nil)
//}
