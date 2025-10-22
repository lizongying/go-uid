package v1

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

// Constants defining bit positions for node ID and sequence number
const (
	node    = 8  // A bit of position for node ID
	seq     = 30 // Bit position for sequence number
	maxBase = 1<<25 - 1
	maxSeq  = 1<<30 - 1
)

type Uid struct {
	nodeId   uint8         // Identifier for the node (supports up to 256 nodes)[0-255]
	base     atomic.Uint32 // Base time in minutes since a reference point
	nextId   atomic.Uint64 // Atomic counter for the next ID to be generated
	savePath string
}

// NewUid creates a new Uid generator
// It supports up to 256 nodes, only allows a maximum of 1 instance per minute,
// and generates up to 1 billion IDs per instance.
// The ID format is as follows:
// 1 bit | 25 bits for base time in minutes (60 years) | 8 bits for node ID | 30 bits for sequence number
// 0--00000000-00000000-00000000-0--00000000--00000000-00000000-00000000-000000
func NewUid(nodeId uint8, baseTime *time.Time) (u *Uid) {
	u = &Uid{
		nodeId:   nodeId,
		savePath: filepath.Join(os.TempDir(), fmt.Sprintf("base_minute_%d.bin", nodeId)),
	}

	if baseTime == nil {
		t := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		baseTime = &t
	}

	// Calculate the base time in minutes since the reference time
	base := uint32(time.Now().UTC().Sub(baseTime.UTC())/time.Minute) & maxBase

	// Try to read persisted baseMinute
	if b, err := u.loadBase(); err == nil {
		if base <= b {
			base = b + 1
		}
	}

	u.base.Store(base)
	_ = u.saveBase()

	// Initialize the nextId with the base time and node ID
	u.nextId.Store(uint64(u.base.Load())<<(node+seq) + uint64(u.nodeId)<<seq)
	return
}

// NodeId returns the node ID of the generator
func (u *Uid) NodeId() uint8 {
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

// Gen generates a new unique ID by incrementing the atomic counter
func (u *Uid) Gen() uint64 {
	nextId := u.nextId.Add(1)
	if nextId&maxSeq == 0 {
		base := u.base.Add(1) & maxBase
		u.base.Store(base)
		_ = u.saveBase()
		nextId = uint64(base)<<(node+seq) + uint64(u.nodeId)<<seq + 1
		u.nextId.Store(nextId)
		return nextId
	}
	return nextId
}

// saveBaseMinute persists baseMinute to file
func (u *Uid) saveBase() error {
	f, err := os.OpenFile(u.savePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, u.base.Load())
	_, err = f.Write(buf)
	return err
}

// loadBase reads persisted baseMinute from file
func (u *Uid) loadBase() (uint32, error) {
	f, err := os.Open(u.savePath)
	if err != nil {
		return 0, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	buf := make([]byte, 4)
	_, err = f.Read(buf)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf), nil
}
