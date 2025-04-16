package payment

import (
	"car-sell-buy-system/internal/ads-service/domain/payment"
	"car-sell-buy-system/internal/ads-service/domain/tariff"
	"time"
)

type Response struct {
	Id               int64         `json:"id"`
	TransactionId    string        `json:"transaction_id"`
	AdId             int64         `json:"ad_id"`
	UserId           int64         `json:"user_id"`
	Tariff           tariff.Tariff `json:"tariff"`
	Status           string        `json:"status"`
	ConfirmationLink string        `json:"confirmation_link"`
	UpdatedAt        time.Time     `json:"updated_at"`
	CreatedAt        time.Time     `json:"created_at"`
	ExpiresAt        time.Time     `json:"expires_at"`
}

func newResponse(p payment.Payment) Response {
	return Response{
		Id:               p.Id,
		TransactionId:    p.TransactionId,
		AdId:             p.AdId,
		UserId:           p.UserId,
		Tariff:           p.Tariff,
		Status:           p.Status,
		ConfirmationLink: p.ConfirmationLink,
		UpdatedAt:        p.UpdatedAt,
		CreatedAt:        p.CreatedAt,
		ExpiresAt:        p.ExpiresAt,
	}
}
