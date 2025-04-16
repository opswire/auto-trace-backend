package payment

import "car-sell-buy-system/internal/ads-service/domain/payment"

type CreatePaymentRequest struct {
	AdId     int64 `json:"ad_id"`
	UserId   int64 `json:"user_id"`
	TariffId int64 `json:"tariff_id"`
}

func (r CreatePaymentRequest) toDTO() payment.CreatePaymentDto {
	return payment.CreatePaymentDto{
		AdId:     r.AdId,
		UserId:   r.UserId,
		TariffId: r.TariffId,
	}
}
