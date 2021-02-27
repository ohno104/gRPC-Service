package main

import (
	"context"
	"log"
	"time"

	pb "felix.bs.com/felix/BeStrongerInGO/gRPC-Service/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func main() {
	ctx := context.Background()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			UnaryContextTimeout(),

			//retry
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithMax(2),
				grpc_retry.WithCodes(
					codes.Unknown,
					codes.Internal,
					codes.DeadlineExceeded,
				),
			),
		),
	))

	opts = append(opts, grpc.WithStreamInterceptor(
		grpc_middleware.ChainStreamClient(
			StreamContextTimeout(),

			//retry
			grpc_retry.StreamClientInterceptor(
				grpc_retry.WithMax(2),
				grpc_retry.WithCodes(
					codes.Unknown,
					codes.Internal,
					codes.DeadlineExceeded,
				),
			),
		),
	))

	clientConn, err := GetClientConn(ctx, "localhost:8001", opts)

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	defer clientConn.Close()
	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(
		ctx,
		&pb.GetTagListRequest{
			Name: "",
		},
	)
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}

func UnaryContextTimeout() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := defaultContextTimeout(ctx)
		if cancel != nil {
			defer cancel()
		}
		return invoker(ctx, method, req, resp, cc, opts...)
	}
}

func StreamContextTimeout() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx, cancel := defaultContextTimeout(ctx)
		if cancel != nil {
			defer cancel()
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		defaultTime := 3 * time.Second
		ctx, cancel = context.WithTimeout(ctx, defaultTime)
	}

	return ctx, cancel
}
