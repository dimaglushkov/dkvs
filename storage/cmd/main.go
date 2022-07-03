package main

import (
	"flag"
	"fmt"
	"github.com/dimaglushkov/dkvs/api/storagepb"
	"github.com/dimaglushkov/dkvs/storage/internal"
	"github.com/dimaglushkov/dkvs/storage/internal/warehouses"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
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
	portFlag := flag.Int64("port", 0, "port number for storage to run on")
	flag.Parse()

	if *portFlag == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if err := run(*portFlag); err != nil {
		log.Fatal(err)
	}
}
