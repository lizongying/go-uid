package main

import (
	"fmt"
	uidv2 "github.com/lizongying/go-uid/v2"
)

func main() {
	ug := uidv2.NewUid(0)
	for i := 0; i < 10; i++ {
		id := ug.Gen() // Generate a new unique ID
		fmt.Println("Generated ID:", id)
	}
}
