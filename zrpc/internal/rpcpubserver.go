package internal

import (
	"github.com/sjclijie/go-zero/core/discov"
)

type rpcPubServer struct {
	publisher discov.Publisher
	Server
}

func NewRpcPubServer(publisher discov.Publisher, listenOn string, opts ...ServerOption) (Server, error) {
	server := rpcPubServer{
		publisher: publisher,
		Server:    NewRpcServer(listenOn, opts...),
	}
	return server, nil
}

func (ags rpcPubServer) Start(fn RegisterFn) error {
	if err := ags.publisher.Register(); err != nil {
		return err
	}
	return ags.Server.Start(fn)
}

func (ags rpcPubServer) Stop() error {
	return ags.publisher.Deregister()
}
