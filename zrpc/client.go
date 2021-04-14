package zrpc

import (
	"github.com/sjclijie/go-zero/core/discov/etcdv3"
	"github.com/sjclijie/go-zero/zrpc/internal/resolver"
	"log"
	"time"

	"github.com/sjclijie/go-zero/zrpc/internal"
	"github.com/sjclijie/go-zero/zrpc/internal/auth"
	"google.golang.org/grpc"
)

var (
	WithDialOption             = internal.WithDialOption
	WithTimeout                = internal.WithTimeout
	WithUnaryClientInterceptor = internal.WithUnaryClientInterceptor
)

type (
	ClientOption = internal.ClientOption

	Client interface {
		Conn() *grpc.ClientConn
	}

	RpcClient struct {
		client Client
	}
)

func MustNewClient(c RpcClientConf, options ...ClientOption) Client {
	cli, err := NewClient(c, options...)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}

func NewClient(c RpcClientConf, options ...ClientOption) (Client, error) {
	var opts []ClientOption
	if c.HasCredential() {
		opts = append(opts, WithDialOption(grpc.WithPerRPCCredentials(&auth.Credential{
			App:   c.App,
			Token: c.Token,
		})))
	}
	if c.Timeout > 0 {
		opts = append(opts, WithTimeout(time.Duration(c.Timeout)*time.Millisecond))
	}
	opts = append(opts, options...)

	var client Client
	var err error

	resolverTarget := resolver.ResolverTarget{
		EtcdConf:   c.Etcd,
		DirectConf: c.Endpoints,
		ConsulConf: c.Consul,
	}

	if len(c.Endpoints) > 0 {
		client, err = internal.NewClient(resolverTarget.Build(resolver.DirectScheme), opts...)
	} else if c.HasEtcd() {
		client, err = internal.NewClient(resolverTarget.Build(resolver.EtcdScheme), opts...)
	} else if c.HasConsul() {
		client, err = internal.NewClient(resolverTarget.Build(resolver.ConsulScheme), opts...)
	}
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

func NewClientNoAuth(c etcdv3.EtcdConf, opts ...ClientOption) (Client, error) {

	return nil, nil

	/*
		//client, err := internal.NewClient(internal.BuildDiscovTarget(c.Hosts, c.Key), opts...)
		//client, err := internal.NewClient(  , opts...)

		if err != nil {
			return nil, err
		}

		return &RpcClient{
			client: client,
		}, nil
	*/
}

func NewClientWithTarget(target string, opts ...ClientOption) (Client, error) {
	return internal.NewClient(target, opts...)
}

func (rc *RpcClient) Conn() *grpc.ClientConn {
	return rc.client.Conn()
}
