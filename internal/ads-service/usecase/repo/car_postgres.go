package repo

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type CarRepo struct {
	*postgres.Postgres
}

func NewCarRepo(pg *postgres.Postgres) *CarRepo {
	return &CarRepo{pg}
}

func (r *CarRepo) GetById(ctx context.Context, id int) (entity.Car, error) {
	sql, args, err := r.Builder.
		Select("*").
		From(carTableName).
		Where(squirrel.Eq{tableColumn(carTableName, "id"): id}).
		ToSql()
	if err != nil {
		return entity.Car{}, fmt.Errorf("CarRepo - GetById - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Car{}, fmt.Errorf("CarRepo - GetById - r.Pool.Query: %w", err)
	}

	car, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Car])
	if err != nil {
		return entity.Car{}, fmt.Errorf("CarRepo - GetById - pgx.CollectOneRow: %w", err)
	}

	return car, nil
}

func (r *CarRepo) GetByVin(ctx context.Context, vin string) (entity.Car, error) {
	sql, args, err := r.Builder.
		Select("*").
		From(carTableName).
		Where(squirrel.Eq{tableColumn(carTableName, "vin"): vin}).
		ToSql()
	if err != nil {
		return entity.Car{}, fmt.Errorf("CarRepo - GetByVin - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Car{}, fmt.Errorf("CarRepo - GetByVin - r.Pool.Query: %w", err)
	}

	car, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Car])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Car{}, nil
		}
		return entity.Car{}, fmt.Errorf("CarRepo - GetById - pgx.CollectOneRow: %w", err)
	}

	return car, nil
}

func (r *CarRepo) Store(ctx context.Context, car entity.Car) (entity.Car, error) {
	sql, args, err := r.Builder.
		Insert(carTableName).
		Columns("vin", "brand", "model", "year_of_release", "image_url").
		Values(car.Vin, car.Brand, car.Model, car.YearOfRelease, car.ImageUrl).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return entity.Car{}, fmt.Errorf("CarRepo - Store - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&car.Id)
	if err != nil {
		return entity.Car{}, fmt.Errorf("CarRepo - Store - row.Scan: %w", err)
	}

	return car, nil
}
