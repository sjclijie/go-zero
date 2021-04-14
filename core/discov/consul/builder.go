package consul

import (
	"fmt"
	"google.golang.org/grpc/resolver"
)

type builder struct {
	scheme          string
	endpointSepChar int32
	subsetSize      int
	Endpoints       []string
	config          ConsulConf
}

func NewBuilder(scheme string) *builder {
	return &builder{
		endpointSepChar: ',',
		subsetSize:      32,
		scheme:          scheme,
	}
}

func (b *builder) SetConfig(c ConsulConf) {
	b.config = c
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {

	r := &Resolver{
		authority: target.Authority,
		token:     b.config.Token,
		name:      target.Endpoint,
		cc:        cc,
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
