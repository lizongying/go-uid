package v1

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

type StoreLocal struct {
	nodeId   uint32
	savePath string
}

func NewStoreLocal(nodeId uint32) (StoreBase, error) {
	return &StoreLocal{
		nodeId:   nodeId,
		savePath: filepath.Join(os.TempDir(), fmt.Sprintf("uid_base_%d.bin", nodeId)),
	}, nil
}

// Load reads persisted baseMinute from file
func (s *StoreLocal) Load() (base uint32, err error) {
	f, err := os.Open(s.savePath)
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

// Save persists baseMinute to file
func (s *StoreLocal) Save(base uint32) error {
	f, err := os.OpenFile(s.savePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, base)
	_, err = f.Write(buf)
	return err
}

// Remove deletes the persisted file if it exists.
func (s *StoreLocal) Remove() error {
	if _, err := os.Stat(s.savePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(s.savePath)
}
