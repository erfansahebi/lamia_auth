package di

import (
	"context"
	"github.com/erfansahebi/lamia_auth/config"
	"github.com/erfansahebi/lamia_auth/svc"
	"github.com/erfansahebi/lamia_shared/log"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DIContainerInterface interface {
	Config() *config.Config

	AuthDAL() svc.AuthDALInterface

	Service() AuthServiceInterface
}

type diContainer struct {
	ctx           context.Context
	configuration *config.Config

	authDAL svc.AuthDALInterface

	service AuthServiceInterface

	pgx   *pgxpool.Pool
	redis *redis.Client
}

func NewDIContainer(ctx context.Context, config *config.Config) DIContainerInterface {
	return &diContainer{
		ctx:           ctx,
		configuration: config,
	}
}

func (d *diContainer) Config() *config.Config {
	return d.configuration
}

func (d *diContainer) Service() AuthServiceInterface {
	if err := d.initService(); err != nil {
		log.WithError(err).Fatalf(d.ctx, "error in init service")
		panic(err)
	}

	return d.service
}

func (d *diContainer) initService() error {
	if d.service != nil {
		return nil
	}

	d.service = NewAuthService(d.ctx, d.configuration)

	return nil
}

func (d *diContainer) getPgxConnection(dbName string) (*pgxpool.Pool, error) {
	if d.pgx != nil {
		return d.pgx, nil
	}

	pgx, err := pgxpool.Connect(d.ctx, d.Config().GetDbUrl(dbName))
	if err != nil {
		log.WithError(err).Fatalf(d.ctx, "error in pgxpool connection")
		return nil, err
	}

	d.pgx = pgx

	return d.pgx, nil
}

func (d *diContainer) AuthDAL() svc.AuthDALInterface {
	if err := d.initAuthDAL(); err != nil {
		log.WithError(err).Fatalf(d.ctx, "error in init auth dal")
		panic(err)
	}

	return d.authDAL
}

func (d *diContainer) initAuthDAL() error {
	if d.authDAL != nil {
		return nil
	}

	pgxConn, err := d.getPgxConnection("app")
	if err != nil {
		log.WithError(err).Fatalf(d.ctx, "error in pgx connection")
		return err
	}

	d.authDAL = svc.NewAuthDAL(pgxConn, d.getRedisClient())

	return nil
}

func (d *diContainer) getRedisClient() *redis.Client {
	if err := d.initRedisClient(); err != nil {
		log.WithError(err).Fatalf(d.ctx, "error in init redis client")
		panic(err)
	}

	return d.redis
}

func (d *diContainer) initRedisClient() error {
	if d.redis != nil {
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     d.configuration.Redis.Host + ":" + d.configuration.Redis.Port,
		Username: d.configuration.Redis.Username,
		Password: d.configuration.Redis.Password,
		DB:       d.configuration.Redis.DB,
	})

	if _, err := client.Ping(d.ctx).Result(); err != nil {
		log.WithError(err).Fatalf(d.ctx, "error in redis connection")
		return err
	}

	d.redis = client

	return nil
}
