package main

import (
	"flag"
	"fmt"
	"github.com/sjclijie/go-zero/core/discov/etcdv3"
	"log"
	"time"
)

var value = flag.String("v", "value", "the value")

func main() {
	flag.Parse()

	client := etcdv3.NewPublisher([]string{"etcd.discovery:2379"}, "028F2C35852D", *value)
	if err := client.KeepAlive(); err != nil {
		log.Fatal(err)
	}
	defer client.Stop()

	for {
		time.Sleep(time.Second)
		fmt.Println(*value)
	}
}
