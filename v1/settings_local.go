package v1

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SettingsLocal struct {
	nodeId   uint32
	savePath string
}

func NewSettingsLocal(nodeId uint32) (Settings, error) {
	return &SettingsLocal{
		nodeId:   nodeId,
		savePath: filepath.Join(os.TempDir(), fmt.Sprintf("uid_settings_%d.bin", nodeId)),
	}, nil
}

func NewSettingsLocalFactory(nodeId uint32) NewSettingsFunc {
	return func() (Settings, error) {
		return NewSettingsLocal(nodeId)
	}
}

func (s *SettingsLocal) CurrentTime() (time.Time, error) {
	return time.Now().UTC(), nil
}

func (s *SettingsLocal) NodeId() (nodeId uint32, err error) {
	return s.nodeId, nil
}

// LoadBase reads persisted baseMinute from file
func (s *SettingsLocal) LoadBase() (base uint32, err error) {
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

// SaveBase persists baseMinute to file
func (s *SettingsLocal) SaveBase(base uint32) error {
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

// RemoveBase deletes the persisted file if it exists.
func (s *SettingsLocal) RemoveBase() error {
	if _, err := os.Stat(s.savePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(s.savePath)
}
