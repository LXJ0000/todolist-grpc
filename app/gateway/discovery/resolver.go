package discovery

import (
	"context"
	"log/slog"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

const schema = "etcd"

type Resolver struct {
	schema      string
	EtcdAddr    []string
	DialTimeout time.Duration

	close   chan struct{}
	watchCh clientv3.WatchChan

	client         *clientv3.Client
	keyPerfix      string
	serverAddrList []resolver.Address

	cc     resolver.ClientConn
	logger *slog.Logger
}

// NewResolver create a new resolver.Builder base on etcd
func NewResolver(etcdAddr []string, logger *slog.Logger) *Resolver {
	return &Resolver{
		schema:      schema,
		EtcdAddr:    etcdAddr,
		DialTimeout: 3 * time.Second,
		logger:      logger,
	}
}

// Scheme returns the scheme supported by this resolver.
func (r *Resolver) Scheme() string {
	return r.schema
}

// Build creates a new resolver.Resolver for the given target
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	server := Server{Name: target.Endpoint(), Version: target.URL.Host}
	r.keyPerfix = server.Path()
	if _, err := r.start(); err != nil {
		return nil, err
	}
	return r, nil
}

// ResolveNow resolver.Resolver interface
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close resolver.Resolver interface
func (r *Resolver) Close() {
	r.close <- struct{}{}
}

// start
func (r *Resolver) start() (chan<- struct{}, error) {
	var err error
	r.client, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddr,
		DialTimeout: r.DialTimeout,
	})
	if err != nil {
		return nil, err
	}
	resolver.Register(r)

	r.close = make(chan struct{})

	if err = r.sync(); err != nil {
		return nil, err
	}

	go r.watch()

	return r.close, nil
}

// watch update events
func (r *Resolver) watch() {
	ticker := time.NewTicker(time.Minute)
	r.watchCh = r.client.Watch(context.Background(), r.keyPerfix, clientv3.WithPrefix())

	for {
		select {
		case <-r.close:
			return
		case res, ok := <-r.watchCh:
			if ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				r.logger.Error("sync failed", err)
			}
		}
	}
}

// update
func (r *Resolver) update(events []*clientv3.Event) {
	for _, ev := range events {
		var info Server
		var err error

		switch ev.Type {
		case clientv3.EventTypePut:
			info, err = UnmarshalServer(ev.Kv.Value)
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
			if !info.exist(r.serverAddrList, addr) {
				r.serverAddrList = append(r.serverAddrList, addr)
				r.cc.UpdateState(resolver.State{Addresses: r.serverAddrList})
			}
		case clientv3.EventTypeDelete:
			info, err = SplitPath(string(ev.Kv.Key))
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr}
			if s, ok := Remove(r.serverAddrList, addr); ok {
				r.serverAddrList = s
				r.cc.UpdateState(resolver.State{Addresses: r.serverAddrList})
			}
		}
	}
}

// sync 同步获取所有地址信息
func (r *Resolver) sync() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := r.client.Get(ctx, r.keyPerfix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	r.serverAddrList = []resolver.Address{}

	for _, v := range res.Kvs {
		info, err := UnmarshalServer(v.Value)
		if err != nil {
			continue
		}
		addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
		r.serverAddrList = append(r.serverAddrList, addr)
	}
	r.cc.UpdateState(resolver.State{Addresses: r.serverAddrList})
	return nil
}
