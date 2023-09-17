package postgres

import (
	"context"
	// "database/sql"
	// "errors"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/jmoiron/sqlx"
	"time"
	userDomain "user/pkg/domains/user"
)

type PostgresAuthToken struct {
	ID           int64  `db:"id"`
	RefreshToken string `db:"refresh_token"`
	AccessToken  string `db:"access_token"`
	UserId       string `db:"user_id"`

	CreatedAt time.Time `db:"created_at"`
	ValidTill time.Time `db:"valid_till"`
}

func getDomainAuthToken(at *PostgresAuthToken) *userDomain.AuthToken {
	return &userDomain.AuthToken{
		ID:           at.ID,
		RefreshToken: at.RefreshToken,
		AccessToken:  at.AccessToken,
		UserId:       at.UserId,

		CreatedAt: at.CreatedAt,
		ValidTill: at.ValidTill,
	}
}

func getPostgresAuthToken(at *userDomain.AuthToken) *PostgresAuthToken {
	return &PostgresAuthToken{
		ID:           at.ID,
		RefreshToken: at.RefreshToken,
		AccessToken:  at.AccessToken,
		UserId:       at.UserId,

		CreatedAt: at.CreatedAt,
		ValidTill: at.ValidTill,
	}
}

func (ps *PostgresStore) CreateToken(ctx context.Context, token *userDomain.AuthToken, validityDurationInSeconds int64) (*userDomain.AuthToken, error) {
	tx, err := ps.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	var newAuthToken *userDomain.AuthToken

	if newAuthToken, err = ps.CreateAuthTokenTx(ctx, tx, token, validityDurationInSeconds); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return newAuthToken, nil
}

func (ps *PostgresStore) CreateAuthTokenTx(ctx context.Context, tx *sqlx.Tx, token *userDomain.AuthToken, validityDurationInSeconds int64) (*userDomain.AuthToken, error) {

	pat := getPostgresAuthToken(token)

	tokenInsertQuery := fmt.Sprintf(`INSERT INTO tokens(
        access_token, refresh_token, user_id, created_at, valid_till
    ) VALUES (:access_token, :refresh_token, :user_id, now(), now() + interval '%d seconds')
    returning id, created_at, valid_till`, validityDurationInSeconds)

	level.Debug(ps.logger).Log("query", tokenInsertQuery)

	insertedRows, err := tx.NamedQuery(tokenInsertQuery, &pat)

	if err != nil {
		return nil, err
	}

	if insertedRows.Next() {
		insertedRows.Scan(&pat.ID, &pat.CreatedAt, &pat.ValidTill)
	}

	return getDomainAuthToken(pat), nil

}
