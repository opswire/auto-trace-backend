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
	}
	if err != nil {
		return nil, err
	}

	return filter, nil
}

func filterByTitle(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.Like{"ads.title": fmt.Sprintf("%%%s%%", value)})
	}, nil
}

func filterByDescription(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.Like{"ads.description": "%" + value + "%"})
	}, nil
}

func filterByBrand(value string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.Like{"ads.brand": fmt.Sprintf("%%%s%%", value)})
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
