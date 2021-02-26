package main

import (
	"context"
	"log"

	pb "felix.bs.com/felix/BeStrongerInGO/gRPC-Service/proto"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	clientConn, _ := GetClientConn(ctx, "localhost:8001", nil)
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, _ := tagServiceClient.GetTagList(
		ctx,
		&pb.GetTagListRequest{
			Name: "",
		},
	)

	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}

/*
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
	//_ = SayList(client)
	//_ = SayRecord(client, &pb.HelloRequest{Name: "felix"})
	_ = SayRoute(client, &pb.HelloRequest{Name: "felix"})
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

func SayRecord(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayRecord(context.Background())
	for n := 0; n < 6; n++ {
		_ = stream.Send(r)
	}
	resp, _ := stream.CloseAndRecv()

	log.Printf("resp err: %v", resp)
	return nil
}

func SayRoute(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayRoute(context.Background())
	for n := 0; n <= 6; n++ {
		_ = stream.Send(r)

		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp err: %v", resp)
	}

	_ = stream.CloseSend()
	return nil
}
*/
