package discovery

import (
	"log/slog"
	"time"

	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	schema      string
	EtcdAddr    []string
	DailTimeout time.Duration

	close chan struct{}
	watch clientv3.WatchChan

	client         *clientv3.Client
	keyPerfix      string
	serverAddrList []resolver.Address

	cc resolver.ClientConn
	logger *slog.Logger
}
