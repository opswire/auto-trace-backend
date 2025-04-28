package sort

import (
	"car-sell-buy-system/pkg/sqlutil"
	"fmt"
	"github.com/Masterminds/squirrel"
)

type AdSorter struct {
}

func (f *AdSorter) GetSorterOptionByField(field string, direction sqlutil.SortingDirection) (sqlutil.SortOption, error) {
	var sorter sqlutil.SortOption
	var err error

	switch field {
	case "id":
		sorter, err = sortById(direction)
	case "title":
		sorter, err = sortByTitle(direction)
	case "price":
		sorter, err = sortByPrice(direction)
	}
	if err != nil {
		return nil, err
	}

	return sorter, nil
}

func sortById(direction sqlutil.SortingDirection) (sqlutil.SortOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.OrderBy(fmt.Sprintf("ads.id %s", direction))
	}, nil
}

func sortByTitle(direction sqlutil.SortingDirection) (sqlutil.SortOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.OrderBy(fmt.Sprintf("ads.title %s", direction))
	}, nil
}

func sortByPrice(direction sqlutil.SortingDirection) (sqlutil.SortOption, error) {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return builder.OrderBy(fmt.Sprintf("ads.price %s", direction))
	}, nil
}
