package svc

import (
	"context"
	"github.com/erfansahebi/lamia_auth/model"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type AuthDALInterface interface {
	StoreUser(ctx context.Context, user model.User) (storedUser model.User, err error)

	FetchUser(ctx context.Context, userID uuid.UUID) (fetchedUser model.User, err error)
	FetchUserByEmail(ctx context.Context, email string) (fetchedUser model.User, err error)

	StoreToken(ctx context.Context, tokenDetail model.Token, expireDuration uint) (tokenString string, err error)
	FetchToken(ctx context.Context, token string) (fetchedToken model.Token, err error)
}

type PgxConn interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type PgxRows interface {
	pgx.Rows
}

type PgxRow interface {
	pgx.Row
}
