# go-uid

分布式UID生成器，基於指定時間（分鐘）遞增

[go-uid](https://github.com/lizongying/go-uid)
___
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

對比Snowflake：

|              | Snowflake                  | go-uid                 |    |
|--------------|----------------------------|------------------------|----|
| 生成最大數量       | 每毫秒 2^12 = 4096，達上限需等待下一毫秒 | 幾乎無限制                  | ✔️ |
| 生成速度         | 快                          | 極快                     | ✔️ |
| Node 節點數量    | 10 bit → 1024 節點           | 8 bit → 256 節點         | ❌️ |
| ID 是否可反應生成時間 | 是                          | 否（只反應初始時間）             | ❌  |
| 單節點順序性       | 良好，但受毫秒 sequence 限制        | 完全順序                   | ✔️ |
| 多節點順序性       | 良好，但受毫秒 sequence 限制        | 差                      | ❌️ |
| 時間回撥問題       | 需要處理                       | 無。基於 baseMinute 設計天然解決 | ✔️ |

適用場景:

- 順序
- 分佈式
- 高性能
- 高可用

不適用場景:

- 无序

## Usage

v1:

- 如果 baseTime 為 nil，則默認使用 2025-01-01 00:00:00 UTC。
- baseTime 不得早於上次生成器使用的時間，以避免回撥導致 ID 重複。
- 如需允許回撥，請手動刪除臨時文件：/tmp/base_minute_{nodeId}.bin
- 最多支持256個節點。
- 同一節點，保證順序。多節點，不保證順序。

```
  1  |25 base minute | |8 node| |30 seq
  0--00000000-00000000-00000000-0--00000000--00000000-00000000-00000000-000000

```

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
	baseTimeStr  = "2025-01-01 00:00:00"
	locationName = "UTC"
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
	fmt.Println("Base Time (minutes since 2025-01-01 00:00:00 UTC):", ug.Base())

	// Generate and print 10 unique IDs
	for i := 0; i < 10; i++ {
		id := ug.Gen() // Generate a new unique ID
		fmt.Println("Generated ID:", id)
	}
}

```