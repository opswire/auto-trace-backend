package sqlutil

import "car-sell-buy-system/pkg/pagination"

type BasicListRequestDTO struct {
	Filter     FiltersRequest
	Sort       SortsRequest
	Pagination pagination.Params
}
