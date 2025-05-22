package app

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/nft-service/controller/http"
	"car-sell-buy-system/internal/nft-service/domain/nft"
	nftpsql "car-sell-buy-system/internal/nft-service/repository/psql/nft"
	"car-sell-buy-system/internal/nft-service/repository/webapi"
	"car-sell-buy-system/pkg/httpserver"
	"car-sell-buy-system/pkg/logger"
	"car-sell-buy-system/pkg/postgres"
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

	// Services
	nftService := nft.NewService(
		nftpsql.NewRepository(pg),
		webapi.NewNftEthereumWebAPI(),
	)

	handler := gin.New()
	http.NewRouter(handler, l, cfg, nftService)
	httpServ := httpserver.New(handler, httpserver.WithPort(cfg.Http.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("nft Service - Run - signal: " + s.String())
	case err = <-httpServ.Notify():
		l.Error(fmt.Errorf("nft Service - Run - httpServer.Notify: %w", err))
	}

	err = httpServ.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("nft Service - Run - httpServer.Shutdown: %w", err))
	}
}
