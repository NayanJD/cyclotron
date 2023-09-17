package user

import (
	"context"
	"time"
	// customErrors "user/pkg/errors"
)

type (
	AuthToken struct {
		ID           int64  `json:"-"`
		RefreshToken string `json:"refreshToken"`
		AccessToken  string `json:"accessToken"`
		UserId       string `json:"-"`

		CreatedAt time.Time `json:"-"`
		ValidTill time.Time `json:"validTill"`
	}

	AuthTokenRepository interface {
		CreateToken(ctx context.Context, token *AuthToken, validityDurationInSeconds int64) (*AuthToken, error)
	}
)
