package payment

import (
	"car-sell-buy-system/internal/ads-service/domain/payment"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"fmt"
)

const (
	tableName = "payments"
)

type Repository struct {
	*postgres.Postgres
}

func NewRepository(pg *postgres.Postgres) *Repository {
	return &Repository{
		pg,
	}
}

func (r *Repository) Store(ctx context.Context, pmnt payment.Payment) (payment.Payment, error) {
	sql, args, err := r.Builder.
		Insert(tableName).
		Columns(
			sqlutil.TableColumn(tableName, "user_id"),
			sqlutil.TableColumn(tableName, "ad_id"),
			sqlutil.TableColumn(tableName, "tariff_id"),
			sqlutil.TableColumn(tableName, "status"),
			sqlutil.TableColumn(tableName, "transaction_id"),
			sqlutil.TableColumn(tableName, "confirmation_link"),
			sqlutil.TableColumn(tableName, "expires_at"),
		).
		Values(
			pmnt.UserId,
			pmnt.AdId,
			pmnt.Tariff.Id,
			pmnt.Status,
			pmnt.TransactionId,
			pmnt.ConfirmationLink,
			pmnt.ExpiresAt,
		).
		Suffix("RETURNING payment_id").
		ToSql()
	if err != nil {
		return payment.Payment{}, fmt.Errorf("paymentRepository - Store - r.Builder: %w", err)
	}

	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&pmnt.Id)
	if err != nil {
		return payment.Payment{}, fmt.Errorf("paymentRepository - Store - row.Scan: %w", err)
	}

	return pmnt, nil
}
