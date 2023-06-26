package handler

import (
	"context"
	"github.com/erfansahebi/lamia_auth/handler/validator"
	"github.com/erfansahebi/lamia_auth/model"
	authProto "github.com/erfansahebi/lamia_shared/services/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (h *Handler) Register(ctx context.Context, request *authProto.RegisterRequest) (*authProto.AuthenticationResponse, error) {
	pendData := validator.RegisterStruct{RegisterRequest: request}

	if err := pendData.Validate(ctx, h.Di); err != nil {
		return nil, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pendData.User.Password), bcrypt.MinCost)

	registeredUser, err := h.Di.AuthDAL().StoreUser(ctx, model.User{
		ID:        uuid.Nil,
		FirstName: pendData.User.FirstName,
		LastName:  pendData.User.LastName,
		Email:     pendData.User.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	})
	if err != nil {
		return nil, err
	}

	tokenString, err := h.Di.AuthDAL().StoreToken(ctx, model.Token{
		UserID:    registeredUser.ID,
		IssuedAt:  time.Time{},
		ExpiredAt: time.Time{},
	}, h.Di.Config().AuthorizationToken.Duration)
	if err != nil {
		return nil, err
	}

	return &authProto.AuthenticationResponse{
		User: &authProto.UserStruct{
			Id:        registeredUser.ID.String(),
			FirstName: registeredUser.FirstName,
			LastName:  registeredUser.LastName,
			Email:     registeredUser.Email,
			Password:  registeredUser.Password,
		},
		AuthorizationToken: tokenString,
	}, nil
}

func (h *Handler) Login(ctx context.Context, request *authProto.LoginRequest) (*authProto.AuthenticationResponse, error) {
	pendData := validator.LoginStruct{LoginRequest: request}
	if err := pendData.Validate(ctx, h.Di); err != nil {
		return nil, err
	}

	tokenString, err := h.Di.AuthDAL().StoreToken(ctx, model.Token{
		UserID:    pendData.FetchedUser.ID,
		IssuedAt:  time.Time{},
		ExpiredAt: time.Time{},
	}, h.Di.Config().AuthorizationToken.Duration)
	if err != nil {
		return nil, err
	}

	return &authProto.AuthenticationResponse{
		User: &authProto.UserStruct{
			Id:        pendData.FetchedUser.ID.String(),
			FirstName: pendData.FetchedUser.FirstName,
			LastName:  pendData.FetchedUser.LastName,
			Email:     pendData.FetchedUser.Email,
			Password:  pendData.FetchedUser.Password,
		},
		AuthorizationToken: tokenString,
	}, nil
}

func (h *Handler) Authenticate(ctx context.Context, request *authProto.AuthenticateRequest) (*authProto.AuthenticateResponse, error) {
	pendData := validator.AuthenticateStruct{AuthenticateRequest: request}
	if err := pendData.Validate(ctx, h.Di); err != nil {
		return nil, err
	}

	return &authProto.AuthenticateResponse{
		Id: pendData.TokenDetail.UserID.String(),
	}, nil
}

func (h *Handler) GetUser(ctx context.Context, request *authProto.GetUserRequest) (*authProto.GetUserResponse, error) {
	pendData := validator.UserStruct{GetUserRequest: request}
	if err := pendData.Validate(ctx, h.Di); err != nil {
		return nil, err
	}

	return &authProto.GetUserResponse{
		User: &authProto.UserStruct{
			Id:        pendData.User.ID.String(),
			FirstName: pendData.User.FirstName,
			LastName:  pendData.User.LastName,
			Email:     pendData.User.Email,
			Password:  pendData.User.Password,
		},
	}, nil
}
