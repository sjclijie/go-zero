package zrpc

import (
	"github.com/sjclijie/go-zero/core/discov/consul"
	"github.com/sjclijie/go-zero/core/discov/etcdv3"
	"github.com/sjclijie/go-zero/core/service"
	"github.com/sjclijie/go-zero/core/stores/redis"
)

type (
	RpcServerConf struct {
		service.ServiceConf
		ListenOn      string
		Etcd          etcdv3.EtcdConf    `json:",optional"`
		Consul        consul.ConsulConf  `json:",optional"`
		Auth          bool               `json:",optional"`
		Redis         redis.RedisKeyConf `json:",optional"`
		StrictControl bool               `json:",optional"`
		// pending forever is not allowed
		// never set it to 0, if zero, the underlying will set to 2s automatically
		Timeout      int64 `json:",default=2000"`
		CpuThreshold int64 `json:",default=900,range=[0:1000]"`
	}

	RpcClientConf struct {
		Etcd      etcdv3.EtcdConf   `json:",optional"`
		Consul    consul.ConsulConf `json:",optional"`
		Endpoints []string          `json:",optional"`
		App       string            `json:",optional"`
		Token     string            `json:",optional"`
		Timeout   int64             `json:",optional"`
	}
)

func NewDirectClientConf(endpoints []string, app, token string) RpcClientConf {
	return RpcClientConf{
		Endpoints: endpoints,
		App:       app,
		Token:     token,
	}
}

func NewEtcdClientConf(hosts []string, key, app, token string) RpcClientConf {
	return RpcClientConf{
		Etcd: etcdv3.EtcdConf{
			Hosts: hosts,
			Key:   key,
		},
		App:   app,
		Token: token,
	}
}

func (sc RpcServerConf) HasEtcd() bool {
	err := sc.Etcd.Validate()
	return err == nil
}

func (sc RpcServerConf) HasConsul() bool {
	err := sc.Consul.Validate()
	return err == nil
}

func (sc RpcServerConf) Validate() error {
	if sc.Auth {
		if err := sc.Redis.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (cc RpcClientConf) HasCredential() bool {
	return len(cc.App) > 0 && len(cc.Token) > 0
}

func (cc RpcClientConf) HasEtcd() bool {
	err := cc.Etcd.Validate()
	return err == nil
}

func (cc RpcClientConf) HasConsul() bool {
	err := cc.Consul.Validate()
	return err == nil
}
