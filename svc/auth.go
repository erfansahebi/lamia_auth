package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/erfansahebi/lamia_auth/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type auth struct {
	pgx   PgxConn
	redis *redis.Client
}

func NewAuthDAL(pgx PgxConn, redis *redis.Client) AuthDALInterface {
	return &auth{
		pgx:   pgx,
		redis: redis,
	}
}

func (a *auth) StoreUser(ctx context.Context, user model.User) (model.User, error) {
	tx, err := a.pgx.Begin(ctx)
	if err != nil {
		return model.User{}, err
	}

	defer tx.Rollback(ctx)

	row := tx.QueryRow(
		ctx,
		`INSERT INTO users (
                   first_name,
                   last_name,
                   email,
                   password
			) VALUES (
				  $1, $2, $3, $4
			) RETURNING id`,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	)

	if err = row.Scan(&user.ID); err != nil {
		return model.User{}, err
	} else if err = tx.Commit(ctx); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (a *auth) FetchUser(ctx context.Context, userID uuid.UUID) (fetchedUser model.User, err error) {
	row, err := a.pgx.Query(
		ctx,
		`SELECT id,
					first_name,
					last_name,
					email,
					password,
					created_at,
					updated_at
			FROM users
			WHERE id = $1`,
		userID,
	)
	if err != nil {
		return model.User{}, err
	}

	defer row.Close()

	if row.Next() {

		fetchedUser, err = model.ScanToUser(row.Scan)
		if err != nil {
			return model.User{}, err
		}

		return fetchedUser, nil
	}

	return model.User{}, ErrEntryNotFound
}

func (a *auth) FetchUserByEmail(ctx context.Context, email string) (fetchedUser model.User, err error) {
	row, err := a.pgx.Query(
		ctx,
		`SELECT id,
					first_name,
					last_name,
					email,
					password,
					created_at,
					updated_at
			FROM users
			WHERE email = $1`,
		email,
	)
	if err != nil {
		return model.User{}, err
	}

	defer row.Close()

	if row.Next() {

		fetchedUser, err = model.ScanToUser(row.Scan)
		if err != nil {
			return model.User{}, err
		}

		return fetchedUser, nil
	}

	return model.User{}, ErrEntryNotFound
}

func (a *auth) StoreToken(ctx context.Context, tokenDetail model.Token, expireDuration uint) (tokenString string, err error) {
	tokenString = a.generateToken()

	tokenDetail.IssuedAt = time.Now()
	tokenDetail.ExpiredAt = time.Now().Add(time.Duration(expireDuration) * time.Minute)

	a.redis.HSet(ctx, a.generateTokenKey(tokenString), tokenDetail, time.Duration(expireDuration)*time.Minute)

	return tokenString, nil
}

func (a *auth) FetchToken(ctx context.Context, token string) (fetchedToken model.Token, err error) {
	fetchedData, err := a.redis.Get(ctx, a.generateTokenKey(token)).Result()
	if err != nil && err.Error() != "redis: nil" {
		return fetchedToken, err
	}

	if fetchedData == "" {
		return fetchedToken, ErrEntryNotFound
	}

	if err = json.Unmarshal([]byte(fetchedData), &fetchedToken); err != nil {
		return fetchedToken, err
	}

	return fetchedToken, nil
}

func (a *auth) generateToken() string {
	return uuid.New().String()
}

func (a *auth) generateTokenKey(token string) string {
	return fmt.Sprintf("token.%s", token)
}
