# go-uid

分布式uid生成器，
id基於當前時間分鐘數遞增
如：

v1:

    123632936791572481
    123632936791572482
    123632936791572483
    123632936791572484
    123632936791572485
    123632936791572486
    123632936791572487
    123632936791572488
    123632936791572489
    123632936791572490

適用場景:

1. 性能要求高
2. ID需要保持順序
3. 分佈式
4. 高可用
5. 程序再次啟動後，依然保持順序

不適用場景:

1. 需要隱藏生成數量
2. 隨機

## Usage

v1:

* 在分佈式場景下，最多支持256個節點。
* 同一個節點，保證順序。多個節點間，不保證順序。
* 同一個節點，一分鐘內只能啟動一個實例，若啟動多個實例，ID會重複。
* 強烈建议使用UnsafeGen，实际上很安全，但會在生成約10億個ID後，Base+1。

The generator supports up to 256 nodes,
only allows a maximum of 1 instance per minute,
and generates up to 1 billion IDs per instance

    1  |25 base minute            |  |8 node|  |30 seq
    0--00000000-00000000-00000000-0--00000000--00000000-00000000-00000000-000000

## Install

```shell
go get -u github.com/lizongying/go-uid
```

## Sample

[samples](samples)

```go
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

```