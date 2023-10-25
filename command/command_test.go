package command

import (
	"gitee.com/quant1x/gox/errors"
	"testing"
)

func TestCommand(t *testing.T) {
	err := errors.New("invalid argument \"f10\" for \"--features\" flag: strconv.ParseBool: parsing \"f10\": invalid syntax")
	parseFlagError(err)
}
