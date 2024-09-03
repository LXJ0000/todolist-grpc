package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"go.etcd.io/etcd/client/v3"
)

type Register struct {
	EtcdAddr    []string
	DialTimeout time.Duration

	close       chan struct{}
	leasesID    clientv3.LeaseID                        // 租约
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse // 心跳

	serverInfo Server
	serverTTL  time.Duration

	client *clientv3.Client

	logger *slog.Logger
}

// NewRegister 基于 etcd 创建一个 register
func NewRegister(etcdAddr []string, dialTimeout time.Duration, logger *slog.Logger) *Register {
	return &Register{
		EtcdAddr:    etcdAddr,
		DialTimeout: dialTimeout,
		close:       make(chan struct{}),
		logger:      logger,
	}
}

// Register 创建封装的 etcd 实例
func (r *Register) Register(serverInfo Server, ttl time.Duration) (chan<- struct{}, error) {
	if strings.Split(serverInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid server addr")
	}
	var err error
	if r.client, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddr,
		DialTimeout: r.DialTimeout,
	}); err != nil {
		return nil, err
	}
	r.serverInfo = serverInfo
	r.serverTTL = ttl
	if err := r.register(); err != nil {
		return nil, err
	}
	go r.keepAlive()
	return r.close, nil
}

// register 创建 etcd 实例
func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.DialTimeout)
	defer cancel()
	leaseResp, err := r.client.Grant(ctx, int64(r.serverTTL.Seconds()))
	if err != nil {
		return err
	}
	r.leasesID = leaseResp.ID
	if r.keepAliveCh, err = r.client.KeepAlive(context.Background(), r.leasesID); err != nil {
		return err
	}
	bytes, err := json.Marshal(r.serverInfo)
	if err != nil {
		return err
	}
	_, err = r.client.Put(context.Background(), r.serverInfo.Path(), string(bytes), clientv3.WithLease(r.leasesID))
	return err
}

func (r *Register) unRegister() error {
	_, err := r.client.Delete(context.Background(), r.serverInfo.Path())
	return err
}

func (r *Register) keepAlive() {
	ticker := time.NewTicker(r.serverTTL)
	for {
		select {
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Warn("keep alive failed", slog.Any("err", err))
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Warn("keep alive failed", slog.Any("err", err))
				}
			}
		case <-r.close:
			if err := r.unRegister(); err != nil {
				r.logger.Warn("unregister server failed", slog.Any("err", err))
			}
			if _, err := r.client.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Warn("revoke lease failed", slog.Any("err", err))
			}
		}
	}
}
