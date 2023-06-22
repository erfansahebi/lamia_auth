package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/erfansahebi/lamia_auth/config"
	"github.com/erfansahebi/lamia_auth/database"
	"github.com/erfansahebi/lamia_auth/di"
	"github.com/erfansahebi/lamia_auth/handler"
	sharedCommon "github.com/erfansahebi/lamia_shared/common"
	"github.com/erfansahebi/lamia_shared/log"
	authProto "github.com/erfansahebi/lamia_shared/services/auth"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	configurations, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatalf(ctx, "failed to load config")
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configurations.HTTP.Port))
	if err != nil {
		log.WithError(err).Fatalf(ctx, "failed to listen tcp")
		panic(err)
	}

	grpcServer := grpc.NewServer()

	h := handler.Handler{
		AppCtx: ctx,
		Di:     di.NewDIContainer(ctx, configurations),
	}

	migrateSteps := flag.Int("migrate", 0, "number of steps to migrate")
	migrateName := flag.String("mname", "", "migration name")
	flag.Parse()

	cmd := flag.Arg(0)

	log.Infof(ctx, "Starting with command: %s", cmd)

	go func() {
		switch cmd {
		case "serve":
			authProto.RegisterAuthServiceServer(grpcServer, &h)

			if err = grpcServer.Serve(lis); err != nil {
				log.WithError(err).Fatalf(ctx, "failed to serve grpc server")
				panic(err)
			}

			cancel()
		case "migrate":
			if err = database.Migrate(ctx, configurations, *migrateSteps); err != nil {
				log.WithError(err).Fatalf(ctx, "failed to migrate")
				panic(err)
			}

			cancel()
		case "makemigration":
			if err = database.MakeMigration(ctx, configurations, *migrateName); err != nil {
				log.WithError(err).Fatalf(ctx, "failed to make migration")
				panic(err)
			}

			cancel()
		default:
			log.WithError(sharedCommon.ErrWrongCommand).Fatalf(ctx, "wrong command")
			panic(sharedCommon.ErrWrongCommand)
		}

	}()

	sig := make(chan os.Signal, 10)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-sig:
		fmt.Printf("Received %s, graceful shut down...", s.String())
		cancel()
	case <-ctx.Done():
	}

	time.Sleep(1 * time.Second)
}
