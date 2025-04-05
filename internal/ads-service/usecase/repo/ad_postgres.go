package repo

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/internal/ads-service/usecase"
	"car-sell-buy-system/internal/ads-service/usecase/repo/filtermodel"
	"car-sell-buy-system/internal/ads-service/usecase/repo/sortermodel"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type AdRepo struct {
	*postgres.Postgres
}

func NewAdRepo(pg *postgres.Postgres) *AdRepo {
	return &AdRepo{pg}
}

func (r *AdRepo) GetById(ctx context.Context, id int) (*entity.Ad, error) {
	sql, args, err := r.Builder.
		Select(
			// ad
			tableColumn(adTableName, "id"),
			tableColumn(adTableName, "car_id"),
			tableColumn(adTableName, "title"),
			tableColumn(adTableName, "description"),
			tableColumn(adTableName, "price"),
			// car
			tableColumn(carTableName, "vin"),
			tableColumn(carTableName, "is_token_minted"),
			tableColumn(carTableName, "brand"),
			tableColumn(carTableName, "model"),
			tableColumn(carTableName, "year_of_release"),
			tableColumn(carTableName, "image_url"),
		).
		From(adTableName).
		InnerJoin(carTableName + " ON " + tableColumn(carTableName, "id") + " = " + tableColumn(adTableName, "car_id")). // util
		Where(squirrel.Eq{tableColumn(adTableName, "id"): id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AdRepo - GetById - r.Builder: %w", err)
	}
	row := r.Pool.QueryRow(ctx, sql, args...)

	ad := &entity.Ad{}
	err = row.Scan(
		&ad.Id,
		&ad.Car.Id,
		&ad.Title,
		&ad.Description,
		&ad.Price,
		&ad.Car.Vin,
		&ad.Car.IsTokenMinted,
		&ad.Car.Brand,
		&ad.Car.Model,
		&ad.Car.YearOfRelease,
		&ad.Car.ImageUrl,
	)
	if err != nil {
		return nil, fmt.Errorf("AdRepo - GetById - row.Scan: %w", err)
	}

	return ad, nil
}

func (r *AdRepo) Store(ctx context.Context, ad entity.Ad) (entity.Ad, error) {
	sql, args, err := r.Builder.
		Insert(adTableName).
		Columns("car_id", "title", "description", "price").
		Values(ad.Car.Id, ad.Title, ad.Description, ad.Price).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return entity.Ad{}, fmt.Errorf("AdRepo - Store - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&ad.Id)
	if err != nil {
		return entity.Ad{}, fmt.Errorf("AdRepo - Store - row.Scan: %w", err)
	}

	return ad, nil
}

func (r *AdRepo) List(ctx context.Context, dto usecase.BasicListRequestDTO) ([]entity.Ad, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	builder := r.Builder.
		Select(
			// ad
			tableColumn(adTableName, "id"),
			tableColumn(adTableName, "car_id"),
			tableColumn(adTableName, "title"),
			tableColumn(adTableName, "description"),
			tableColumn(adTableName, "price"),
			// car
			tableColumn(carTableName, "vin"),
			tableColumn(carTableName, "is_token_minted"),
			tableColumn(carTableName, "brand"),
			tableColumn(carTableName, "model"),
			tableColumn(carTableName, "year_of_release"),
			tableColumn(carTableName, "image_url"),
			// favorite
			//squirrel.Case().When("user_favorites.user_id = 1").to,
			"CASE WHEN user_favorites.user_id = @userId THEN 1 ELSE 0 END AS is_favorite",
		).
		From(adTableName).
		InnerJoin(carTableName + " ON " + tableColumn(carTableName, "id") + " = " + tableColumn(adTableName, "car_id")). // util
		LeftJoin(userFavoritesTableName + " ON " + tableColumn(userFavoritesTableName, "ad_id") + " = " + tableColumn(adTableName, "id")).
		From(adTableName)

	if true {
		//builder = builder.Where(squirrel.Eq{tableColumn(userFavoritesTableName, "user_id"): userId})
	}

	builder, err := sqlutil.ApplyFilters(builder, &filtermodel.AdFilter{}, dto.Filter)
	if err != nil {
		return nil, fmt.Errorf("AdRepo - List - sqlutil.ApplyFilters: %w", err)
	}

	builder, err = sqlutil.ApplySorts(builder, &sortermodel.AdSorter{}, dto.Sort)
	if err != nil {
		return nil, fmt.Errorf("AdRepo - List - sqlutil.ApplySorts: %w", err)
	}

	builder = dto.Pagination.ApplyPagination(builder)

	sql, _, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("AdRepo - List - r.Builder: %w", err)
	}
	fmt.Println(sql)

	args := pgx.NamedArgs{
		"userId": userId,
	}
	rows, err := r.Pool.Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("AdRepo - List - r.Pool.Query: %w", err)
	}

	var ads []entity.Ad
	for rows.Next() {
		var ad entity.Ad
		err = rows.Scan(
			&ad.Id,
			&ad.Car.Id,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Car.Vin,
			&ad.Car.IsTokenMinted,
			&ad.Car.Brand,
			&ad.Car.Model,
			&ad.Car.YearOfRelease,
			&ad.Car.ImageUrl,
			&ad.IsFavorite,
		)
		if err != nil {
			return nil, fmt.Errorf("AdRepo - List - rows.Scan: %w", err)
		}
		ads = append(ads, ad)
	}

	return ads, nil
}

func (r *AdRepo) HandleFavorite(ctx context.Context, adId, userId int) error {
	isFavorite, err := r.isFavorite(ctx, adId, userId)
	if err != nil {
		return err
	}

	fmt.Println("is fav: ", isFavorite)

	if isFavorite {
		sql, args, err := r.Builder.
			Delete("user_favorites").
			Where(squirrel.Eq{"user_id": userId, "ad_id": adId}).
			ToSql()
		if err != nil {
			return fmt.Errorf("AdRepo - HandleFavorite - r.Builder: %w", err)
		}

		_, err = r.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("AdRepo - HandleFavorite - r.Pool.Exec: %w", err)
		}

		return nil
	}

	sql, args, err := r.Builder.
		Insert("user_favorites").
		Columns("user_id", "ad_id").
		Values(userId, adId).
		ToSql()
	if err != nil {
		return fmt.Errorf("AdRepo - HandleFavorite - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AdRepo - HandleFavorite - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *AdRepo) isFavorite(ctx context.Context, adId, userId int) (bool, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("user_favorites").
		Where(squirrel.Eq{"user_id": userId, "ad_id": adId}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("AdRepo - isFavorite - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("adRepo - isFavorite - r.Pool.Query: %v", err)
	}

	rowsProcessed := 0
	for rows.Next() {
		rowsProcessed++
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("adRepo - isFavorite - rows.Err: %v", err)
		return false, err
	}

	if rowsProcessed == 0 {
		fmt.Println("adRepo - isFavorite - rows == 0")
		return false, nil
	}

	fmt.Println("adRepo - isFavorite - rows > 0")
	return true, nil
}
