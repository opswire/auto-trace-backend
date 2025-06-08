package payment_test

import (
	"car-sell-buy-system/internal/payments-service/domain/payment"
	"context"
	"errors"
	"go.uber.org/mock/gomock"
	"testing"

	"car-sell-buy-system/internal/payments-service/domain/tariff"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*gomock.Controller, *payment.Service, *MockRepository, *MockTariffRepository, *MockApiRepository) {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	tariffRepo := NewMockTariffRepository(ctrl)
	apiRepo := NewMockApiRepository(ctrl)
	svc := payment.NewService(repo, tariffRepo, apiRepo)
	return ctrl, svc, repo, tariffRepo, apiRepo
}

func TestService_CreatePayment_Success(t *testing.T) {
	ctrl, svc, repo, tariffRepo, apiRepo := setup(t)
	defer ctrl.Finish()

	dto := payment.CreatePaymentDto{
		AdId:     1,
		UserId:   2,
		TariffId: 3,
	}

	tariffObj := tariff.Tariff{Id: 3, Name: "Premium"}
	paymentInput := payment.Payment{
		AdId:   1,
		UserId: 2,
		Tariff: tariffObj,
	}

	apiPayment := paymentInput
	apiPayment.Id = 100

	finalPayment := apiPayment

	tariffRepo.EXPECT().GetById(gomock.Any(), dto.TariffId).Return(tariffObj, nil)
	apiRepo.EXPECT().CreatePayment(gomock.Any(), paymentInput).Return(apiPayment, nil)
	repo.EXPECT().Store(gomock.Any(), apiPayment).Return(finalPayment, nil)

	result, err := svc.CreatePayment(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, finalPayment, result)
}

func TestService_CreatePayment_TariffError(t *testing.T) {
	ctrl, svc, _, tariffRepo, _ := setup(t)
	defer ctrl.Finish()

	dto := payment.CreatePaymentDto{TariffId: 42}
	tariffRepo.EXPECT().GetById(gomock.Any(), dto.TariffId).Return(tariff.Tariff{}, errors.New("not found"))

	_, err := svc.CreatePayment(context.Background(), dto)
	assert.Error(t, err)
}

func TestService_CreatePayment_CreateApiError(t *testing.T) {
	ctrl, svc, _, tariffRepo, apiRepo := setup(t)
	defer ctrl.Finish()

	dto := payment.CreatePaymentDto{AdId: 1, UserId: 2, TariffId: 3}
	tariffObj := tariff.Tariff{Id: 3}

	tariffRepo.EXPECT().GetById(gomock.Any(), dto.TariffId).Return(tariffObj, nil)
	apiRepo.EXPECT().
		CreatePayment(gomock.Any(), gomock.Any()).
		Return(payment.Payment{}, errors.New("api error"))

	_, err := svc.CreatePayment(context.Background(), dto)
	assert.Error(t, err)
}

func TestService_CreatePayment_StoreError(t *testing.T) {
	ctrl, svc, repo, tariffRepo, apiRepo := setup(t)
	defer ctrl.Finish()

	dto := payment.CreatePaymentDto{AdId: 1, UserId: 2, TariffId: 3}
	tariffObj := tariff.Tariff{Id: 3}
	pmnt := payment.Payment{AdId: 1, UserId: 2, Tariff: tariffObj}
	apiResp := pmnt

	tariffRepo.EXPECT().GetById(gomock.Any(), dto.TariffId).Return(tariffObj, nil)
	apiRepo.EXPECT().CreatePayment(gomock.Any(), pmnt).Return(apiResp, nil)
	repo.EXPECT().Store(gomock.Any(), apiResp).Return(payment.Payment{}, errors.New("db error"))

	_, err := svc.CreatePayment(context.Background(), dto)
	assert.Error(t, err)
}

func TestService_ProcessWebhook_StatusConfirmed(t *testing.T) {
	ctrl, svc, repo, _, apiRepo := setup(t)
	defer ctrl.Finish()

	dto := payment.ProcessWebhookPaymentDto{
		TransactionId: "tx123",
		Status:        "waiting_for_capture",
	}

	apiResp := payment.Payment{Status: "succeeded"}

	apiRepo.EXPECT().ConfirmPayment(gomock.Any(), dto.TransactionId).Return(apiResp, nil)
	repo.EXPECT().UpdateStatusByTransactionId(gomock.Any(), dto.TransactionId, "succeeded").Return(nil)

	status, err := svc.ProcessWebhook(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, "succeeded", status)
}

func TestService_ProcessWebhook_UpdateFails(t *testing.T) {
	ctrl, svc, repo, _, _ := setup(t)
	defer ctrl.Finish()

	dto := payment.ProcessWebhookPaymentDto{
		TransactionId: "tx999",
		Status:        "succeeded",
	}

	repo.EXPECT().UpdateStatusByTransactionId(gomock.Any(), dto.TransactionId, "succeeded").Return(errors.New("update error"))

	_, err := svc.ProcessWebhook(context.Background(), dto)
	assert.Error(t, err)
}
