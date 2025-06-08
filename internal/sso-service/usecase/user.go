package usecase

import (
	"car-sell-buy-system/internal/sso-service/entity"
	"context"
	"errors"
	"fmt"
	"strconv"
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

func (uc *UserUseCase) List(ctx context.Context) ([]*entity.User, error) {
	id, err := strconv.Atoi(ctx.Value("userId").(string))
	if err != nil {
		return nil, err
	}

	currentUser, err := uc.repo.GetById(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	if currentUser.Role != "admin" {
		return nil, fmt.Errorf("Пользователь не имеет доступа!")
	}

	users, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *UserUseCase) HandleActive(ctx context.Context, userId int64) error {
	id, err := strconv.Atoi(ctx.Value("userId").(string))
	if err != nil {
		return err
	}

	currentUser, err := uc.repo.GetById(ctx, int64(id))
	if err != nil {
		return err
	}

	if currentUser.Role != "admin" {
		return fmt.Errorf("Пользователь не имеет доступа!")
	}

	err = uc.repo.HandleActive(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}
