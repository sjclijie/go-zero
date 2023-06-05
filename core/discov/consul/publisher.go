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

	if p.config.HealthCheck {
		registration.Check = &api.AgentServiceCheck{
			CheckID:                        p.serviceId,
			Interval:                       (5 * time.Second).String(),
			Timeout:                        (5 * time.Second).String(),
			GRPC:                           fmt.Sprintf("%v:%v/%v", p.listenHost, p.listenPort, p.config.Key),
			DeregisterCriticalServiceAfter: (30 * time.Second).String(),
		}
	}

	ttl := fmt.Sprintf("%ds", 20)
	expireTTL := fmt.Sprintf("%ds", 60)

	registration.Checks = []*api.AgentServiceCheck{
		{
			CheckID:                        p.serviceId,
			TTL:                            ttl,
			Status:                         "passing",
			DeregisterCriticalServiceAfter: expireTTL,
		},
	}

	if err := p.client.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	check := api.AgentServiceCheck{TTL: ttl, Status: "passing", DeregisterCriticalServiceAfter: expireTTL}
	err := p.client.Agent().CheckRegister(&api.AgentCheckRegistration{ID: p.serviceId, Name: p.config.Key, ServiceID: p.serviceId, AgentServiceCheck: check})
	if err != nil {
		return fmt.Errorf("initial register service check to consul error: %s", err.Error())
	}

	/*
		target := fmt.Sprintf("%s://%s/%s", "consul", p.config.Host, p.config.Key)
		conn, err := grpc.Dial(target, grpc.WithInsecure())
		if err != nil {
			return err
		}

		healthClient := grpc_health_v1.NewHealthClient(conn)
	*/

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err = p.client.Agent().UpdateTTL(p.serviceId, "", "passing")
				logx.Info("update ttl")
				if err != nil {
					logx.Infof("update ttl of service error: %v", err.Error())
				}
				/*
					resp, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{
						Service: p.serviceId,
					})
					if err != nil || resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
						fmt.Printf("Service instance is not serving: %v\n", err)
					} else {
						fmt.Println("Service instance is serving")
					}
				*/
			}
		}
	}()

	proc.AddWrapUpListener(func() {
		p.Deregister()
	})

	return nil
}

func (p *Publisher) Deregister() error {
	return p.client.Agent().ServiceDeregister(p.serviceId)
}
