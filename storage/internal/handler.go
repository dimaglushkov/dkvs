package storage

import (
	context "context"
	"github.com/dimaglushkov/dkvs/internal/rpc"
)

type handler struct {
	rpc.StorageServer
	w Warehouse
}

func NewHandler(w Warehouse) *handler {
	return &handler{w: w}
}

func (s handler) Get(ctx context.Context, in *rpc.Key) (*rpc.Response, error) {
	v, err := s.w.Get(in.Key)
	if err != nil {
		return &rpc.Response{Success: false, Value: err.Error()}, nil
	}
	return &rpc.Response{Success: true, Value: v}, nil
}

func (s handler) Put(ctx context.Context, in *rpc.KeyValue) (*rpc.Response, error) {
	if err := s.w.Put(in.Key, in.Value); err != nil {
		return &rpc.Response{Success: false, Value: err.Error()}, nil
	}
	return &rpc.Response{Success: true}, nil
}

func (s handler) Delete(ctx context.Context, in *rpc.Key) (*rpc.Response, error) {
	if err := s.w.Delete(in.Key); err != nil {
		return &rpc.Response{Success: false, Value: err.Error()}, nil
	}
	return &rpc.Response{Success: true}, nil
}
