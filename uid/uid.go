package uid

import (
	"sync/atomic"
	"time"
)

const (
	NODE = 30
	SEQ  = 30
)

type Uid struct {
	nodeId uint8
	base   uint32
	nextId atomic.Uint64
}

// NewUid
// The generator supports up to 256 nodes,
// only allows a maximum of 1 instance per minute,
// and generates up to 1 billion IDs per instance
// 1  | 25                       |  | 8    |  | 30
// 0--00000000-00000000-00000000-0--00000000--00000000-00000000-00000000-000000
func NewUid(nodeId uint8) (u *Uid, err error) {
	utcTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	u = &Uid{
		nodeId: nodeId,
		base:   uint32(time.Now().Sub(utcTime).Minutes()),
	}
	u.nextId.Store(uint64(u.base)<<38 + uint64(u.nodeId)<<30)
	return
}
func (u *Uid) NodeId() uint8 {
	return u.nodeId
}
func (u *Uid) Base() uint32 {
	return u.base
}
func (u *Uid) CurrentId() uint64 {
	return u.nextId.Load()
}
func (u *Uid) Gen() uint64 {
	u.nextId.Add(1)
	return u.nextId.Load()
}
