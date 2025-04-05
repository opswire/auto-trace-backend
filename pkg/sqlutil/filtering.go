package sqlutil

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

type FiltersRequest map[string]string

type FilterOption func(builder squirrel.SelectBuilder) squirrel.SelectBuilder

type Filterable interface {
	GetFilterOptionByField(field, value string) (FilterOption, error)
}

func ApplyFilters(builder squirrel.SelectBuilder, filterModel Filterable, filtersRequest FiltersRequest) (squirrel.SelectBuilder, error) {
	filterOptions, err := getFilterOptions(filterModel, filtersRequest)
	if err != nil {
		return squirrel.SelectBuilder{}, err
	}

	for _, filter := range filterOptions {
		if filter != nil {
			builder = filter(builder)
		}
	}

	return builder, nil
}

func getFilterOptions(filterModel Filterable, filtersRequest FiltersRequest) ([]FilterOption, error) {
	var filters []FilterOption

	for field, value := range filtersRequest {
		filterOption, err := filterModel.GetFilterOptionByField(field, value)
		if err != nil {
			return nil, fmt.Errorf("[sqlutil.filtering] error when getting the filtering option: %w", err)
		}

		filters = append(filters, filterOption)
	}

	return filters, nil
}
