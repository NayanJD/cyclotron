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
	At  string `json:"at"`
	Err error  `json:"err"`
}

func MakeLoginEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		at, err := s.Login(ctx, req.Username, req.Password)
		return LoginResponse{
			At:  at,
			Err: err,
		}, nil
	}
}

func (r LoginResponse) Failed() error {
	return r.Err
}

type Failure interface {
	Failed() error
}

func (e Endpoints) Login(ctx context.Context, username string, password string) (at string, err error) {
	request := LoginRequest{
		Password: password,
		Username: username,
	}
	response, err := e.LoginEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(LoginResponse).At, response.(LoginResponse).Err
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
