package etcdv3

import (
	"fmt"
	"github.com/sjclijie/go-zero/core/utils"
	"google.golang.org/grpc/resolver"
	"strings"
)

type builder struct {
	endpointSepChar int32
	subsetSize      int
	scheme          string
}

func NewBuilder(scheme string) *builder {
	return &builder{
		endpointSepChar: ',',
		subsetSize:      32,
		scheme:          scheme,
	}
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	hosts := strings.FieldsFunc(target.Authority, func(r rune) bool {
		return r == b.endpointSepChar
	})

	sub, err := NewSubscriber(hosts, target.Endpoint)
	if err != nil {
		return nil, err
	}

	update := func() {
		var addrs []resolver.Address
		for _, val := range utils.Subset(sub.Values(), b.subsetSize) {
			addrs = append(addrs, resolver.Address{
				Addr: val,
			})
		}
		cc.UpdateState(resolver.State{
			Addresses: addrs,
		})
	}
	sub.AddListener(update)
	update()

	return &Resolver{cc: cc}, nil
}

func (b *builder) Scheme() string {
	return b.scheme
}

func (b *builder) Target(conf EtcdConf) string {
	return fmt.Sprintf("%s://%s/%s", b.scheme,
		strings.Join(conf.Hosts, fmt.Sprintf("%c", b.endpointSepChar)), conf.Key)
}
