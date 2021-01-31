package main

import (
	"fmt"
	"time"

	"github.com/sjclijie/go-zero/core/discov"
	"github.com/sjclijie/go-zero/core/logx"
)

func main() {
	sub, err := discov.NewSubscriber([]string{"etcd.discovery:2379"}, "028F2C35852D", discov.Exclusive())
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
