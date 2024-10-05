package grpc

import (
	"context"

	grpc "github.com/go-kit/kit/transport/grpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	userDomain "cyclotron/user/pkg/domains/user"
	endpoint "cyclotron/user/pkg/endpoint"
	pb "cyclotron/user/pkg/grpc/pb"
)

func makeLoginHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeLoginResponse,
		options...)
}

func decodeLoginRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoginRequest)
	return endpoint.LoginRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

func encodeLoginResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.LoginResponse)

	if resp.Err != nil {
		return nil, resp.Err
	}

	return &pb.LoginReply{
		AccessToken:  resp.Token.AccessToken,
		RefreshToken: resp.Token.RefreshToken,
		ValidTill:    timestamppb.New(resp.Token.ValidTill),
	}, nil
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
	return grpc.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeRegisterResponse,
		options...)
}

func decodeRegisterRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.RegisterRequest)

	return endpoint.RegisterRequest{
		User: userDomain.User{
			FirstName: string(req.FirstName),
			LastName:  string(req.LastName),
			Username:  string(req.Username),
			Password:  string(req.Password),
			Dob:       req.Dob.AsTime(),
		},
	}, nil
}

func encodeRegisterResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.RegisterResponse)

	if resp.Err != nil {
		return nil, resp.Err
	}

	var deletedAt *timestamppb.Timestamp

	if resp.NewUser.DeletedAt != nil {
		deletedAt = timestamppb.New(*resp.NewUser.DeletedAt)
	}
	return &pb.RegisterReply{
		Id:        resp.NewUser.ID.String(),
		FirstName: resp.NewUser.FirstName,
		LastName:  resp.NewUser.LastName,
		Username:  resp.NewUser.Username,
		Dob:       timestamppb.New(resp.NewUser.Dob),
		CreatedAt: timestamppb.New(resp.NewUser.CreatedAt),
		DeletedAt: deletedAt,
	}, nil
}

func (g *grpcServer) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.RegisterReply, error) {
	_, rep, err := g.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.RegisterReply), nil
}

func makeGetUserFromTokenHandler(
	endpoints endpoint.Endpoints,
	options []grpc.ServerOption,
) grpc.Handler {
	return grpc.NewServer(
		endpoints.GetUserFromTokenEndpoint,
		decodeGetUserFromTokenRequest,
		encodeGetUserFromTokenResponse,
		options...)
}

func decodeGetUserFromTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetUserFromTokenRequest)
	return endpoint.GetUserFromTokenRequest{Token: req.AccessToken}, nil
}

func encodeGetUserFromTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.GetUserFromTokenResponse)

	if resp.Err != nil {
		return nil, resp.Err
	}

	var deletedAt *timestamppb.Timestamp

	if resp.User.DeletedAt != nil {
		deletedAt = timestamppb.New(*resp.User.DeletedAt)
	}

	return &pb.GetUserFromTokenReply{
		Id:        resp.User.ID.String(),
		FirstName: resp.User.FirstName,
		LastName:  resp.User.LastName,
		Username:  resp.User.Username,
		Dob:       timestamppb.New(resp.User.Dob),
		CreatedAt: timestamppb.New(resp.User.CreatedAt),
		DeletedAt: deletedAt,
	}, nil
}

func (g *grpcServer) GetUserFromToken(
	ctx context.Context,
	req *pb.GetUserFromTokenRequest,
) (*pb.GetUserFromTokenReply, error) {
	_, rep, err := g.getUserFromToken.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetUserFromTokenReply), nil
}

func makeRefreshAccessTokenHandler(
	endpoints endpoint.Endpoints,
	options []grpc.ServerOption,
) grpc.Handler {
	return grpc.NewServer(
		endpoints.RefreshAccessTokenEndpoint,
		decodeRefreshAccessTokenRequest,
		encodeRefreshAccessTokenResponse,
		options...)
}

func decodeRefreshAccessTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.RefreshAccessTokenRequest)
	return endpoint.RefreshAccessTokenRequest{RefreshToken: req.RefreshToken}, nil
}

func encodeRefreshAccessTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.RefreshAccessTokenResponse)

	if resp.Err != nil {
		return nil, resp.Err
	}

	return &pb.RefreshAccessTokenReply{
		AccessToken:  resp.Token.AccessToken,
		RefreshToken: resp.Token.RefreshToken,
		ValidTill:    timestamppb.New(resp.Token.ValidTill),
	}, nil
}

func (g *grpcServer) RefreshAccessToken(
	ctx context.Context,
	req *pb.RefreshAccessTokenRequest,
) (*pb.RefreshAccessTokenReply, error) {
	_, rep, err := g.refreshAccessToken.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.RefreshAccessTokenReply), nil
}
