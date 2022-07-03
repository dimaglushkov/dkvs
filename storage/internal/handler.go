package internal

import (
	"context"
	"github.com/dimaglushkov/dkvs/api/storagepb"
)

type handler struct {
	storagepb.StorageServer
	w Warehouse
}

func NewHandler(w Warehouse) *handler {
	return &handler{w: w}
}

func (s handler) Get(ctx context.Context, in *storagepb.Key) (*storagepb.Response, error) {
	v, err := s.w.Get(in.Key)
	if err != nil {
		return &storagepb.Response{Success: false, Value: err.Error()}, nil
	}
	return &storagepb.Response{Success: true, Value: v}, nil
}

func (s handler) Put(ctx context.Context, in *storagepb.KeyValue) (*storagepb.Response, error) {
	if err := s.w.Put(in.Key, in.Value); err != nil {
		return &storagepb.Response{Success: false, Value: err.Error()}, nil
	}
	return &storagepb.Response{Success: true}, nil
}

func (s handler) Delete(ctx context.Context, in *storagepb.Key) (*storagepb.Response, error) {
	if err := s.w.Delete(in.Key); err != nil {
		return &storagepb.Response{Success: false, Value: err.Error()}, nil
	}
	return &storagepb.Response{Success: true}, nil
}
