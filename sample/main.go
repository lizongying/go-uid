package main

import (
	"fmt"
	"github.com/lizongying/go-uid/uid"
)

func main() {
	ug, err := uid.NewUid(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(ug.NodeId())
	fmt.Println(ug.Base())
	for i := 0; i < 10; i++ {
		id := ug.Gen()
		fmt.Println(id)
	}
}
