package v1

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	maxBase = 1<<25 - 1
)

type NewSettingsFunc func() (Settings, error)

type Settings interface {
	NodeId() (nodeId uint32, err error)
	CurrentTime() (currentTime time.Time, err error)
	LoadBase() (base uint32, err error)
	SaveBase(base uint32) (err error)
}

func NewUidLocal(sinceTime *time.Time, nodeBits uint8, nodeId uint32) (u *Uid, err error) {
	return NewUid(sinceTime, nodeBits, NewSettingsLocalFactory(nodeId))
}

func NewUidEtcd(sinceTime *time.Time, nodeBits uint8, endpoints []string) (u *Uid, err error) {
	return NewUid(sinceTime, nodeBits, NewSettingsEtcdFactory(endpoints))
}

type Uid struct {
	nodeBits uint8
	seqBits  uint8  // Bit position for sequence number
	nodeId   uint32 // Identifier for the node
	maxSeq   uint64
	base     atomic.Uint32 // Base time in minutes since a reference point
	nextId   atomic.Uint64 // Atomic counter for the next ID to be generated
	settings Settings
}

// NewUid creates a new Uid generator
// The ID format is as follows:
// 1 bit | 25 bits for base time in minutes (60 years) | 16 bits for node ID | 22 bits for sequence number
// 0--00000000-00000000-00000000-0--00000000--00000000-00000000-00000000-000000
// nodeBits A bit of position for node ID
func NewUid(sinceTime *time.Time, nodeBits uint8, newSettingsFunc NewSettingsFunc) (u *Uid, err error) {
	if nodeBits == 0 {
		nodeBits = 16
	}
	if nodeBits > 32 || nodeBits < 6 {
		return nil, fmt.Errorf("invalid nodeBits %d", nodeBits)
	}
	seqBits := 38 - nodeBits
	maxSeq := 1<<seqBits - 1
	u = &Uid{
		nodeBits: nodeBits,
		seqBits:  seqBits,
		maxSeq:   uint64(maxSeq),
	}

	if newSettingsFunc == nil {
		return nil, fmt.Errorf("newSettingsFunc is nil")
	}

	var settings Settings
	if settings, err = newSettingsFunc(); err != nil {
		return nil, err
	}

	if settings == nil {
		return nil, fmt.Errorf("settings is nil")
	}

	u.settings = settings
	if u.nodeId, err = settings.NodeId(); err != nil {
		return nil, err
	}

	if sinceTime == nil {
		t := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		sinceTime = &t
	}

	// Calculate the base time in minutes since the reference time
	currentTime, err := settings.CurrentTime()
	if err != nil {
		return nil, err
	}

	base := uint32(currentTime.Sub(sinceTime.UTC())/time.Minute) & maxBase

	// Try to read persisted baseMinute
	var baseMinute uint32
	if baseMinute, err = u.settings.LoadBase(); err == nil {
		if base <= baseMinute {
			base = baseMinute + 1
		}
	}

	u.base.Store(base)
	_ = u.settings.SaveBase(base)

	// Initialize the nextId with the base time and node ID
	u.nextId.Store(uint64(u.base.Load())<<(u.nodeBits+u.seqBits) + uint64(u.nodeId)<<u.seqBits)
	return
}

// NodeId returns the node ID of the generator
func (u *Uid) NodeId() uint32 {
	return u.nodeId
}

// Base returns the base time of the generator in minutes since the reference time
func (u *Uid) Base() uint32 {
	return u.base.Load()
}

// CurrentId returns the current value of the ID generator
func (u *Uid) CurrentId() uint64 {
	return u.nextId.Load()
}

// NextId generates a new unique ID by incrementing the atomic counter
func (u *Uid) NextId() uint64 {
	nextId := u.nextId.Add(1)
	if nextId&u.maxSeq == 0 {
		base := u.base.Add(1) & maxBase
		u.base.Store(base)
		_ = u.settings.SaveBase(base)
		nextId = uint64(base)<<(u.nodeBits+u.seqBits) + uint64(u.nodeId)<<u.seqBits + 1
		u.nextId.Store(nextId)
		return nextId
	}
	return nextId
}
