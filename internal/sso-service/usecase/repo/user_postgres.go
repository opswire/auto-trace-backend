package repo

import (
	"car-sell-buy-system/internal/sso-service/entity"
	"car-sell-buy-system/pkg/logger"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
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
		return nil, fmt.Errorf("UserRepo - GetByTransactionId - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByTransactionId - r.Pool.Query: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[entity.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByTransactionId - pgx.CollectOneRow: %w", err)
	}

	return user, nil
}

func (r *UserRepo) List(ctx context.Context) ([]*entity.User, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - List - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - List - r.Pool.Query: %w", err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[entity.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepo - List - pgx.CollectOneRow: %w", err)
	}

	return users, nil
}

func (r *UserRepo) HandleActive(ctx context.Context, userId int64) error {
	user, err := r.GetById(ctx, userId)
	if err != nil {
		return err
	}

	fmt.Println("is active: ", user.IsActive)

	if user.Role == "admin" {
		return fmt.Errorf("Нельзя заблокировать админа")
	}

	sql, args, err := r.Builder.
		Update("users").
		Set("is_active", user.IsActive == false).
		Where(squirrel.Eq{sqlutil.TableColumn("users", "id"): userId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Store - row.Scan: %w", err)
	}

	return nil
}
