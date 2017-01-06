package vinna

import (
	"net"

	"github.com/akresling/vinna/pb"
	"google.golang.org/grpc"
)

type server struct {
	ser pb.VinnaServer
}

// StartVinnaServer will create our grpc service
func StartVinnaServer() {
	listen, _ := net.Listen("tcp", ":3333")

	grpcServer := grpc.NewServer()

	service := New()
	pb.RegisterVinnaServer(grpcServer, service)
	grpcServer.Serve(listen)
}
