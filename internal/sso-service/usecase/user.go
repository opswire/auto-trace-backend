package usecase

import (
	"car-sell-buy-system/internal/sso-service/entity"
	"context"
	"errors"
)

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(r UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (uc *UserUseCase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	foundUser, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (uc *UserUseCase) Register(ctx context.Context, user entity.User) (entity.User, error) {
	foundUser, err := uc.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return entity.User{}, err
	}

	if foundUser != nil {
		return entity.User{}, errors.New("пользователь с такой почтой уже существует")
	}

	storedUser, err := uc.repo.Store(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	return storedUser, nil
}
