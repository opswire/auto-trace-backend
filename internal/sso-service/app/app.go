package app

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/sso-service/controller/grpc/v1"
	httpv1 "car-sell-buy-system/internal/sso-service/controller/http/v1"
	"car-sell-buy-system/internal/sso-service/usecase"
	"car-sell-buy-system/internal/sso-service/usecase/repo"
	"car-sell-buy-system/pkg/httpserver"
	"car-sell-buy-system/pkg/logger"
	"car-sell-buy-system/pkg/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New("info")

	// Postgres
	pg, err := postgres.New(cfg.Pg.URL)
	if err != nil {
		log.Fatal("Failed to connect DB.", err)
	}
	defer pg.Pool.Close()

	// Use cases
	userUseCase := usecase.NewUserUseCase(repo.NewUserRepo(pg, l))

	// grpc
	grpcServer := grpc.NewServer()
	v1.NewRouter(grpcServer, l, userUseCase)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("listening to port %s. press ctrl+c to cancel.", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	// --

	// http
	handler := gin.New()
	httpv1.NewRouter(handler, l, userUseCase)
	httpServ := httpserver.New(handler, httpserver.WithPort(cfg.Http.Port))
	// --

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("ads - Run - signal: " + s.String())
	case err = <-httpServ.Notify():
		l.Error(fmt.Errorf("ads - Run - httpServer.Notify: %w", err))
	}

	err = httpServ.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("ads - Run - httpServer.Shutdown: %w", err))
	}
}
