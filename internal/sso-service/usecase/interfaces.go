package usecase

import (
	"car-sell-buy-system/internal/sso-service/entity"
	"context"
)

type (
	User interface {
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
		Register(ctx context.Context, user entity.User) (entity.User, error)
		List(ctx context.Context) ([]*entity.User, error)
		HandleActive(ctx context.Context, userId int64) error
	}

	UserRepo interface {
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
		GetById(ctx context.Context, id int64) (*entity.User, error)
		Store(ctx context.Context, user entity.User) (entity.User, error)
		List(ctx context.Context) ([]*entity.User, error)
		HandleActive(ctx context.Context, userId int64) error
	}
)
