package payment

import (
	"car-sell-buy-system/internal/ads-service/domain/payment"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
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
			"user_id",
			"ad_id",
			"tariff_id",
			"status",
			"transaction_id",
			"confirmation_link",
			"expires_at",
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
		//Suffix("RETURNING payment_id").
		ToSql()
	if err != nil {
		return payment.Payment{}, fmt.Errorf("paymentRepository - Store - r.Builder: %w", err)
	}
	fmt.Println("sql: ", sql)
	fmt.Println("args: ", args)

	_, err = r.Pool.Exec(ctx, sql, args...)
	//err = row.Scan(&pmnt.Id)
	if err != nil {
		return payment.Payment{}, fmt.Errorf("paymentRepository - Store - row.Scan: %w", err)
	}

	return pmnt, nil
}

func (r *Repository) UpdateStatusByTransactionId(ctx context.Context, transactionId string, status string) error {
	sql, args, err := r.Builder.
		Update(tableName).
		Set("status", status).
		Where(squirrel.Eq{
			sqlutil.TableColumn(tableName, "transaction_id"): transactionId,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("paymentRepository - UpdateStatusByTransactionId - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("paymentRepository - UpdateStatusByTransactionId - row.Exec: %w", err)
	}

	return nil
}
