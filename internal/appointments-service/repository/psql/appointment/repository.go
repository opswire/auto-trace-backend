package appointment

import (
	"car-sell-buy-system/internal/appointments-service/domain/appointment"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"time"
)

const (
	appointmentsTableName = "appointments"
)

type Repository struct {
	*postgres.Postgres
}

func NewRepository(pg *postgres.Postgres) *Repository {
	return &Repository{
		pg,
	}
}

func (r *Repository) StoreAppointment(ctx context.Context, app *appointment.Appointment) error {
	buyerUserId := ctx.Value("userId")
	fmt.Println("userId: ", buyerUserId)

	query, args, err := r.Builder.
		Insert(appointmentsTableName).
		Columns(
			"start_time",
			"duration",
			"location",
			"ad_id",
			"buyer_id",
			"is_confirmed",
			"is_canceled",
		).
		Values(
			app.Start,
			app.Duration,
			app.Location,
			app.AdId,
			buyerUserId,
			app.IsConfirmed,
			app.IsCanceled,
		).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	return r.Pool.QueryRow(ctx, query, args...).Scan(&app.ID)
}

func (r *Repository) CheckTimeConflict(
	ctx context.Context,
	dto appointment.CheckTimeConflictDTO,
) (bool, error) {
	end := dto.StartTime.Add(time.Duration(dto.Duration) * time.Minute)

	subQuery := r.Builder.
		Select("1").
		From(appointmentsTableName).
		Where(squirrel.And{
			squirrel.Eq{"ad_id": dto.AdId},
			squirrel.Expr("start_time < ?", end),
			squirrel.Expr("(start_time + (duration * interval '1 minute')) > ?", dto.StartTime),
		})

	query, args, err := r.Builder.
		Select().
		Column(squirrel.Expr("EXISTS (?)", subQuery)).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool
	err = r.Pool.QueryRow(ctx, query, args...).Scan(&exists)
	return exists, err
}

func (r *Repository) GetAllAppointmentsByUserId(ctx context.Context) ([]*appointment.Appointment, error) {
	buyerUserId := ctx.Value("userId")
	fmt.Println("userId: ", buyerUserId)

	query, args, err := r.Builder.
		Select(
			"app.id",
			"app.start_time",
			"app.duration",
			"app.location",
			"app.ad_id",
			"ads.user_id",
			"app.buyer_id",
			"app.is_confirmed",
			"app.is_canceled",
			"ads.title",
			"buyers.name",
			"sellers.name",
		).
		From("appointments app").
		InnerJoin("ads ON app.ad_id = ads.id").
		InnerJoin("users as buyers on buyers.id = app.buyer_id").
		InnerJoin("users as sellers on sellers.id = ads.user_id").
		Where(
			squirrel.Or{
				squirrel.Eq{"app.buyer_id": buyerUserId},
				squirrel.Eq{"ads.user_id": buyerUserId},
			},
		).
		OrderBy("app.start_time").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	fmt.Println("sql: ", query)
	fmt.Println("args: ", args)

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*appointment.Appointment
	for rows.Next() {
		var app appointment.Appointment
		err = rows.Scan(
			&app.ID,
			&app.Start,
			&app.Duration,
			&app.Location,
			&app.AdId,
			&app.SellerId,
			&app.BuyerId,
			&app.IsConfirmed,
			&app.IsCanceled,
			&app.AdTitle,
			&app.BuyerName,
			&app.SellerName,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &app)
	}

	return appointments, rows.Err()
}

func (r *Repository) GetAppointmentsByDateRange(ctx context.Context, dto appointment.GetAppointmentsByDateRangeDTO) ([]*appointment.Appointment, error) {
	buyerUserId := ctx.Value("userId")
	fmt.Println("userId: ", buyerUserId)

	query, args, err := r.Builder.
		Select(
			"id", "start_time", "duration", "location", "ad_id",
			"seller_id", "buyer_id", "is_confirmed", "is_canceled",
		).
		From(appointmentsTableName).
		Where(squirrel.And{
			squirrel.Or{
				squirrel.Eq{"app.buyer_id": buyerUserId},
				squirrel.Eq{"ads.user_id": buyerUserId},
			},
			squirrel.Eq{"ad_id": dto.AdId},
			squirrel.GtOrEq{"start_time": dto.StartDate},
			squirrel.LtOrEq{"start_time": dto.EndDate},
		}).
		OrderBy("start_time").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*appointment.Appointment
	for rows.Next() {
		var app appointment.Appointment
		err = rows.Scan(
			&app.ID,
			&app.Start,
			&app.Duration,
			&app.Location,
			&app.AdId,
			&app.SellerId,
			&app.BuyerId,
			&app.IsConfirmed,
			&app.IsCanceled,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &app)
	}

	return appointments, rows.Err()
}

func (r *Repository) ConfirmAppointment(ctx context.Context, id int64) error {
	query, args, err := r.Builder.
		Update(appointmentsTableName).
		Set("is_confirmed", true).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build ConfirmAppointment query: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	return err
}

func (r *Repository) MarkAppointmentAsCanceled(ctx context.Context, id int64) error {
	query, args, err := r.Builder.
		Update(appointmentsTableName).
		Set("is_canceled", true).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build MarkAppointmentAsCanceled query: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	return err
}

func (r *Repository) FindById(ctx context.Context, id int64) (appointment.Appointment, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	sql, args, err := r.Builder.
		Select(
			sqlutil.TableColumn(appointmentsTableName, "id"),
			sqlutil.TableColumn(appointmentsTableName, "start_time"),
			sqlutil.TableColumn(appointmentsTableName, "duration"),
			sqlutil.TableColumn(appointmentsTableName, "location"),
			sqlutil.TableColumn(appointmentsTableName, "ad_id"),
			sqlutil.TableColumn("ads", "user_id"),
			sqlutil.TableColumn(appointmentsTableName, "buyer_id"),
			sqlutil.TableColumn(appointmentsTableName, "is_confirmed"),
			sqlutil.TableColumn(appointmentsTableName, "is_canceled"),
		).
		From(appointmentsTableName).
		InnerJoin("ads on ads.id = appointments.ad_id").
		Where(squirrel.Eq{sqlutil.TableColumn(appointmentsTableName, "id"): id}).
		ToSql()
	if err != nil {
		return appointment.Appointment{}, fmt.Errorf("AppointmentRepository - GetByTransactionId - r.Builder: %w", err)
	}
	fmt.Println("sql: ", sql)
	fmt.Println("args: ", args)

	var app appointment.Appointment
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&app.ID,
			&app.Start,
			&app.Duration,
			&app.Location,
			&app.AdId,
			&app.SellerId,
			&app.BuyerId,
			&app.IsConfirmed,
			&app.IsCanceled,
		)
	if err != nil {
		return appointment.Appointment{}, fmt.Errorf("AppointmentRepository - FindById - row.Scan: %w", err)
	}

	return app, nil
}
