package svc

import (
	"bookstore/api/internal/config"
	"bookstore/api/internal/middleware"
	"bookstore/rpc/add/adder"
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/sjclijie/go-zero/rest"
	"github.com/sjclijie/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	AdminCheck rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	rpc := adder.NewAdder(zrpc.MustNewClient(c.Add))
	resp, err := rpc.Add(context.Background(), &adder.AddReq{Book: "hello", Price: 1})
	if err != nil {
		spew.Dump(err)
	}
	spew.Dump(resp)

	return &ServiceContext{
		Config:     c,
		AdminCheck: middleware.NewAdminCheck().Handle,
	}
}
