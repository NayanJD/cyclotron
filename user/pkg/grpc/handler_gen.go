// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package grpc

import (
	grpc "github.com/go-kit/kit/transport/grpc"
	endpoint "user/pkg/endpoint"
	pb "user/pkg/grpc/pb"
)

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer
type grpcServer struct {
	login grpc.Handler
	pb.UnimplementedUserServer
}

func NewGRPCServer(endpoints endpoint.Endpoints, options map[string][]grpc.ServerOption) pb.UserServer {
	return &grpcServer{login: makeLoginHandler(endpoints, options["Login"])}
}
