package consul

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"net"
)

type builder struct {
	scheme          string
	endpointSepChar int32
	subsetSize      int
	Endpoints       []string
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

	host, port, err := net.SplitHostPort(target.Authority)
	if err != nil {
		return nil, err
	}

	r := &Resolver{
		host: host,
		port: port,
		name: target.Endpoint,
		cc:   cc,
	}

	go r.Watcher()
	r.ResolveNow(resolver.ResolveNowOptions{})

	return r, nil
}

func (b *builder) Scheme() string {
	return b.scheme
}

func (b *builder) Target(conf ConsulConf) string {
	return fmt.Sprintf("%s://%s/%s", b.scheme, conf.Host, conf.Key)
}
