package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"

	"github.com/dimaglushkov/dkvs/api/storagepb"
	"github.com/dimaglushkov/dkvs/storage/internal"
	"github.com/dimaglushkov/dkvs/storage/internal/warehouses"
)

func run(port int64) error {
	listener, err := net.Listen("tcp", ":"+strconv.FormatInt(port, 10))
	if err != nil {
		return fmt.Errorf("error while setting listener: %s", err)
	}

	storageHandler := internal.NewHandler(warehouses.NewHashTable())
	grpcServer := grpc.NewServer()
	storagepb.RegisterStorageServer(grpcServer, storageHandler)

	log.Printf("starting storage-handler-server listener on port %d\n", port)
	if err = grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("error while serving grpc server: %s", err)
	}

	return nil
}

func main() {
	var appPortVar = os.Getenv("APP_PORT")
	if appPortVar == "" {
		log.Fatal("APP_PORT env variable is missing")
	}

	appPort, err := strconv.ParseInt(appPortVar, 10, 64)
	if err != nil {
		log.Fatal("error while parsing APP_PORT env variable")
	}

	if err := run(appPort); err != nil {
		log.Fatal(err)
	}
}
