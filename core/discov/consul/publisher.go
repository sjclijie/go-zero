package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/sjclijie/go-zero/core/lang"
	"github.com/sjclijie/go-zero/core/logx"
	"github.com/sjclijie/go-zero/core/proc"
	"github.com/sjclijie/go-zero/core/syncx"
	"net"
	"strconv"
	"time"
)

type (
	PublisherOption func(client *Publisher)
	Publisher       struct {
		client     *api.Client
		config     ConsulConf
		listenHost string
		listenPort int
		serviceId  string
		quit       *syncx.DoneChan
		pauseChan  chan lang.PlaceholderType
		resumeChan chan lang.PlaceholderType
	}
)

func NewPublisher(config ConsulConf, listenOn string, opts ...PublisherOption) (*Publisher, error) {

	host, portStr, err := net.SplitHostPort(listenOn)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	publisher := &Publisher{
		config:     config,
		serviceId:  fmt.Sprintf("%v-%v:%v", config.Key, host, port),
		listenHost: host,
		listenPort: port,
		quit:       syncx.NewDoneChan(),
		pauseChan:  make(chan lang.PlaceholderType),
		resumeChan: make(chan lang.PlaceholderType),
	}

	for _, opt := range opts {
		opt(publisher)
	}

	c := api.DefaultConfig()
	c.Address = config.Host
	c.Token = config.Token

	if client, err := api.NewClient(c); err != nil {
		return nil, err
	} else {
		publisher.client = client
	}

	return publisher, nil
}

func (p *Publisher) Register() error {

	registration := &api.AgentServiceRegistration{
		ID:      p.serviceId,
		Name:    p.config.Key,
		Tags:    []string{p.config.Key},
		Address: p.listenHost,
		Port:    p.listenPort,
	}

	ttl := fmt.Sprintf("%ds", 20)
	expiredTTL := fmt.Sprintf("%ds", 50)

	registration.Checks = []*api.AgentServiceCheck{
		{
			CheckID:                        p.serviceId,
			TTL:                            ttl,
			Status:                         "passing",
			DeregisterCriticalServiceAfter: expiredTTL,
		},
	}

	if err := p.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("[consul] service registration failed, %s", err.Error())
	} else {
		logx.Infof("[consul] service registration successful, %s", p.serviceId)
	}

	check := api.AgentServiceCheck{TTL: ttl, Status: "passing", DeregisterCriticalServiceAfter: expiredTTL}
	err := p.client.Agent().CheckRegister(&api.AgentCheckRegistration{ID: p.serviceId, Name: p.config.Key, ServiceID: p.serviceId, AgentServiceCheck: check})
	if err != nil {
		return fmt.Errorf("[consul] initial register service check to consul error: %s", err.Error())
	}

	ticker := time.NewTicker(time.Second * 10)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err = p.client.Agent().UpdateTTL(p.serviceId, "", "passing"); err != nil {
					logx.Errorf("[consul] update ttl of service error: %v", err.Error())

					if err := p.client.Agent().ServiceRegister(registration); err != nil {
						logx.Errorf("[consul] registration the service failed, service: %s,  error: %s", p.serviceId, err.Error())
					} else {
						logx.Infof("[consul] registration the service successful, service: %s", p.serviceId)
					}
				} else {
					logx.Info("[consul] update ttl successful")
				}
			}
		}
	}()

	proc.AddWrapUpListener(func() {
		ticker.Stop()
		p.Deregister()
	})

	return nil
}

func (p *Publisher) Deregister() error {
	return p.client.Agent().ServiceDeregister(p.serviceId)
}
