package payment

import (
	"car-sell-buy-system/internal/ads-service/domain/tariff"
	"context"
)

type Repository interface {
	Store(ctx context.Context, payment Payment) (Payment, error)
}

type TariffRepository interface {
	GetById(ctx context.Context, id int64) (tariff.Tariff, error)
}

type ApiRepository interface {
	CreatePayment(ctx context.Context, payment Payment) (Payment, error)
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
