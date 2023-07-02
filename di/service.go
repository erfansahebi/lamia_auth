package di

import (
	"context"
	"github.com/erfansahebi/lamia_auth/config"
	"github.com/erfansahebi/lamia_shared/go/log"
	authProto "github.com/erfansahebi/lamia_shared/go/proto/auth"
	"google.golang.org/grpc"
)

type AuthServiceInterface interface {
	Configuration() *config.Config

	Client() authProto.AuthServiceClient
}

type authService struct {
	ctx           context.Context
	configuration *config.Config
	client        authProto.AuthServiceClient
}

func NewAuthService(ctx context.Context, configuration *config.Config) AuthServiceInterface {
	return &authService{
		ctx:           ctx,
		configuration: configuration,
	}
}

func (s *authService) Configuration() *config.Config {
	return s.configuration
}

func (s *authService) Client() authProto.AuthServiceClient {
	if err := s.initClient(); err != nil {
		log.WithError(err).Fatalf(s.ctx, "error in init client service")
		panic(err)
	}

	return s.client
}

func (s *authService) initClient() error {
	if s.client != nil {
		return nil
	}

	cc, err := grpc.Dial(s.Configuration().GetHttpUrl(), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf(s.ctx, "error in crete client connection")
		return err
	}

	s.client = authProto.NewAuthServiceClient(cc)

	return nil
}
