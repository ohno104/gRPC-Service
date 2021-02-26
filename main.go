package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "felix.bs.com/felix/BeStrongerInGO/gRPC-Service/proto"
	"felix.bs.com/felix/BeStrongerInGO/gRPC-Service/server"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8001", "通訊埠編號")
	flag.Parse()
}

func main() {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())

	//using for grpcurl to debug
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("server.Serve err: %v", err)
	}
}
