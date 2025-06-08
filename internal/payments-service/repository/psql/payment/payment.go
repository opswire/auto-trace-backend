package payment

import (
	"car-sell-buy-system/internal/payments-service/domain/payment"
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

func (r *Repository) GetByTransactionId(ctx context.Context, id string) (payment.Payment, error) {
	sql, args, err := r.Builder.
		Select(
			// ad
			sqlutil.TableColumn(tableName, "user_id"),
			sqlutil.TableColumn(tableName, "ad_id"),
			sqlutil.TableColumn(tableName, "tariff_id"),
			sqlutil.TableColumn(tableName, "status"),
			sqlutil.TableColumn("ads", "title"),
			sqlutil.TableColumn("tariffs", "price"),
			sqlutil.TableColumn("users", "email"),
		).
		From(tableName).
		InnerJoin("ads on ads.id = payments.ad_id").
		InnerJoin("tariffs on tariffs.tariff_id = payments.tariff_id").
		InnerJoin("users on users.id = payments.user_id").
		Where(squirrel.Eq{sqlutil.TableColumn(tableName, "transaction_id"): id}).
		ToSql()
	if err != nil {
		return payment.Payment{}, fmt.Errorf("paymentRepository - GetByTransactionId - r.Builder: %w", err)
	}
	fmt.Println(sql, args, id)

	var pmnt payment.Payment
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&pmnt.UserId,
			&pmnt.AdId,
			&pmnt.Tariff.Id,
			&pmnt.Status,
			&pmnt.AdTitle,
			&pmnt.Tariff.Price,
			&pmnt.UserEmail,
		)
	if err != nil {
		return payment.Payment{}, fmt.Errorf("paymentRepository - GetByTransactionId - row.Scan: %w", err)
	}

	return pmnt, nil
}
