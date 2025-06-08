package appointment

import (
	"car-sell-buy-system/internal/ads-service/domain/ad"
	"context"
	"fmt"
)

type Repository interface {
	StoreAppointment(ctx context.Context, app *Appointment) error
	CheckTimeConflict(ctx context.Context, dto CheckTimeConflictDTO) (bool, error)
	GetAllAppointmentsByUserId(ctx context.Context) ([]*Appointment, error)
	GetAppointmentsByDateRange(ctx context.Context, dto GetAppointmentsByDateRangeDTO) ([]*Appointment, error)
	ConfirmAppointment(ctx context.Context, id int64) error
	MarkAppointmentAsCanceled(ctx context.Context, id int64) error
	FindById(ctx context.Context, id int64) (Appointment, error)
}

type AdRepository interface {
	GetById(ctx context.Context, id int64) (ad.Ad, error)
}

type Service struct {
	repository   Repository
	adRepository AdRepository
}

func NewService(
	repository Repository,
	adRepository AdRepository,
) *Service {
	return &Service{
		repository:   repository,
		adRepository: adRepository,
	}
}

func (s *Service) StoreAppointment(ctx context.Context, dto StoreDTO) (*Appointment, error) {
	appointment := &Appointment{
		Start:    dto.Start,
		Duration: dto.Duration,
		Location: dto.Location,
		AdId:     dto.AdId,
		SellerId: dto.SellerId,
		BuyerId:  dto.BuyerId,
	}

	err := s.repository.StoreAppointment(ctx, appointment)
	if err != nil {
		return nil, err
	}

	return appointment, nil
}

func (s *Service) CheckTimeConflict(ctx context.Context, dto CheckTimeConflictDTO) (bool, error) {
	flag, err := s.repository.CheckTimeConflict(ctx, dto)
	if err != nil {
		return true, err
	}

	return flag, nil
}

func (s *Service) GetAllAppointmentsByUserId(ctx context.Context) ([]*Appointment, error) {
	apps, err := s.repository.GetAllAppointmentsByUserId(ctx)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (s *Service) GetAppointmentsByDateRange(ctx context.Context, dto GetAppointmentsByDateRangeDTO) ([]*Appointment, error) {
	apps, err := s.repository.GetAppointmentsByDateRange(ctx, dto)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (s *Service) ConfirmAppointment(ctx context.Context, id int64) error {
	app, err := s.repository.FindById(ctx, id)
	if err != nil {
		return err
	}

	if app.SellerId != ctx.Value("userId").(int64) {
		return fmt.Errorf("Только продавец может подтвердить встречу!")
	}

	err = s.repository.ConfirmAppointment(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) MarkAppointmentAsCanceled(ctx context.Context, id int64) error {
	err := s.repository.MarkAppointmentAsCanceled(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
