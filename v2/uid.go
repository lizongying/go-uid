package v1

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/google/uuid"
)

type Uid struct {
	nodeId uint32 // Identifier for the node (supports up to 4294967295 nodes)[0-4294967295]
}

// NewUid
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                           time_high                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |           time_mid            |      time_low_and_version     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |clk_seq_hi_res |  clk_seq_low  |         node (0-1)            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         node (2-5)                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
func NewUid(nodeId uint32) (u *Uid) {
	u = &Uid{
		nodeId: nodeId,
	}

	b := make([]byte, 6)
	binary.BigEndian.PutUint32(b, nodeId)
	uuid.SetNodeID(b)
	return
}

func (u *Uid) Gen() (str string) {
	uuidValue, err := uuid.NewV6()
	if err != nil {
		return
	}

	bytes, _ := uuidValue.MarshalBinary()
	str = hex.EncodeToString(bytes)
	return
}
