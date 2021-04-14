package redis

import (
	red "github.com/go-redis/redis"
	"github.com/sjclijie/go-zero/core/syncx"
	"io"
	"strings"
)

var clusterManager = syncx.NewResourceManager()

func getCluster(server, pass string) (*red.ClusterClient, error) {
	val, err := clusterManager.GetResource(server, func() (io.Closer, error) {
		store := red.NewClusterClient(&red.ClusterOptions{
			Addrs:        strings.Split(server, ","),
			Password:     pass,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
		})
		store.WrapProcess(process)
		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.ClusterClient), nil
}
