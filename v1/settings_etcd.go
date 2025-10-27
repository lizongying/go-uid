package v1

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type SettingsEtcd struct {
	nodeId  uint32
	key     string
	client  *clientv3.Client
	timeout time.Duration
	ttl     time.Duration
}

func NewSettingsEtcd(client *clientv3.Client) (Settings, error) {
	s := &SettingsEtcd{
		client:  client,
		timeout: 3 * time.Second,
		ttl:     time.Hour,
	}
	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

func NewSettingsEtcdFactory(endpoints []string) NewSettingsFunc {
	var (
		once   sync.Once
		client *clientv3.Client
		err    error
	)
	return func() (Settings, error) {
		once.Do(func() {
			client, err = clientv3.New(clientv3.Config{
				Endpoints:   endpoints,
				DialTimeout: 5 * time.Second,
			})
		})
		if err != nil {
			return nil, err
		}

		return NewSettingsEtcd(client)
	}
}

func (s *SettingsEtcd) init() (err error) {
	key := "/uid/nodes/"

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	resp, err := s.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("get nodes failed: %w", err)
	}

	ids := make([]int, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		id, _ := strconv.Atoi(strings.TrimPrefix(string(kv.Key), key))
		ids = append(ids, id)
	}

	nodeId := 0
	for _, id := range ids {
		if id != nodeId {
			break
		}
		nodeId++
	}

	leaseResp, err := s.client.Grant(ctx, int64(s.ttl/time.Second))
	if err != nil {
		return fmt.Errorf("grant lease failed: %w", err)
	}

	_, err = s.client.Put(ctx, path.Join(key, fmt.Sprintf("%d", nodeId)), "active", clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("register nodeId failed: %w", err)
	}

	ch, err := s.client.KeepAlive(context.Background(), leaseResp.ID)
	if err == nil {
		go func() {
			for range ch {
			}
		}()
	}

	s.nodeId = uint32(nodeId)
	s.key = fmt.Sprintf("/uid/settings/%d", nodeId)
	return
}

func (s *SettingsEtcd) CurrentTime() (time.Time, error) {
	return time.Now().UTC(), nil
}

func (s *SettingsEtcd) NodeId() (nodeId uint32, err error) {
	return s.nodeId, nil
}

// LoadBase reads persisted baseMinute from etcd
func (s *SettingsEtcd) LoadBase() (base uint32, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	resp, err := s.client.Get(ctx, s.key)
	if err != nil {
		return 0, err
	}

	if len(resp.Kvs) == 0 {
		return 0, nil
	}

	value := resp.Kvs[0].Value
	if len(value) != 4 {
		return 0, fmt.Errorf("invalid data length: %d", len(value))
	}

	return binary.BigEndian.Uint32(value), nil
}

// SaveBase persists baseMinute to etcd
func (s *SettingsEtcd) SaveBase(base uint32) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, base)

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_, err := s.client.Put(ctx, s.key, string(buf))
	return err
}

// RemoveBase deletes the persisted etcd if it exists.
func (s *SettingsEtcd) RemoveBase() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	resp, err := s.client.Delete(ctx, s.key)
	if err != nil {
		return err
	}
	if resp.Deleted == 0 {
		return os.ErrNotExist
	}
	return nil
}
