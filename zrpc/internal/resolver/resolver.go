package resolver

import (
	"github.com/sjclijie/go-zero/core/discov/consul"
	"github.com/sjclijie/go-zero/core/discov/direct"
	"github.com/sjclijie/go-zero/core/discov/etcdv3"
	"google.golang.org/grpc/resolver"
)

const (
	DirectScheme = "direct"
	EtcdScheme   = "etcd"
	ConsulScheme = "consul"
)

var (
	directBuilder = direct.NewBuilder(DirectScheme)
	etcdV3Builder = etcdv3.NewBuilder(EtcdScheme)
	consulBuilder = consul.NewBuilder(ConsulScheme)
)

func RegisterResolver() {
	resolver.Register(directBuilder)
	resolver.Register(etcdV3Builder)
	resolver.Register(consulBuilder)
}

type ResolverTarget struct {
	EtcdConf   etcdv3.EtcdConf
	DirectConf []string
	ConsulConf consul.ConsulConf
}

func (r *ResolverTarget) Build(scheme string) (target string) {
	switch scheme {
	case DirectScheme:
		target = r.buildDirectTarget()
	case EtcdScheme:
		target = r.buildEtcdTarget()
	case ConsulScheme:
		consulBuilder.SetConfig(r.ConsulConf)
		target = r.buildConsulTarget()
	}
	return
}

func (r *ResolverTarget) buildDirectTarget() string {
	return directBuilder.Target(r.DirectConf)
}

func (r *ResolverTarget) buildEtcdTarget() string {
	return etcdV3Builder.Target(r.EtcdConf)
}

func (r *ResolverTarget) buildConsulTarget() string {
	return consulBuilder.Target(r.ConsulConf)
}
