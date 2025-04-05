package filtermodel

import (
	"car-sell-buy-system/pkg/sqlutil"
	"fmt"
	"github.com/Masterminds/squirrel"
)

type AdFilter struct {
}

func (f *AdFilter) GetFilterOptionByField(field, value string) (sqlutil.FilterOption, error) {
	var filter sqlutil.FilterOption
	var err error

	switch field {
	case "title":
		filter, err = filterByTitle(value)
	case "brand":
		filter, err = filterByBrand(value)
	}
	if err != nil {
		return nil, err
	}

	return filter, nil
}

func filterByTitle(title string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.Like{"ads.title": fmt.Sprintf("%%%s%%", title)})
	}, nil
}

func filterByBrand(brand string) (sqlutil.FilterOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.Where(squirrel.Like{"ads.brand": fmt.Sprintf("%%%s%%", brand)})
	}, nil
}
