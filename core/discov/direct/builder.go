package direct

import (
	"fmt"
	"github.com/sjclijie/go-zero/core/utils"
	"google.golang.org/grpc/resolver"
	"strings"
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
	var addrs []resolver.Address
	endpoints := strings.FieldsFunc(target.Endpoint, func(r rune) bool {
		return r == b.endpointSepChar
	})

	for _, val := range utils.Subset(endpoints, b.subsetSize) {
		addrs = append(addrs, resolver.Address{
			Addr: val,
		})
	}
	cc.UpdateState(resolver.State{
		Addresses: addrs,
	})

	return &Resolver{
		cc: cc,
	}, nil
}

func (b *builder) Scheme() string {
	return b.scheme
}

func (b *builder) Target(endpoints []string) string {
	return fmt.Sprintf("%s:///%s", b.scheme,
		strings.Join(endpoints, fmt.Sprintf("%c", b.endpointSepChar)))
}
