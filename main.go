package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/erfansahebi/lamia_auth/config"
	"github.com/erfansahebi/lamia_auth/database"
	"github.com/erfansahebi/lamia_auth/di"
	"github.com/erfansahebi/lamia_auth/handler"
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
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configurations.HTTP.Port))
	if err != nil {
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
	go func() {
		switch cmd {
		case "serve":
			authProto.RegisterAuthServiceServer(grpcServer, &h)

			if err = grpcServer.Serve(lis); err != nil {
				panic(err)
			}

			cancel()
		case "migrate":
			if err = database.Migrate(ctx, configurations, *migrateSteps); err != nil {
				fmt.Println("failed to migrate")
				panic(err)
			}

			cancel()
		case "makemigration":
			if err = database.MakeMigration(ctx, configurations, *migrateName); err != nil {
				fmt.Println("failed to make migration")
				panic(err)
			}

			cancel()
		default:
			panic(errors.New("wrong command"))
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
