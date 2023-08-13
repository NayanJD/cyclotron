package grpc

import (
	"context"
	endpoint "user/pkg/endpoint"
	pb "user/pkg/grpc/pb"

	grpc "github.com/go-kit/kit/transport/grpc"
)

// makeLoginHandler creates the handler logic
func makeLoginHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.LoginEndpoint, decodeLoginRequest, encodeLoginResponse, options...)
}

// decodeLoginResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Login request.
func decodeLoginRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoginRequest)
	return endpoint.LoginRequest{Username: string(req.Username), Password: string(req.Password)}, nil
}

// encodeLoginResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
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
