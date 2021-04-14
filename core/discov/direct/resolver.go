package direct

import (
	"google.golang.org/grpc/resolver"
)

//implements grpc.naming.Resolver
type Resolver struct {
	cc resolver.ClientConn
}

func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {
}

func (r *Resolver) Close() {
}
