package persistent

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/internal/ads-service/repository/persistent/filter"
	"car-sell-buy-system/internal/ads-service/repository/persistent/sort"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	adTableName            = "ads"
	userFavoritesTableName = "user_favorites"
)

type AdRepository struct {
	*postgres.Postgres
}

func NewAdRepository(pg *postgres.Postgres) *AdRepository {
	return &AdRepository{
		pg,
	}
}

func (r *AdRepository) GetById(ctx context.Context, id int64) (entity.Ad, error) {
	sql, args, err := r.Builder.
		Select(
			// ad
			sqlutil.TableColumn(adTableName, "id"),
			sqlutil.TableColumn(adTableName, "title"),
			sqlutil.TableColumn(adTableName, "description"),
			sqlutil.TableColumn(adTableName, "price"),
			sqlutil.TableColumn(adTableName, "vin"),
			sqlutil.TableColumn(adTableName, "is_token_minted"),
			sqlutil.TableColumn(adTableName, "brand"),
			sqlutil.TableColumn(adTableName, "model"),
			sqlutil.TableColumn(adTableName, "year_of_release"),
		).
		From(adTableName).
		Where(squirrel.Eq{sqlutil.TableColumn(adTableName, "id"): id}).
		ToSql()
	if err != nil {
		return entity.Ad{}, fmt.Errorf("AdRepository - GetById - r.Builder: %w", err)
	}

	var ad entity.Ad
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&ad.Id,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Vin,
			&ad.IsTokenMinted,
			&ad.Brand,
			&ad.Model,
			&ad.YearOfRelease,
		)
	if err != nil {
		return entity.Ad{}, fmt.Errorf("AdRepository - GetById - row.Scan: %w", err)
	}

	return ad, nil
}

func (r *AdRepository) Store(ctx context.Context, dto entity.AdStoreDTO) (entity.Ad, error) {
	sql, args, err := r.Builder.
		Insert(adTableName).
		Columns(
			"title",
			"description",
			"price",
			"vin",
			"brand",
			"model",
			"year_of_release",
			"is_token_minted",
		).
		Values(
			dto.Title,
			dto.Description,
			dto.Price,
			dto.Vin,
			dto.Brand,
			dto.Model,
			dto.YearOfRelease,
			false,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return entity.Ad{}, fmt.Errorf("AdRepository - Store - r.Builder: %w", err)
	}

	ad := entity.Ad{
		Title:         dto.Title,
		Description:   dto.Description,
		Price:         dto.Price,
		Vin:           dto.Vin,
		Brand:         dto.Brand,
		Model:         dto.Model,
		YearOfRelease: dto.YearOfRelease,
		IsTokenMinted: false,
		IsFavorite:    false,
	}
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&ad.Id)
	if err != nil {
		return entity.Ad{}, fmt.Errorf("AdRepository - Store - row.Scan: %w", err)
	}

	return ad, nil
}

func (r *AdRepository) List(ctx context.Context, dto entity.AdListDTO) ([]entity.Ad, uint64, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	builder := r.Builder.
		Select(
			sqlutil.TableColumn(adTableName, "id"),
			sqlutil.TableColumn(adTableName, "title"),
			sqlutil.TableColumn(adTableName, "description"),
			sqlutil.TableColumn(adTableName, "price"),
			sqlutil.TableColumn(adTableName, "vin"),
			sqlutil.TableColumn(adTableName, "is_token_minted"),
			sqlutil.TableColumn(adTableName, "brand"),
			sqlutil.TableColumn(adTableName, "model"),
			sqlutil.TableColumn(adTableName, "year_of_release"),
			"CASE WHEN user_favorites.user_id = @userId THEN true ELSE false END AS is_favorite",
			"COUNT(*) OVER() AS total",
		).
		From(adTableName).
		LeftJoin(userFavoritesTableName + " ON " + sqlutil.TableColumn(userFavoritesTableName, "ad_id") + " = " + sqlutil.TableColumn(adTableName, "id")).
		From(adTableName)

	builder, err := sqlutil.ApplyFilters(builder, &filter.AdFilter{}, dto.Filter)
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - sqlutil.ApplyFilters: %w", err)
	}

	builder, err = sqlutil.ApplySorts(builder, &sort.AdSorter{}, dto.Sort)
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - sqlutil.ApplySorts: %w", err)
	}

	builder = dto.Pagination.ApplyPaginationToBuilder(builder)

	sql, _, err := builder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - r.Builder: %w", err)
	}
	fmt.Println(sql)

	args := pgx.NamedArgs{
		"userId": userId,
	}
	rows, err := r.Pool.Query(ctx, sql, args)
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - r.Pool.Query: %w", err)
	}

	var ads []entity.Ad
	var count uint64
	for rows.Next() {
		var ad entity.Ad
		err = rows.Scan(
			&ad.Id,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Vin,
			&ad.IsTokenMinted,
			&ad.Brand,
			&ad.Model,
			&ad.YearOfRelease,
			&ad.IsFavorite,
			&count,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("AdRepository - List - rows.Scan: %w", err)
		}
		ads = append(ads, ad)
	}

	return ads, count, nil
}

func (r *AdRepository) HandleFavorite(ctx context.Context, adId, userId int64) error {
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
			return fmt.Errorf("AdRepository - HandleFavorite - r.Builder: %w", err)
		}

		_, err = r.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("AdRepository - HandleFavorite - r.Pool.Exec: %w", err)
		}

		return nil
	}

	sql, args, err := r.Builder.
		Insert("user_favorites").
		Columns("user_id", "ad_id").
		Values(userId, adId).
		ToSql()
	if err != nil {
		return fmt.Errorf("AdRepository - HandleFavorite - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AdRepository - HandleFavorite - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *AdRepository) isFavorite(ctx context.Context, adId, userId int64) (bool, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("user_favorites").
		Where(squirrel.Eq{"user_id": userId, "ad_id": adId}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("AdRepository - isFavorite - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("AdRepository - isFavorite - r.Pool.Query: %v", err)
	}

	rowsProcessed := 0
	for rows.Next() {
		rowsProcessed++
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("AdRepository - isFavorite - rows.Err: %v", err)
		return false, err
	}

	if rowsProcessed == 0 {
		fmt.Println("AdRepository - isFavorite - rows == 0")
		return false, nil
	}

	fmt.Println("AdRepository - isFavorite - rows > 0")
	return true, nil
}
