package main

import (
	"fmt"
	"github.com/sjclijie/go-zero/core/discov/etcdv3"
	"time"

	"github.com/sjclijie/go-zero/core/logx"
)

func main() {
	sub, err := etcdv3.NewSubscriber([]string{"etcd.discovery:2379"}, "028F2C35852D", etcdv3.Exclusive())
	logx.Must(err)

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("values:", sub.Values())
		}
	}
}
