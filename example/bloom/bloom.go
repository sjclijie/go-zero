package main

import (
	"fmt"

	"github.com/sjclijie/go-zero/core/bloom"
	"github.com/sjclijie/go-zero/core/stores/redis"
)

func main() {
	store := redis.NewRedis("localhost:6379", "node")
	filter := bloom.New(store, "testbloom", 64)
	filter.Add([]byte("kevin"))
	filter.Add([]byte("wan"))
	fmt.Println(filter.Exists([]byte("kevin")))
	fmt.Println(filter.Exists([]byte("wan")))
	fmt.Println(filter.Exists([]byte("nothing")))
}
