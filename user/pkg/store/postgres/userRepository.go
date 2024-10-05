package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	userDomain "cyclotron/user/pkg/domains/user"
)

type PostgresUser struct {
	ID             string    `db:"id"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	Username       string    `db:"username"`
	HashedPassword string    `db:"hashed_password"`
	Dob            time.Time `db:"dob"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func getDomainUser(pu *PostgresUser) *userDomain.User {
	return &userDomain.User{
		ID:             uuid.MustParse(pu.ID),
		FirstName:      pu.FirstName,
		LastName:       pu.LastName,
		Username:       pu.Username,
		HashedPassword: pu.HashedPassword,
		Dob:            pu.Dob,

		CreatedAt: pu.CreatedAt,
		UpdatedAt: pu.UpdatedAt,
		DeletedAt: pu.DeletedAt,
	}
}

func getPostgresUser(u *userDomain.User) *PostgresUser {
	return &PostgresUser{
		ID:             u.ID.URN(),
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
		Dob:            u.Dob,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}

func (ps *PostgresStore) CreateUser(
	ctx context.Context,
	user *userDomain.User,
) (*userDomain.User, error) {
	tx, err := ps.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	doesUserExistsWithUsername := false

	userExistsQuery := "SELECT EXISTS (SELECT 1 FROM users where username = $1 and deleted_at is null limit 1)"

	level.Debug(ps.logger).Log("query", userExistsQuery)

	if err := tx.GetContext(ctx, &doesUserExistsWithUsername, userExistsQuery, user.Username); err != nil {
		return nil, err
	}

	if doesUserExistsWithUsername {
		return nil, userDomain.UserAlreadyExistsErr
	}

	pu := getPostgresUser(user)

	userInsertQuery := `INSERT INTO users (
	id, first_name, last_name, username, hashed_password, dob, created_at, updated_at
	) VALUES (gen_random_uuid(), :first_name, :last_name, :username, :hashed_password, :dob, now(), now()) returning
    id, created_at, updated_at`

	level.Debug(ps.logger).Log("query", userInsertQuery)

	insertedRows, err := tx.NamedQuery(userInsertQuery, &pu)
	if err != nil {
		return nil, err
	}

	if insertedRows.Next() {
		insertedRows.Scan(&pu.ID, &pu.CreatedAt, &pu.UpdatedAt)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return getDomainUser(pu), nil
}

// func (ps *PostgresStore) FindByID(ctx *context.Context, user *user.User) (*user.User, error) {

// }

func (ps *PostgresStore) FindByUsername(
	ctx context.Context,
	username string,
) (*userDomain.User, error) {
	tx, err := ps.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var user *userDomain.User

	if user, err = ps.FindByUsernameTx(ctx, tx, username); err != nil {
		return nil, err
	} else {
		if err = tx.Commit(); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return user, err
}

func (ps *PostgresStore) FindByUsernameTx(
	ctx context.Context,
	tx *sqlx.Tx,
	username string,
) (*userDomain.User, error) {
	pu := PostgresUser{}

	getUserQuery := "SELECT id, first_name, last_name, username, hashed_password, dob, created_at, updated_at, deleted_at FROM users where username = $1"

	if err := tx.GetContext(ctx, &pu, getUserQuery, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.UserDoesNotExistsErr
		} else {
			return nil, err
		}
	}

	return getDomainUser(&pu), nil
}

func (ps *PostgresStore) FindByID(ctx context.Context, id string) (*userDomain.User, error) {
	tx, err := ps.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return ps.FindByIDTx(ctx, tx, id)
}

func (ps *PostgresStore) FindByIDTx(
	ctx context.Context,
	tx *sqlx.Tx,
	id string,
) (*userDomain.User, error) {
	pu := PostgresUser{}

	getUserQuery := "SELECT id, first_name, last_name, username, hashed_password, dob, created_at, updated_at, deleted_at FROM users where id = $1 and deleted_at is null"

	if err := tx.GetContext(ctx, &pu, getUserQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userDomain.UserDoesNotExistsErr
		} else {
			return nil, err
		}
	}

	return getDomainUser(&pu), nil
}

// func (ps *PostgresStore) UpdateUser(ctx *context.Context, user *user.User) (*user.User, error) {

// }

// func (ps *PostgresStore) DeleteUser(ctx *context.Context, user *user.User) error {

// }
