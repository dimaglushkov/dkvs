package internal

import (
	"github.com/dimaglushkov/dkvs/api/storagepb"
	"google.golang.org/grpc"
)

type handler struct {
	storagepb.StorageClient

	grpcConnector func(ip, port string) (*grpc.ClientConn, error)
}

func NewHandler() *handler {
	h := handler{}

	h.grpcConnector = grpcConnector

	return &h
}
