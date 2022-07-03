package internal

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
)

type serverAddr struct {
	ipAddr, port string
}

func (addr serverAddr) validate() bool {
	if ip := net.ParseIP(addr.ipAddr); ip == nil {
		return false
	}
	if _, err := strconv.ParseInt(addr.port, 10, 32); err != nil {
		return false
	}
	return true
}

//TODO: replace insecure.NewCredentials() with the proper authentication
func grpcConnector(addr, port string) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(addr+":"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
