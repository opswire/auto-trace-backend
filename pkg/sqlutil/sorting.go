package sqlutil

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

type SortingDirection string

const (
	ASC  SortingDirection = "asc"
	DESC SortingDirection = "desc"
)

type SortsRequest map[string]string

type SortOption func(builder squirrel.SelectBuilder) squirrel.SelectBuilder

type Sortable interface {
	GetSorterOptionByField(field string, direction SortingDirection) (SortOption, error)
}

func ApplySorts(builder squirrel.SelectBuilder, sortModel Sortable, sortsRequest SortsRequest) (squirrel.SelectBuilder, error) {
	sortOptions, err := getSorterOptions(sortModel, sortsRequest)
	if err != nil {
		return squirrel.SelectBuilder{}, err
	}

	for _, sorter := range sortOptions {
		if sorter != nil {
			builder = sorter(builder)
		}
	}

	return builder, nil
}

func getSorterOptions(sortModel Sortable, sortsRequest SortsRequest) ([]SortOption, error) {
	var sorts []SortOption

	for field, direction := range sortsRequest {
		if direction != "asc" && direction != "desc" {
			return nil, fmt.Errorf("sorting direction is not valid")
		}

		sortOption, err := sortModel.GetSorterOptionByField(field, SortingDirection(direction))
		if err != nil {
			return nil, fmt.Errorf("[sqlutil.sorting] error when getting the filtering option: %w", err)
		}

		sorts = append(sorts, sortOption)
	}

	return sorts, nil
}
