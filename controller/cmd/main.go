package main

import (
	"context"
	"github.com/dimaglushkov/dkvs/api/storagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func run() error {
	conn, err := grpc.Dial("127.0.0.1:10031",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := storagepb.NewStorageClient(conn)
	client.Put(context.Background(), &storagepb.KeyValue{Key: "1", Value: "5"})

	res, err := client.Get(context.Background(), &storagepb.Key{Key: "1"})
	_ = res
	defer conn.Close()
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
