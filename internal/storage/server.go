package storage

import (
	context "context"
	"github.com/dimaglushkov/dkvs/internal/rpc"
	"google.golang.org/grpc"
)

type Server struct {
	rpc.StorageServer
	w Warehouse
}

func NewServer(w Warehouse) *Server {
	return &Server{w: w}
}

func (s Server) Get(ctx context.Context, in *rpc.Key, opts ...grpc.CallOption) (*rpc.Response, error) {
	var msg string
	v, err := s.w.Get(in.Key)
	if err != nil {
		msg = err.Error()
	}
	return &rpc.Response{Key: in.Key, Value: v, Msg: msg}, nil
}

func (s Server) Put(ctx context.Context, in *rpc.KeyValue, opts ...grpc.CallOption) (*rpc.Response, error) {
	var msg string

	err := s.w.Put(in.Key, in.Value)

	if err != nil {
		msg = err.Error()
	}
	return &rpc.Response{Key: in.Key, Value: in.Value, Msg: msg}, nil
}

func (s Server) Delete(ctx context.Context, in *rpc.Key, opts ...grpc.CallOption) (*rpc.Response, error) {
	var msg string

	err := s.w.Delete(in.Key)

	if err != nil {
		msg = err.Error()
	}
	return &rpc.Response{Key: in.Key, Msg: msg}, nil
}
