package psql

import (
	"car-sell-buy-system/internal/ads-service/domain/ad"
	"car-sell-buy-system/internal/ads-service/repository/psql/filter"
	"car-sell-buy-system/internal/ads-service/repository/psql/sort"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

const (
	AdTableName            = "ads"
	UserFavoritesTableName = "user_favorites"
)

type AdRepository struct {
	*postgres.Postgres
}

func NewAdRepository(pg *postgres.Postgres) *AdRepository {
	return &AdRepository{
		pg,
	}
}

func (r *AdRepository) GetById(ctx context.Context, id int64) (ad.Ad, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	selectBuilder := r.Builder.
		Select(
			// ad
			sqlutil.TableColumn(AdTableName, "id"),
			sqlutil.TableColumn(AdTableName, "title"),
			sqlutil.TableColumn(AdTableName, "description"),
			sqlutil.TableColumn(AdTableName, "price"),
			sqlutil.TableColumn(AdTableName, "vin"),
			sqlutil.TableColumn(AdTableName, "brand"),
			sqlutil.TableColumn(AdTableName, "model"),
			sqlutil.TableColumn(AdTableName, "year_of_release"),
			sqlutil.TableColumn(AdTableName, "image_url"),
			sqlutil.TableColumn(AdTableName, "user_id"),
			sqlutil.TableColumn(AdTableName, "created_at"),
			sqlutil.TableColumn(AdTableName, "updated_at"),
			sqlutil.TableColumn(AdTableName, "category"),
			sqlutil.TableColumn(AdTableName, "reg_number"),
			sqlutil.TableColumn(AdTableName, "type"),
			sqlutil.TableColumn(AdTableName, "color"),
			sqlutil.TableColumn(AdTableName, "hp"),
			sqlutil.TableColumn(AdTableName, "full_weight"),
			sqlutil.TableColumn(AdTableName, "solo_weight"),
		).
		Column(squirrel.Alias(squirrel.Case().When(squirrel.Eq{"nfts.is_minted": "true"}, "true").Else("false"), "is_favorite"))

	if userId != nil {
		selectBuilder = selectBuilder.Column(
			squirrel.Alias(
				squirrel.
					Case().
					When(
						squirrel.
							Select("1").
							Prefix("EXISTS (").
							From("chats").
							Where(squirrel.And{
								squirrel.Eq{"chats.buyer_id": userId},
								squirrel.Expr("chats.seller_id = ads.user_id"),
								squirrel.Expr("chats.ad_id = ads.id"),
							}).
							Suffix(")"),
						"true",
					).
					Else("false"),
				"chat_exists",
			),
		)
	} else {
		selectBuilder = selectBuilder.Column("false AS chat_exists")
	}

	sql, args, err := selectBuilder.
		Column("p1.status").
		Column("p1.expires_at").
		Column("p1.tariff_id").
		From(AdTableName).
		LeftJoin("payments p1 on p1.ad_id = ads.id").
		LeftJoin("nfts on nfts.vin = ads.vin").
		JoinClause("LEFT OUTER JOIN payments p2 ON (ads.id = p2.ad_id AND (p1.expires_at < p2.expires_at OR (p1.expires_at = p2.expires_at AND p1.payment_id < p2.payment_id)))").
		Where("p2.payment_id IS NULL").
		Where(squirrel.Eq{sqlutil.TableColumn(AdTableName, "id"): id}).
		ToSql()
	if err != nil {
		return ad.Ad{}, fmt.Errorf("AdRepository - GetById - r.Builder: %w", err)
	}
	fmt.Println("sql: ", sql)
	fmt.Println("args: ", args)

	var adv ad.Ad
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&adv.Id,
			&adv.Title,
			&adv.Description,
			&adv.Price,
			&adv.Vin,
			&adv.Brand,
			&adv.Model,
			&adv.YearOfRelease,
			&adv.ImageUrl,
			&adv.UserId,
			&adv.CreatedAt,
			&adv.UpdatedAt,
			&adv.Category,
			&adv.RegNumber,
			&adv.Type,
			&adv.Color,
			&adv.Hp,
			&adv.FullWeight,
			&adv.SoloWeight,
			&adv.IsTokenMinted,
			&adv.ChatExists,
			&adv.Promotion.Status,
			&adv.Promotion.ExpiresAt,
			&adv.Promotion.TariffId,
		)
	if err != nil {
		return ad.Ad{}, fmt.Errorf("AdRepository - GetById - row.Scan: %w", err)
	}

	return adv, nil
}

func (r *AdRepository) Store(ctx context.Context, dto ad.StoreDTO) (ad.Ad, error) {
	sql, args, err := r.Builder.
		Insert(AdTableName).
		Columns(
			"title",
			"description",
			"price",
			"vin",
			"brand",
			"model",
			"year_of_release",
			"is_token_minted",
			"image_url",
			"category",
			"reg_number",
			"type",
			"color",
			"hp",
			"full_weight",
			"solo_weight",
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
			dto.CurrentImageUrl,
			dto.Category,
			dto.RegNumber,
			dto.Type,
			dto.Color,
			dto.Hp,
			dto.FullWeight,
			dto.SoloWeight,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return ad.Ad{}, fmt.Errorf("AdRepository - Store - r.Builder: %w", err)
	}

	adv := ad.Ad{
		Title:         dto.Title,
		Description:   dto.Description,
		Price:         dto.Price,
		Vin:           dto.Vin,
		Brand:         dto.Brand,
		Model:         dto.Model,
		YearOfRelease: dto.YearOfRelease,
		IsTokenMinted: false,
		IsFavorite:    false,
		ImageUrl:      dto.Image.Path,
	}
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&adv.Id)
	if err != nil {
		return ad.Ad{}, fmt.Errorf("AdRepository - Store - row.Scan: %w", err)
	}

	return adv, nil
}

func (r *AdRepository) Update(ctx context.Context, id int64, dto ad.UpdateDTO) error {
	sql, args, err := r.Builder.
		Update(AdTableName).
		Set("title", dto.Title).
		Set("description", dto.Description).
		Set("price", dto.Price).
		Set("vin", dto.Vin).
		Set("brand", dto.Brand).
		Set("model", dto.Model).
		Set("year_of_release", dto.YearOfRelease).
		Set("image_url", dto.CurrentImageUrl).
		Set("category", dto.Category).
		Set("reg_number", dto.RegNumber).
		Set("type", dto.Type).
		Set("color", dto.Color).
		Set("hp", dto.Hp).
		Set("full_weight", dto.FullWeight).
		Set("solo_weight", dto.SoloWeight).
		Where(squirrel.Eq{sqlutil.TableColumn(AdTableName, "id"): id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("AdRepository - Update - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AdRepository - Store - row.Scan: %w", err)
	}

	return nil
}

func (r *AdRepository) List(ctx context.Context, dto ad.ListDTO) ([]ad.Ad, uint64, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	builder := r.Builder.
		Select(
			sqlutil.TableColumn(AdTableName, "id"),
			sqlutil.TableColumn(AdTableName, "title"),
			sqlutil.TableColumn(AdTableName, "description"),
			sqlutil.TableColumn(AdTableName, "price"),
			sqlutil.TableColumn(AdTableName, "vin"),
			sqlutil.TableColumn(AdTableName, "brand"),
			sqlutil.TableColumn(AdTableName, "model"),
			sqlutil.TableColumn(AdTableName, "year_of_release"),
			sqlutil.TableColumn(AdTableName, "image_url"),
			sqlutil.TableColumn(AdTableName, "user_id"),
			sqlutil.TableColumn(AdTableName, "created_at"),
			sqlutil.TableColumn(AdTableName, "updated_at"),
			sqlutil.TableColumn(AdTableName, "category"),
			sqlutil.TableColumn(AdTableName, "reg_number"),
			sqlutil.TableColumn(AdTableName, "type"),
			sqlutil.TableColumn(AdTableName, "color"),
			sqlutil.TableColumn(AdTableName, "hp"),
			sqlutil.TableColumn(AdTableName, "full_weight"),
			sqlutil.TableColumn(AdTableName, "solo_weight"),
		).
		Column(squirrel.Alias(squirrel.Case().When(squirrel.Expr("user_favorites.user_id = ?", userId), "true").Else("false"), "is_favorite")).
		Column(squirrel.Alias(squirrel.Expr("COUNT(*) OVER()"), "total")).
		Column(squirrel.Alias(squirrel.Case().When(squirrel.Eq{"nfts.is_minted": "true"}, "true").Else("false"), "is_favorite")).
		Column("p1.status").
		Column("p1.expires_at").
		Column("p1.tariff_id").
		LeftJoin("payments p1 on p1.ad_id = ads.id").
		LeftJoin("nfts on nfts.vin = ads.vin").
		JoinClause("LEFT OUTER JOIN payments p2 ON (ads.id = p2.ad_id AND (p1.expires_at < p2.expires_at OR (p1.expires_at = p2.expires_at AND p1.payment_id < p2.payment_id)))").
		Where("p2.payment_id IS NULL").
		From(AdTableName).
		LeftJoin(UserFavoritesTableName + " ON " + sqlutil.TableColumn(UserFavoritesTableName, "ad_id") + " = " + sqlutil.TableColumn(AdTableName, "id")).
		From(AdTableName)

	builder, err := sqlutil.ApplyFilters(builder, &filter.AdFilter{}, dto.Filter)
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - sqlutil.ApplyFilters: %w", err)
	}

	builder = builder.
		OrderBy("p1.tariff_id DESC nulls last")

	builder, err = sqlutil.ApplySorts(builder, &sort.AdSorter{}, dto.Sort)
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - sqlutil.ApplySorts: %w", err)
	}

	builder = dto.Pagination.ApplyPaginationToBuilder(builder)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - r.Builder: %w", err)
	}
	fmt.Println("sql: ", sql)
	fmt.Println("args: ", args)

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - r.Pool.Query: %w", err)
	}

	var ads []ad.Ad
	var count uint64
	for rows.Next() {
		var adv ad.Ad
		err = rows.Scan(
			&adv.Id,
			&adv.Title,
			&adv.Description,
			&adv.Price,
			&adv.Vin,
			&adv.Brand,
			&adv.Model,
			&adv.YearOfRelease,
			&adv.ImageUrl,
			&adv.UserId,
			&adv.CreatedAt,
			&adv.UpdatedAt,
			&adv.Category,
			&adv.RegNumber,
			&adv.Type,
			&adv.Color,
			&adv.Hp,
			&adv.FullWeight,
			&adv.SoloWeight,
			&adv.IsFavorite,
			&count,
			&adv.IsTokenMinted,
			&adv.Promotion.Status,
			&adv.Promotion.ExpiresAt,
			&adv.Promotion.TariffId,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("AdRepository - List - rows.Scan: %w", err)
		}
		ads = append(ads, adv)
	}

	return ads, count, nil
}

func (r *AdRepository) Delete(ctx context.Context, id int64) error {
	sql, args, err := r.Builder.
		Delete(AdTableName).
		Where(squirrel.Eq{sqlutil.TableColumn(AdTableName, "id"): id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("AdRepository - Delete - r.Builder: %w", err)
	}

	if _, err = r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("AdRepository - Delete - r.Pool.Exec: %w", err)
	}

	return nil
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
