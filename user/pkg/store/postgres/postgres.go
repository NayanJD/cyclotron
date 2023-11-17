package postgres

import (
	log "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type PostgresStore struct {
	db     *sqlx.DB
	logger log.Logger
}

func NewPostgresStore(connectionUrl string, logger log.Logger) (*PostgresStore, error) {

	// db, err := sqlx.Connect("postgres", connectionUrl)

	db, err := otelsqlx.Connect("postgres", connectionUrl, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))

	if err != nil {
		return nil, err
	}

	db.DB.SetMaxOpenConns(1000)
	db.DB.SetMaxIdleConns(10)   // defaultMaxIdleConns = 2
	db.DB.SetConnMaxLifetime(0) // 0, connections are reused forever.

	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db:     db,
		logger: logger,
	}, nil
}
