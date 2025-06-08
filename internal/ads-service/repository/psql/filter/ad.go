package filter

import (
	"car-sell-buy-system/pkg/sqlutil"
	"fmt"
	"github.com/Masterminds/squirrel"
)

type AdFilter struct {
}

func (f *AdFilter) GetFilterOptionByField(field string, value string) (sqlutil.FilterOption, error) {
	var filter sqlutil.FilterOption
	var err error

	switch field {
	case "title":
		filter, err = filterByTitle(value)
	case "description":
		filter, err = filterByDescription(value)
	case "brand":
		filter, err = filterByBrand(value)
	case "is_favorite":
		filter, err = filterByFavorite(value)
	case "price_min":
		filter, err = filterByPriceMin(value)
	case "price_max":
		filter, err = filterByPriceMax(value)
	case "year_min":
		filter, err = filterByYearMin(value)
	case "year_max":
		filter, err = filterByYearMax(value)
	case "car_category":
		filter, err = filterByCarCategory(value)
	case "driver_category":
		filter, err = filterByDriverCategory(value)
	}
	if err != nil {
		return nil, err
	}

	return filter, nil
}

func filterByTitle(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.ILike{"ads.title": fmt.Sprintf("%%%s%%", value)})
	}, nil
}

func filterByDescription(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.ILike{"ads.description": "%" + value + "%"})
	}, nil
}

func filterByBrand(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.ILike{"ads.brand": fmt.Sprintf("%%%s%%", value)})
	}, nil
}

func filterByPriceMin(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.GtOrEq{"ads.price": value})
	}, nil
}

func filterByPriceMax(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.LtOrEq{"ads.price": value})
	}, nil
}

func filterByYearMin(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.GtOrEq{"ads.year_of_release": value})
	}, nil
}

func filterByYearMax(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.LtOrEq{"ads.year_of_release": value})
	}, nil
}

func filterByCarCategory(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.ILike{"ads.category": value})
	}, nil
}

func filterByDriverCategory(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.ILike{"ads.type": value})
	}, nil
}

func filterByFavorite(isFavorite string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		if isFavorite == "true" {
			return builder.Where("user_favorites.ad_id IS NOT NULL")
		}

		return builder
	}, nil
}
