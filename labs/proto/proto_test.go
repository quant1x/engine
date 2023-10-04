package proto

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	"os"
	"testing"
)

func TestProtobuf(t *testing.T) {
	filename := "t1.txt"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := History{Name: "1", Date: "2023-10-04", Price: 1.23}
	buff, err := proto.Marshal(&msg)
	defer file.Close()
	file.Write(buff)
	msg1 := History{}
	err1 := proto.Unmarshal(buff, &msg1)
	fmt.Println(err1)
	protoimpl.X.MessageOf(msg1)
}
