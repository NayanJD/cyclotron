package user

import (
	"github.com/google/uuid"
	"time"
)

type (
	User struct {
		ID             uuid.UUID `json:"id"`
		FirstName      string    `json:"firstName"`
		LastName       string    `json:"LastName"`
		Username       string    `json:"username"`
		HashedPassword string    `json:"-"`
		Dob            time.Time `json:"dob"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		DeletedAt time.Time `json:"deletedAt"`
	}

	AuthToken struct {
		ID           int64  `json:"id"`
		RefreshToken string `json:"refreshToken"`
		AccessToken  string `json:"accessToken"`
		UserId       string `json:"userId"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	UserRepository interface {
		CreateUser(user *User) (*User, error)

		FindByID(id uuid.UUID) (*User, error)
		FindByUsername(username string) (*User, error)

		UpdateUser(user *User) (*User, error)

		DeleteUser(user *User) error
	}
)
