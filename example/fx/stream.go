package main

import (
	"fmt"
	"github.com/sjclijie/go-zero/core/fx"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ch := make(chan int)

	go inputStream(ch)
	go outputStream(ch)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
}

func inputStream(ch chan int) {
	count := 0
	for {
		ch <- count
		time.Sleep(time.Millisecond * 500)
		count++
	}
}

func outputStream(ch chan int) {
	fx.From(func(source chan<- interface{}) {
		for c := range ch {
			source <- c
		}
	}).Walk(func(item interface{}, pipe chan<- interface{}) {
		count := item.(int)
		pipe <- count
	}).Filter(func(item interface{}) bool {
		itemInt := item.(int)
		if itemInt%2 == 0 {
			return true
		}
		return false
	}).Head(5).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}
