package endpoint

import (
	"context"
	user "user/pkg/domains/user"
	service "user/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token user.AuthToken `json:"token"`
	Err   error          `json:"err"`
}

func MakeLoginEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		token, err := s.Login(ctx, req.Username, req.Password)
		return LoginResponse{
			Token: token,
			Err:   err,
		}, nil
	}
}

func (r LoginResponse) Failed() error {
	return r.Err
}

type Failure interface {
	Failed() error
}

func (e Endpoints) Login(ctx context.Context, username string, password string) (token user.AuthToken, err error) {
	request := LoginRequest{
		Password: password,
		Username: username,
	}
	response, err := e.LoginEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(LoginResponse).Token, response.(LoginResponse).Err
}

type RegisterRequest struct {
	User user.User `json:"user"`
}

type RegisterResponse struct {
	NewUser user.User `json:"new_user"`
	Err     error     `json:"err"`
}

func MakeRegisterEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterRequest)
		newUser, err := s.Register(ctx, req.User)
		return RegisterResponse{
			Err:     err,
			NewUser: newUser,
		}, nil
	}
}

func (r RegisterResponse) Failed() error {
	return r.Err
}

func (e Endpoints) Register(ctx context.Context, user user.User) (newUser user.User, err error) {
	request := RegisterRequest{User: user}
	response, err := e.RegisterEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(RegisterResponse).NewUser, response.(RegisterResponse).Err
}

// GetUserFromTokenRequest collects the request parameters for the GetUserFromToken method.
type GetUserFromTokenRequest struct {
	Token string `json:"token"`
}

// GetUserFromTokenResponse collects the response parameters for the GetUserFromToken method.
type GetUserFromTokenResponse struct {
	User user.User `json:"user"`
	Err  error     `json:"err"`
}

// MakeGetUserFromTokenEndpoint returns an endpoint that invokes GetUserFromToken on the service.
func MakeGetUserFromTokenEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserFromTokenRequest)
		user, err := s.GetUserFromToken(ctx, req.Token)
		return GetUserFromTokenResponse{
			Err:  err,
			User: user,
		}, nil
	}
}

// Failed implements Failer.
func (r GetUserFromTokenResponse) Failed() error {
	return r.Err
}

// GetUserFromToken implements Service. Primarily useful in a client.
func (e Endpoints) GetUserFromToken(ctx context.Context, token string) (user user.User, err error) {
	request := GetUserFromTokenRequest{Token: token}
	response, err := e.GetUserFromTokenEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetUserFromTokenResponse).User, response.(GetUserFromTokenResponse).Err
}

// RefreshAccessTokenRequest collects the request parameters for the RefreshAccessToken method.
type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshAccessTokenResponse collects the response parameters for the RefreshAccessToken method.
type RefreshAccessTokenResponse struct {
	Token user.AuthToken `json:"token"`
	Err   error          `json:"err"`
}

// MakeRefreshAccessTokenEndpoint returns an endpoint that invokes RefreshAccessToken on the service.
func MakeRefreshAccessTokenEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RefreshAccessTokenRequest)
		token, err := s.RefreshAccessToken(ctx, req.RefreshToken)
		return RefreshAccessTokenResponse{
			Err:   err,
			Token: token,
		}, nil
	}
}

// Failed implements Failer.
func (r RefreshAccessTokenResponse) Failed() error {
	return r.Err
}

// RefreshAccessToken implements Service. Primarily useful in a client.
func (e Endpoints) RefreshAccessToken(ctx context.Context, refreshToken string) (token user.AuthToken, err error) {
	request := RefreshAccessTokenRequest{RefreshToken: refreshToken}
	response, err := e.RefreshAccessTokenEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(RefreshAccessTokenResponse).Token, response.(RefreshAccessTokenResponse).Err
}
