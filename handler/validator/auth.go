package validator

import (
	"context"
	"github.com/erfansahebi/lamia_auth/di"
	"github.com/erfansahebi/lamia_auth/model"
	"github.com/erfansahebi/lamia_auth/svc"
	authProto "github.com/erfansahebi/lamia_shared/services/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterStruct struct {
	*authProto.RegisterRequest
}

func (rs *RegisterStruct) Validate(ctx context.Context, di di.DIContainerInterface) (err error) {
	if _, err = di.AuthDAL().FetchUserByEmail(ctx, rs.User.Email); err == nil {
		return svc.ErrUserExists
	}

	return nil
}

type LoginStruct struct {
	FetchedUser model.User
	*authProto.LoginRequest
}

func (ls *LoginStruct) Validate(ctx context.Context, di di.DIContainerInterface) (err error) {
	ls.FetchedUser, err = di.AuthDAL().FetchUserByEmail(ctx, ls.Email)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(ls.FetchedUser.Password), []byte(ls.Password)); err != nil {
		return err
	}

	return nil
}

type AuthenticateStruct struct {
	TokenDetail model.Token
	*authProto.AuthenticateRequest
}

func (as *AuthenticateStruct) Validate(ctx context.Context, di di.DIContainerInterface) (err error) {
	as.TokenDetail, err = di.AuthDAL().FetchToken(ctx, as.Jwt)
	if err != nil {
		return err
	}

	return nil
}

type UserStruct struct {
	*authProto.GetUserRequest
	User model.User
}

func (us *UserStruct) Validate(ctx context.Context, di di.DIContainerInterface) (err error) {
	userID, err := uuid.Parse(us.UserId)
	if err != nil {
		return err
	}

	us.User, err = di.AuthDAL().FetchUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
