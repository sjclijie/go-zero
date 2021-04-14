package consul

import "errors"

type ConsulConf struct {
	Host        string
	Key         string
	HealthCheck bool `json:",default=true,optional"`
}

func (c ConsulConf) Validate() error {
	if len(c.Host) == 0 {
		return errors.New("empty consul host")
	} else if len(c.Key) == 0 {
		return errors.New("empty consul key")
	} else {
		return nil
	}
}
