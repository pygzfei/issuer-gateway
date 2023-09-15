package driver

import (
	"cert-gateway/bus/grpc_server"
	"cert-gateway/bus/pb"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewGrpcServiceAndListen(addr string) {
	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		// 注册一元拦截器
		grpc.UnaryInterceptor(grpc_server.UnaryInterceptor),
		// 注册流拦截器
		grpc.StreamInterceptor(grpc_server.StreamInterceptor),
	)

	publisherServer := grpc_server.NewPublisherServer()

	// 监听 gRPC 服务端口
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 注册发布者服务到 gRPC 服务器
	pb.RegisterPubSubServiceServer(grpcServer, publisherServer)

	fmt.Printf("Server is listening on : %s", addr)
	// 启动 gRPC 服务器
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}