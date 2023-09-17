package user

import (
	"context"
	"github.com/google/uuid"
	"time"
	customErrors "user/pkg/errors"
)

type (
	User struct {
		ID             uuid.UUID `json:"id"`
		FirstName      string    `json:"firstName"`
		LastName       string    `json:"lastName"`
		Username       string    `json:"username"`
		Password       string    `json:"password"`
		HashedPassword string    `json:"-"`
		Dob            time.Time `json:"dob"`

		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt time.Time  `json:"updatedAt"`
		DeletedAt *time.Time `json:"deletedAt"`
	}

	UserRepository interface {
		CreateUser(ctx context.Context, user *User) (*User, error)

		FindByID(ctx context.Context, id string) (*User, error)

		FindByUsername(ctx context.Context, username string) (*User, error)

		// UpdateUser(ctx *context.Context, user *User) (*User, error)

		// DeleteUser(ctx *context.Context, user *User) error
	}
)

const (
	UserAlreadyExistsErr = customErrors.ConstError("The user with this username already exists")
	UserDoesNotExistsErr = customErrors.ConstError("The user does not exists")
)
