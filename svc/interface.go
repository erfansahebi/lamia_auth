package svc

import (
	"context"
	"github.com/erfansahebi/lamia_auth/model"
	"github.com/google/uuid"
)

type AuthDALInterface interface {
	StoreUser(ctx context.Context, user model.User) (storedUser model.User, err error)

	FetchUser(ctx context.Context, userID uuid.UUID) (fetchedUser model.User, err error)
	FetchUserByEmail(ctx context.Context, email string) (fetchedUser model.User, err error)

	StoreToken(ctx context.Context, tokenDetail model.Token, expireDuration uint) (tokenString string, err error)
	FetchToken(ctx context.Context, token string) (fetchedToken model.Token, err error)
	DeleteToken(ctx context.Context, token string)
}
