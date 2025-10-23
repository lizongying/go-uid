# go-uid

分布式UID生成器，順序、極速。

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

|            | Snowflake           | go-uid     |    |
|------------|---------------------|------------|----|
| 生成最大數量     | 每毫秒4096，達上限後需等待下一毫秒 | 幾乎無限制      | ✔️ |
| 生成速度       | 快                   | 極快         | ✔️ |
| Node節點數量   | 1024 節點             | 65536 節點   | ✔️ |
| ID是否反應生成時間 | 是                   | 否（只反應初始時間） | ❌  |
| 單節點順序性     | 良好，但受毫秒 sequence 限制 | 完全順序       | ✔️ |
| 多節點順序性     | 良好，但受毫秒 sequence 限制 | 差          | ❌️ |
| 時間回撥問題     | 有                   | 無          | ✔️ |

適用場景:

- 順序
- 分佈式
- 高性能
- 高可用

不適用場景:

- 无序

## Usage

v1:

- 32 >= nodeBits >= 6
- 如果 sinceTime 為 nil，則默認使用 2025-01-01 00:00:00 UTC。
- sinceTime 不得早於上次生成器使用的時間，避免回撥導致 ID 重複。
- 如需允許回撥，請手動在存儲中刪除。
- 同一節點，保證順序。多節點，不保證順序。
- 為了防止一分鐘內重啟ID重複問題，默認使用了StoreLocal。
- 為了防止nodeId轉移問題，可以使用StoreEtcd，也可以自己實現。
- 儘管生成數量幾乎無限制，但在默認配置下建議一分鐘內生成數量不要超過400萬，可以調整nodeBits來優化。

    - | nodeBits | 節點數量 | 建議單節點每分鐘最大生成數量 |
                  |----------|------------|----------------|
      | 6 | 64 | 4294967296     |
      | 8 | 256 | 1073741824     |
      | 16 | 65536 | 4194304        |
      | 32 | 4294967296 | 64             |

```
  1  |25 base minute            |  |16 node        |  |22 seq
  0--00000000-00000000-00000000-0--00000000-00000000--00000000-00000000-000000

```

## Install

```shell
go get -u github.com/lizongying/go-uid
```

## Sample

[samples](samples)

```go
ug := uidv1.Default(nodeId)
id := ug.NextId()
```

```go
package main

import (
	"fmt"
	uidv1 "github.com/lizongying/go-uid/v1"
	"time"
)

const (
	nodeId       = 0
	sinceTimeStr = "2025-01-01 00:00:00"
	locationName = "UTC"
)

func main() {
	location, err := time.LoadLocation(locationName)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	sinceTime, err := time.ParseInLocation(time.DateTime, sinceTimeStr, location)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	// Create a new Uid generator for node
	ug, _ := uidv1.NewUid(nodeId, &sinceTime, 16, nil, nil)

	// Print the node ID of the generator
	fmt.Println("Node ID:", ug.NodeId())

	// Print the base time in minutes since the reference time
	fmt.Println("Base Time (minutes since 2025-01-01 00:00:00 UTC):", ug.Base())

	// Generate and print 10 unique IDs
	for i := 0; i < 10; i++ {
		id := ug.NextId() // Generate a new unique ID
		fmt.Println("Generated ID:", id)
	}
}

```