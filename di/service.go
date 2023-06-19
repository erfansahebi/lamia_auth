package di

import (
	"github.com/erfansahebi/lamia_auth/config"
	authProto "github.com/erfansahebi/lamia_shared/services/auth"
	"google.golang.org/grpc"
)

type AuthServiceInterface interface {
	Configuration() *config.Config

	Client() authProto.AuthServiceClient
}

type authService struct {
	configuration *config.Config
	client        authProto.AuthServiceClient
}

func NewAuthService(configuration *config.Config) AuthServiceInterface {
	return &authService{
		configuration: configuration,
	}
}

func (s *authService) Configuration() *config.Config {
	return s.configuration
}

func (s *authService) Client() authProto.AuthServiceClient {
	if err := s.initClient(); err != nil {
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
		return err
	}

	s.client = authProto.NewAuthServiceClient(cc)

	return nil
}
