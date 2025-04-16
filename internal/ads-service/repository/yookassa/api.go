package yookassa

import (
	"bytes"
	"car-sell-buy-system/internal/ads-service/domain/payment"
	"car-sell-buy-system/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strconv"
)

type Repository struct {
	logger logger.Interface
}

func NewRepository(logger logger.Interface) *Repository {
	return &Repository{
		logger: logger,
	}
}

type CreatePaymentRequest struct {
	Amount            Amount            `json:"amount"`
	PaymentMethodData PaymentMethodData `json:"payment_method_data"`
	Confirmation      Confirmation      `json:"confirmation"`
	Description       string            `json:"description"`
}

type PaymentMethodData struct {
	Type string `json:"type"`
}

func (r *Repository) CreatePayment(ctx context.Context, pmnt payment.Payment) (payment.Payment, error) {
	request := CreatePaymentRequest{
		Amount: Amount{
			Value:    strconv.FormatFloat(pmnt.Tariff.Price, 'f', -1, 64),
			Currency: pmnt.Tariff.Currency,
		},
		PaymentMethodData: PaymentMethodData{
			Type: "bank_card",
		},
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnUrl: "https://www.google.com/",
		},
		Description: pmnt.Tariff.Description,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		r.logger.Error("Repository - yookassa - createPayment - json.Marshal: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to make request payment: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.yookassa.ru/v3/payments",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		r.logger.Error("Repository - yookassa - createPayment - http.NewRequest: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to create payment: %w", err)
	}

	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("Idempotence-Key", uuid.New().String())
	req.SetBasicAuth("1070206", "test_FSVL5LucilxLPlZd1siu7o5Bu8flLlJlCiKv8E2wR1A")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		r.logger.Error("Repository - yookassa - createPayment - client.Post: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to create payment: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error("Repository - yookassa - createPayment - io.ReadAll: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to read response create payment: %w", err)
	}
	r.logger.Info(fmt.Sprintf("Получены данные из Yookassa при регистрации платежа: %s", body))

	if err = json.Unmarshal(body, &pmnt); err != nil {
		r.logger.Error("Repository - yookassa - createPayment - json.Unmarshal: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to serialize response payment: %w", err)
	}

	return pmnt, nil
}
