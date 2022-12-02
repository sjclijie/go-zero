package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"time"
)

//implements grpc.naming.Resolver
type Resolver struct {
	authority string
	token     string

	cc resolver.ClientConn

	name      string
	lastIndex uint64
}

func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {
}

func (r *Resolver) Watcher() {
	config := api.DefaultConfig()
	config.Address = r.authority
	config.Token = r.token
	client, err := api.NewClient(config)

	if err != nil {
		fmt.Printf("error create consul client: %v\n", err)
		return
	}

	for {
		services, metaInfo, err := client.Health().Service(r.name, "", true, &api.QueryOptions{
			WaitIndex: r.lastIndex,
		})
		if err != nil {
			fmt.Printf("error retrieving instances from consul: %v\n", err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		r.lastIndex = metaInfo.LastIndex

		var addresses []resolver.Address
		for _, service := range services {
			addresses = append(addresses, resolver.Address{
				Addr: fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port),
			})
		}

		r.cc.UpdateState(resolver.State{
			Addresses: addresses,
		})
	}
}

func (r *Resolver) Close() {
}
