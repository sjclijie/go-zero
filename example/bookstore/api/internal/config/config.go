package config

import (
	"github.com/sjclijie/go-zero/rest"
	"github.com/sjclijie/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Add zrpc.RpcClientConf
}
