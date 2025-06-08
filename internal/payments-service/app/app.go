package app

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/payments-service/controller/http"
	"car-sell-buy-system/internal/payments-service/domain/payment"
	paymentpsql "car-sell-buy-system/internal/payments-service/repository/psql/payment"
	"car-sell-buy-system/internal/payments-service/repository/psql/tariff"
	"car-sell-buy-system/internal/payments-service/repository/yookassa"
	"car-sell-buy-system/pkg/httpserver"
	"car-sell-buy-system/pkg/logger"
	"car-sell-buy-system/pkg/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
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

	paymentService := payment.NewService(paymentpsql.NewRepository(pg), tariff.NewRepository(pg), yookassa.NewRepository(l))

	publisher := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "payments",
		Balancer: &kafka.LeastBytes{},
	}

	handler := gin.New()
	http.NewRouter(handler, l, cfg, paymentService, publisher)
	httpServ := httpserver.New(handler, httpserver.WithPort(cfg.Http.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("Payments Service - Run - signal: " + s.String())
	case err = <-httpServ.Notify():
		l.Error(fmt.Errorf("payments Service - Run - httpServer.Notify: %w", err))
	}

	err = httpServ.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("payments Service - Run - httpServer.Shutdown: %w", err))
	}
}
