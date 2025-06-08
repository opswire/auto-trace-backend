package tariff

import (
	"car-sell-buy-system/internal/payments-service/domain/tariff"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

const (
	tableName = "tariffs"
)

type Repository struct {
	*postgres.Postgres
}

func NewRepository(pg *postgres.Postgres) *Repository {
	return &Repository{
		pg,
	}
}

func (r *Repository) GetById(ctx context.Context, id int64) (tariff.Tariff, error) {
	sql, args, err := r.Builder.
		Select(
			// ad
			sqlutil.TableColumn(tableName, "tariff_id"),
			sqlutil.TableColumn(tableName, "name"),
			sqlutil.TableColumn(tableName, "description"),
			sqlutil.TableColumn(tableName, "price"),
			sqlutil.TableColumn(tableName, "currency"),
			sqlutil.TableColumn(tableName, "duration_min"),
			sqlutil.TableColumn(tableName, "is_active"),
			sqlutil.TableColumn(tableName, "created_at"),
			sqlutil.TableColumn(tableName, "updated_at"),
		).
		From(tableName).
		Where(squirrel.Eq{sqlutil.TableColumn(tableName, "tariff_id"): id}).
		ToSql()
	if err != nil {
		return tariff.Tariff{}, fmt.Errorf("TariffRepository - GetByTransactionId - r.Builder: %w", err)
	}
	fmt.Println(sql, args, id)

	var t tariff.Tariff
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&t.Id,
			&t.Name,
			&t.Description,
			&t.Price,
			&t.Currency,
			&t.DurationMin,
			&t.IsActive,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
	if err != nil {
		return tariff.Tariff{}, fmt.Errorf("TariffRepository - GetByTransactionId - row.Scan: %w", err)
	}

	return t, nil
}
