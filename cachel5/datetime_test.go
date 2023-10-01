package cachel5

import (
	"fmt"
	"testing"
)

func TestCorrectDate(t *testing.T) {
	date := "1991-12-19"
	c, s := CorrectDate(date)
	fmt.Println(c, s)
	date = "2023-07-03"
	c, s = CorrectDate(date)
	fmt.Println(c, s)
	date = "2023-07-10"
	c, s = CorrectDate(date)
	fmt.Println(c, s)
	date = "2023-07-09"
	c, s = CorrectDate(date)
	fmt.Println(c, s)
	date = "2023-07-08"
	c, s = CorrectDate(date)
	fmt.Println(c, s)
	date = "2023-07-07"
	c, s = CorrectDate(date)
	fmt.Println(c, s)
}
