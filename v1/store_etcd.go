package v1

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type StoreEtcd struct {
	nodeId  uint32
	key     string
	client  *clientv3.Client
	timeout time.Duration
}

func NewStoreEtcd(client *clientv3.Client, nodeId uint32) (StoreBase, error) {
	key := fmt.Sprintf("/uid/base/%d", nodeId)

	return &StoreEtcd{
		nodeId:  nodeId,
		key:     key,
		client:  client,
		timeout: 3 * time.Second,
	}, nil
}

func NewStoreEtcdFactory(endpoints []string) NewStoreBaseFunc {
	var (
		once   sync.Once
		client *clientv3.Client
		err    error
	)
	return func(nodeId uint32) (StoreBase, error) {
		once.Do(func() {
			client, err = clientv3.New(clientv3.Config{
				Endpoints:   endpoints,
				DialTimeout: 5 * time.Second,
			})
		})
		if err != nil {
			return nil, err
		}

		return NewStoreEtcd(client, nodeId)
	}
}

// Load reads persisted baseMinute from etcd
func (s *StoreEtcd) Load() (base uint32, err error) {
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

// Save persists baseMinute to etcd
func (s *StoreEtcd) Save(base uint32) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, base)

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_, err := s.client.Put(ctx, s.key, string(buf))
	return err
}

// Remove deletes the persisted etcd if it exists.
func (s *StoreEtcd) Remove() error {
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
