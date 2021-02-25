package main

import (
	"context"
	"flag"
	"net"

	pb "felix.bs.com/felix/BeStrongerInGO/gRPC-Service/proto"
	"google.golang.org/grpc"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "通訊埠編號")
	flag.Parse()
}

type GreeterServer struct{}

func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello.world"}, nil
}

func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	for n := 0; n < 6; n++ {
		_ = stream.Send(&pb.HelloReply{Message: "hello.list"})
	}

	return nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":"+port)
	server.Serve(lis)
}
