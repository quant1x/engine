package command

import (
	"errors"
	"fmt"
	"testing"
)

func TestCommand(t *testing.T) {
	fmt.Println("Quant1X-" + FirstUpper("stock"))
	err := errors.New("invalid argument \"f10\" for \"--features\" flag: strconv.ParseBool: parsing \"f10\": invalid syntax")
	parseFlagError(err)
}
