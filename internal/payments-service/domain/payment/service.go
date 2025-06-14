package payment

import (
	"car-sell-buy-system/internal/payments-service/domain/tariff"
	"context"
)

type Repository interface {
	Store(ctx context.Context, payment Payment) (Payment, error)
	UpdateStatusByTransactionId(ctx context.Context, transactionId string, status string) error
	GetByTransactionId(ctx context.Context, id string) (Payment, error)
}

type TariffRepository interface {
	GetById(ctx context.Context, id int64) (tariff.Tariff, error)
}

type ApiRepository interface {
	CreatePayment(ctx context.Context, payment Payment) (Payment, error)
	ConfirmPayment(ctx context.Context, transactionId string) (Payment, error)
}

type Service struct {
	repository       Repository
	tariffRepository TariffRepository
	apiRepository    ApiRepository
}

func NewService(repository Repository, tariffRepository TariffRepository, apiRepository ApiRepository) *Service {
	return &Service{
		repository:       repository,
		tariffRepository: tariffRepository,
		apiRepository:    apiRepository,
	}
}

func (s *Service) CreatePayment(ctx context.Context, dto CreatePaymentDto) (Payment, error) {
	trf, err := s.tariffRepository.GetById(ctx, dto.TariffId)
	if err != nil {
		return Payment{}, err
	}

	pmnt := Payment{
		AdId:   dto.AdId,
		UserId: dto.UserId,
		Tariff: trf,
	}

	p, err := s.apiRepository.CreatePayment(ctx, pmnt)
	if err != nil {
		return Payment{}, err
	}

	p, err = s.repository.Store(ctx, p)
	if err != nil {
		return Payment{}, err
	}

	return p, nil
}

func (s *Service) ProcessWebhook(ctx context.Context, dto ProcessWebhookPaymentDto) (string, error) {
	status := dto.Status

	if status == "waiting_for_capture" {
		response, err := s.apiRepository.ConfirmPayment(ctx, dto.TransactionId)
		if err != nil {
			return "", nil
		}

		status = response.Status
	}

	if err := s.repository.UpdateStatusByTransactionId(ctx, dto.TransactionId, status); err != nil {
		return "", err
	}

	return status, nil
}

func (s *Service) CreateEvent(ctx context.Context, id string) (ConfirmedEvent, error) {
	pmnt, err := s.repository.GetByTransactionId(ctx, id)
	if err != nil {
		return ConfirmedEvent{}, err
	}

	event := ConfirmedEvent{
		PaymentID: id,
		UserEmail: pmnt.UserEmail,
		Amount:    int(pmnt.Tariff.Price),
		AdTitle:   pmnt.AdTitle,
	}

	return event, nil
}
