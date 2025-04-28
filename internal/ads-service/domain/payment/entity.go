package payment

import (
	"car-sell-buy-system/internal/ads-service/domain/tariff"
	"time"
)

type (
	Payment struct {
		Id               int64
		UserId           int64
		AdId             int64
		TransactionId    string
		Tariff           tariff.Tariff
		Status           string
		ConfirmationLink string
		UpdatedAt        time.Time `json:"updated_at"`
		CreatedAt        time.Time `json:"created_at"`
		ExpiresAt        time.Time `json:"expires_at"`
	}
)
