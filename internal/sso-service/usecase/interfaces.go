package usecase

import (
	"car-sell-buy-system/internal/sso-service/entity"
	"context"
)

type (
	User interface {
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
		Register(ctx context.Context, user entity.User) (entity.User, error)
	}

	UserRepo interface {
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
		Store(ctx context.Context, user entity.User) (entity.User, error)
	}
)
