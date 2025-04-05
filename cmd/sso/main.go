package main

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/sso-service/app"
)

func main() {
	// Configuration
	cfg := config.NewConfig()

	// Run
	app.Run(cfg)
}
