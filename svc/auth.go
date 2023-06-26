package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/erfansahebi/lamia_auth/model"
	"github.com/erfansahebi/lamia_shared/log"
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
	user.ID = uuid.New()

	row := a.pgx.QueryRow(
		ctx,
		`INSERT INTO users (
                	id,
	              	first_name,
	              	last_name,
	              	email,
	              	password
			) VALUES (
				  	$1, $2, $3, $4, $5
			) RETURNING id`,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	)

	if err := row.Scan(&user.ID); err != nil {
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

	return model.User{}, ErrUserDoesNotExists
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

	return model.User{}, ErrUserDoesNotExists
}

func (a *auth) StoreToken(ctx context.Context, tokenDetail model.Token, expireDuration uint) (tokenString string, err error) {
	tokenString = a.generateToken()

	tokenDetail.IssuedAt = time.Now()
	tokenDetail.ExpiredAt = tokenDetail.IssuedAt.Add(time.Duration(expireDuration) * time.Minute)

	data, err := json.Marshal(tokenDetail)
	if err != nil {
		log.WithError(err).Fatalf(ctx, "error in store data on redis %v", tokenDetail)
		return "", err
	}

	a.redis.Set(ctx, a.generateTokenKey(tokenString), data, time.Duration(expireDuration)*time.Minute)

	return tokenString, nil
}

func (a *auth) FetchToken(ctx context.Context, token string) (fetchedToken model.Token, err error) {
	fetchedData, err := a.redis.Get(ctx, a.generateTokenKey(token)).Result()
	if err != nil && err.Error() != "redis: nil" {
		log.WithError(err).Fatalf(ctx, "error in fetch data from redis")
		return fetchedToken, err
	}

	if fetchedData == "" {
		return fetchedToken, ErrEntryNotFound
	}

	if err = json.Unmarshal([]byte(fetchedData), &fetchedToken); err != nil {
		log.WithError(err).Fatalf(ctx, "error in unmarshal data from redis")
		return fetchedToken, err
	}

	return fetchedToken, nil
}

func (a *auth) DeleteToken(ctx context.Context, token string) {
	a.redis.Del(ctx, a.generateTokenKey(token))
}

func (a *auth) generateToken() string {
	return uuid.New().String()
}

func (a *auth) generateTokenKey(token string) string {
	return fmt.Sprintf("token.%s", token)
}
