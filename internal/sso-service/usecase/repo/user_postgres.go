package repo

import (
	"car-sell-buy-system/internal/sso-service/entity"
	"car-sell-buy-system/pkg/logger"
	"car-sell-buy-system/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	*postgres.Postgres
	logger *logger.Logger
}

func NewUserRepo(pg *postgres.Postgres, l *logger.Logger) *UserRepo {
	return &UserRepo{pg, l}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"users.email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByEmail - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByEmail - r.Pool.Query: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[entity.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByEmail - pgx.CollectOneRow: %w", err)
	}

	return user, nil
}

func (r *UserRepo) Store(ctx context.Context, user entity.User) (entity.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - Store - bcrypt.GenerateFromPassword: %w", err)
	}

	sql, args, err := r.Builder.
		Insert("users").
		Columns("email", "password", "name").
		Values(user.Email, password, user.Name).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - Store - r.Builder: %w", err)
	}

	r.logger.Info(fmt.Sprintf("sql: %s", sql))
	r.logger.Info(fmt.Sprintf("args: %s", args))

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&user.Id)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - Store - row.Scan: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetById(ctx context.Context, id int64) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"users.id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetById - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetById - r.Pool.Query: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[entity.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetById - pgx.CollectOneRow: %w", err)
	}

	return user, nil
}
