package config

import (
	"github.com/sjclijie/go-zero/core/stores/cache"
	"github.com/sjclijie/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Table      string
	Cache      cache.CacheConf
}
