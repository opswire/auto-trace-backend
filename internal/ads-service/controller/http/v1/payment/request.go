package payment

import "car-sell-buy-system/internal/ads-service/domain/payment"

type CreatePaymentRequest struct {
	AdId     int64 `json:"ad_id"`
	TariffId int64 `json:"tariff_id"`
}

func (r CreatePaymentRequest) toDTO(userId int64) payment.CreatePaymentDto {
	return payment.CreatePaymentDto{
		AdId:     r.AdId,
		UserId:   userId,
		TariffId: r.TariffId,
	}
}

type WebhookPaymentRequest struct {
	Type   string `json:"type"`
	Event  string `json:"event"`
	Object struct {
		Id     string `json:"id"`
		Status string `json:"status"`
		Paid   bool   `json:"paid"`
	} `json:"object"`
}

func (r WebhookPaymentRequest) toDTO() payment.ProcessWebhookPaymentDto {
	return payment.ProcessWebhookPaymentDto{
		TransactionId: r.Object.Id,
		Status:        r.Object.Status,
		Paid:          r.Object.Paid,
	}
}
