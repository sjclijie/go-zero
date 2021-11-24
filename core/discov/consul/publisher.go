package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/sjclijie/go-zero/core/lang"
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

	if p.config.HealthCheck {
		registration.Check = &api.AgentServiceCheck{
			CheckID:                        p.serviceId,
			Interval:                       (5 * time.Second).String(),
			Timeout:                        (5 * time.Second).String(),
			GRPC:                           fmt.Sprintf("%v:%v/%v", p.listenHost, p.listenPort, p.config.Key),
			DeregisterCriticalServiceAfter: (30 * time.Second).String(),
		}
	}

	if err := p.client.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	proc.AddWrapUpListener(func() {
		p.Deregister()
	})

	return nil
}

func (p *Publisher) Deregister() error {
	return p.client.Agent().ServiceDeregister(p.serviceId)
}
