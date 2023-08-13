package service

import "context"

// UserService describes the service.
type UserService interface {
	// Add your methods here
	Login(ctx context.Context, username, password string) (at string, err error)
}

type basicUserService struct{}

func (b *basicUserService) Login(ctx context.Context, username string, password string) (at string, err error) {
	// TODO implement the business logic of Login
	return "abcd", err
}

// NewBasicUserService returns a naive, stateless implementation of UserService.
func NewBasicUserService() UserService {
	return &basicUserService{}
}

// New returns a UserService with all of the expected middleware wired in.
func New(middleware []Middleware) UserService {
	var svc UserService = NewBasicUserService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
