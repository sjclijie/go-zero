package etcdv3

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
	"sync"
)

//implements grpc.naming.Resolver
type Resolver struct {
	client    *clientv3.Client
	prefix    string
	target    resolver.Target
	cc        resolver.ClientConn
	addresses map[string]resolver.Address
	sync.RWMutex
}

func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {
}

func (r *Resolver) Close() {
}

func (r *Resolver) Watcher() {
	r.addresses = make(map[string]resolver.Address)

	//主动拉取一次
	response, err := r.client.Get(context.Background(), r.prefix, clientv3.WithPrefix())
	if err == nil {
		for _, kv := range response.Kvs {
			r.setAddress(string(kv.Key), string(kv.Value))
		}
		r.cc.UpdateState(resolver.State{
			Addresses: r.getAddress(),
		})
	}

	spew.Dump(r.prefix)
	spew.Dump(r.getAddress())

	//watch
	watchChan := r.client.Watch(context.Background(), r.prefix, clientv3.WithPrefix())
	for response := range watchChan {
		for _, event := range response.Events {
			switch event.Type {
			case mvccpb.PUT:
				r.setAddress(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				r.delAddress(string(event.Kv.Key))
			}
		}

		r.cc.UpdateState(resolver.State{
			Addresses: r.getAddress(),
		})

		spew.Dump(r.getAddress())
	}
}

func (r *Resolver) setAddress(key, address string) {
	r.Lock()
	defer r.Unlock()
	r.addresses[key] = resolver.Address{
		Addr: address,
	}
}

func (r *Resolver) delAddress(key string) {
	r.Lock()
	defer r.Unlock()
	delete(r.addresses, key)
}

func (r *Resolver) getAddress() []resolver.Address {

	addresses := make([]resolver.Address, 0, len(r.addresses))

	for _, address := range r.addresses {
		addresses = append(addresses, address)
	}

	return addresses
}
