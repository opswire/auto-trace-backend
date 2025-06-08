package app

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/ads-service/controller/http"
	"car-sell-buy-system/internal/ads-service/domain/ad"
	"car-sell-buy-system/internal/ads-service/repository/psql"
	"car-sell-buy-system/pkg/httpserver"
	"car-sell-buy-system/pkg/logger"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/storage/local"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New("debug")

	// Postgres
	pg, err := postgres.New(cfg.Pg.URL)
	if err != nil {
		log.Fatal("Failed to connect DB.", err)
	}
	defer pg.Pool.Close()

	adService := ad.NewService(
		psql.NewAdRepository(pg, l),
		local.NewFileStorage("./storage"),
	)

	handler := gin.New()
	http.NewRouter(handler, l, cfg, adService)
	httpServ := httpserver.New(handler, httpserver.WithPort(cfg.Http.Port))

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
