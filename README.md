# go-uid

分布式uid生成器，
id基于当前时间分钟数递增
如：

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

适用场景:

1. 性能要求高
2. 顺序递增

不适用场景:

1. 需要隐藏数量信息

## Usage

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

[sample](./sample)

```go
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

```