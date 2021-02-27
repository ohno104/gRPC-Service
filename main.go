package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	pb "felix.bs.com/felix/BeStrongerInGO/gRPC-Service/proto"
	"felix.bs.com/felix/BeStrongerInGO/gRPC-Service/server"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8001", "通訊埠編號")
	flag.Parse()
}

func main() {
	l, err := RunTCPServer(port)
	if err != nil {
		log.Fatalf("Run TCP Server err: %v", err)
	}

	m := cmux.New(l)
	grpcL := m.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldPrefixSendSettings(
			"content-type",
			"application/grpc",
		),
	)
	httpL := m.Match(cmux.HTTP1Fast())

	grpcS := RunGrpcServer()
	httpS := RunHttpServer(port)
	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	err = m.Serve()
	if err != nil {
		log.Fatalf("Run Serve err: %v", err)
	}
}

func RunHttpServer(port string) *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}

func RunGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{
		//grpc.UnaryInterceptor(HelloInterceptor),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			HelloInterceptor,
			WorldInterceptor,
			AccessLog,
			ErrorLog,
		)),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterTagServiceServer(s, server.NewTagServer())

	//using for grpcurl to debug
	reflection.Register(s)
	return s
}

func RunTCPServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好~")
	resp, err := handler(ctx, req)
	log.Println("再見~")
	return resp, err
}

func WorldInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好 WORLD~")
	resp, err := handler(ctx, req)
	log.Println("再見 WORLD~")
	return resp, err
}

func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestLog := "access request log: method: %s, begin_time: %d. request: %v"
	beginTime := time.Now().Local().Unix()
	log.Printf(requestLog, info.FullMethod, beginTime, req)

	resp, err := handler(ctx, req)
	responseLog := "access response Log: method: %s, begin_time: %d, end_time: %d, response: %v"
	endTime := time.Now().Local().Unix()
	log.Printf(responseLog, info.FullMethod, beginTime, endTime, resp)
	return resp, err
}

func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("error log")
	}
	return resp, err
}
