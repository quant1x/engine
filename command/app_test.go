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

func TestUpdateApplicationName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateApplicationName(tt.args.name)
		})
	}
}

func TestInitCommands(t *testing.T) {
	InitCommands()

}
