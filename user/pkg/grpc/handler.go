package grpc

import (
	"context"
	"errors"
	endpoint "user/pkg/endpoint"
	pb "user/pkg/grpc/pb"

	grpc "github.com/go-kit/kit/transport/grpc"
)

func makeLoginHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.LoginEndpoint, decodeLoginRequest, encodeLoginResponse, options...)
}

func decodeLoginRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoginRequest)
	return endpoint.LoginRequest{Username: string(req.Username), Password: string(req.Password)}, nil
}

func encodeLoginResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.LoginResponse)
	return &pb.LoginReply{At: string(resp.At), Err: err2str(resp.Err)}, nil
}
func (g *grpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	_, rep, err := g.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginReply), nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func makeRegisterHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.RegisterEndpoint, decodeRegisterRequest, encodeRegisterResponse, options...)
}

func decodeRegisterRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'User' Decoder is not impelemented")
}

func encodeRegisterResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'User' Encoder is not impelemented")
}
func (g *grpcServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	_, rep, err := g.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.RegisterReply), nil
}
