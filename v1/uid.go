package v1

import (
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
	nodeId uint8         // Identifier for the node (supports up to 256 nodes)[0-255]
	base   atomic.Uint32 // Base time in minutes since a reference point
	nextId atomic.Uint64 // Atomic counter for the next ID to be generated
}

// NewUid creates a new Uid generator
// It supports up to 256 nodes, only allows a maximum of 1 instance per minute,
// and generates up to 1 billion IDs per instance.
// The ID format is as follows:
// 1 bit | 25 bits for base time in minutes (60 years) | 8 bits for node ID | 30 bits for sequence number
// 0--00000000-00000000-00000000-0--00000000--00000000-00000000-00000000-000000
func NewUid(nodeId uint8, baseTime *time.Time) (u *Uid) {
	u = &Uid{
		nodeId: nodeId,
	}

	if baseTime != nil {
		// Calculate the base time in minutes since the reference time
		base := uint32(time.Now().Sub(*baseTime) / time.Minute)
		if base > maxBase {
			base = 0
		}
		u.base.Store(base)
	}

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
	return u.nextId.Add(1) // Increment the atomic counter. Return the new unique ID
}

// UnsafeGen generates a new unique ID by incrementing the atomic counter
func (u *Uid) UnsafeGen() uint64 {
	next := u.nextId.Add(1)
	if next&maxSeq == 0 {
		base := u.base.Add(1)
		if base > maxBase {
			base = 0
		}
		next = uint64(base)<<(node+seq) + uint64(u.nodeId)<<seq
		u.nextId.Store(next)
		return next
	}
	return next
}
