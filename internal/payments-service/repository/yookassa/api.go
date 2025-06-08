package yookassa

import (
	"bytes"
	"car-sell-buy-system/internal/payments-service/domain/payment"
	"car-sell-buy-system/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Repository struct {
	logger logger.Interface
	//config config.Config
}

func NewRepository(logger logger.Interface) *Repository {
	return &Repository{
		logger: logger,
	}
}

func (r *Repository) CreatePayment(ctx context.Context, pmnt payment.Payment) (payment.Payment, error) {
	request := CreatePaymentRequest{
		Amount: AmountResponse{
			Value:    strconv.FormatFloat(pmnt.Tariff.Price, 'f', -1, 64),
			Currency: pmnt.Tariff.Currency,
		},
		PaymentMethodData: PaymentMethodDataRequest{
			Type: "bank_card",
		},
		Confirmation: ConfirmationResponse{
			Type: "redirect",
			ReturnUrl: fmt.Sprintf(
				"http://localhost:3000/ads/%d?success=%s",
				pmnt.AdId,
				url.QueryEscape(fmt.Sprintf("Объявление c ID %d было успешно продвинуто", pmnt.AdId)),
			),
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

	p := PaymentResponse{}
	if err = json.Unmarshal(body, &p); err != nil {
		r.logger.Error("Repository - yookassa - createPayment - json.Unmarshal: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to serialize response payment: %w", err)
	}

	pmnt.TransactionId = p.Id
	pmnt.ExpiresAt = p.CreatedAt.Add(time.Duration(pmnt.Tariff.DurationMin) * time.Minute)
	pmnt.Status = p.Status
	pmnt.ConfirmationLink = p.Confirmation.ConfirmationUrl

	return pmnt, nil
}

func (r *Repository) ConfirmPayment(ctx context.Context, transactionId string) (payment.Payment, error) {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s/capture", transactionId),
		nil,
	)
	if err != nil {
		r.logger.Error("Repository - yookassa - ConfirmPayment - http.NewRequest: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to create payment: %w", err)
	}

	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("Idempotence-Key", uuid.New().String())
	req.SetBasicAuth("1070206", "test_FSVL5LucilxLPlZd1siu7o5Bu8flLlJlCiKv8E2wR1A")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		r.logger.Error("Repository - yookassa - ConfirmPayment - client.Post: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to create payment: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error("Repository - yookassa - ConfirmPayment - io.ReadAll: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to read response create payment: %w", err)
	}
	r.logger.Info(fmt.Sprintf("Получены данные из Yookassa при подтверждении платежа: %s", body))

	p := PaymentResponse{}
	if err = json.Unmarshal(body, &p); err != nil {
		r.logger.Error("Repository - yookassa - ConfirmPayment - json.Unmarshal: %s", err)
		return payment.Payment{}, fmt.Errorf("failed to serialize response payment: %w", err)
	}

	pmnt := payment.Payment{
		TransactionId:    p.Id,
		Status:           p.Status,
		ConfirmationLink: p.Confirmation.ConfirmationUrl,
		ExpiresAt:        p.ExpiresAt,
	}

	return pmnt, nil
}
