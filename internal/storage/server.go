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
	v, err := s.w.Get(in.Key)
	if err != nil {
		return &rpc.Response{Success: false, Value: err.Error()}, nil
	}
	return &rpc.Response{Success: true, Value: v}, nil
}

func (s Server) Put(ctx context.Context, in *rpc.KeyValue, opts ...grpc.CallOption) (*rpc.Response, error) {
	if err := s.w.Put(in.Key, in.Value); err != nil {
		return &rpc.Response{Success: false, Value: err.Error()}, nil
	}
	return &rpc.Response{Success: true}, nil
}

func (s Server) Delete(ctx context.Context, in *rpc.Key, opts ...grpc.CallOption) (*rpc.Response, error) {
	if err := s.w.Delete(in.Key); err != nil {
		return &rpc.Response{Success: false, Value: err.Error()}, nil
	}
	return &rpc.Response{Success: true}, nil
}
