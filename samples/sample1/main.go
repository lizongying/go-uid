package main

import (
	"fmt"
	uidv1 "github.com/lizongying/go-uid/v1"
	"time"
)

const (
	nodeId       = 0
	baseTimeStr  = "2023-01-01 00:00:00"
	locationName = "Asia/Shanghai"
)

func main() {
	location, err := time.LoadLocation(locationName)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	baseTime, err := time.ParseInLocation(time.DateTime, baseTimeStr, location)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	// Create a new Uid generator for node
	ug := uidv1.NewUid(nodeId, &baseTime)

	// Print the node ID of the generator
	fmt.Println("Node ID:", ug.NodeId())

	// Print the base time in minutes since the reference time
	fmt.Println("Base Time (minutes since 2023-01-01 00:00:00 UTC):", ug.Base())

	// Generate and print 10 unique IDs
	for i := 0; i < 10; i++ {
		id := ug.UnsafeGen() // Generate a new unique ID
		fmt.Println("Generated ID:", id)
	}
}
