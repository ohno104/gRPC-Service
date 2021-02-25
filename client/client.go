package main

import (
	"context"
	"flag"
	"io"
	"log"

	pb "felix.bs.com/felix/BeStrongerInGO/gRPC-Service/proto"
	"google.golang.org/grpc"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "通訊埠編號")
	flag.Parse()
}

func main() {
	conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	//_ = SayHello(client)
	_ = SayList(client)
}

func SayHello(client pb.GreeterClient) error {
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: "felix"})

	log.Printf("client.sayHello resop: %s", resp.Message)
	return nil
}

func SayList(client pb.GreeterClient) error {
	steam, _ := client.SayList(context.Background(), &pb.HelloRequest{Name: "felix"})
	for {
		resp, err := steam.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: %v", resp)
	}

	return nil
}
