package svc

import (
	"bookstore/api/internal/config"
	"bookstore/api/internal/middleware"
	"github.com/sjclijie/go-zero/rest"
)

type ServiceContext struct {
	Config     config.Config
	AdminCheck rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		AdminCheck: middleware.NewAdminCheck().Handle,
	}
}
