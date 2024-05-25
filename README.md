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
3. 單點或分佈式
4. 高可用
5. 程序再次啟動後，依然保持順序
6. 可以搭配服務發布或者程序內直接使用

不適用場景:

1. 需要隱藏生成數量
2. 隨機

## Usage

v1:

* 在分佈式場景下，256個節點可以作為同一業務的生成器，保證高可用。當然也可以分組甚至單點使用。
* 同一個節點，保證順序。多個節點間，不能保證順序。
* 同一個節點，一分鐘內只能啟動一個實例，若啟動多個實例，ID會重複。

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
)

func main() {
	ug := uidv1.NewUid(0, nil)
	fmt.Println(ug.NodeId())
	fmt.Println(ug.Base())
	for i := 0; i < 10; i++ {
		id := ug.Gen()
		fmt.Println(id)
	}
}

```