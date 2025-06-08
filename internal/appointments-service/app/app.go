package app

import (
	"car-sell-buy-system/config"
	adPsql "car-sell-buy-system/internal/ads-service/repository/psql"
	"car-sell-buy-system/internal/appointments-service/controller/http"
	"car-sell-buy-system/internal/appointments-service/domain/appointment"
	appointmentPsql "car-sell-buy-system/internal/appointments-service/repository/psql/appointment"
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
	appointmentService := appointment.NewService(
		appointmentPsql.NewRepository(pg),
		adPsql.NewAdRepository(pg, l),
	)

	handler := gin.New()
	http.NewRouter(handler, l, cfg, appointmentService)
	httpServ := httpserver.New(handler, httpserver.WithPort(cfg.Http.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("appointments - Run - signal: " + s.String())
	case err = <-httpServ.Notify():
		l.Error(fmt.Errorf("appointments - Run - httpServer.Notify: %w", err))
	}

	err = httpServ.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("appointments - Run - httpServer.Shutdown: %w", err))
	}
}
