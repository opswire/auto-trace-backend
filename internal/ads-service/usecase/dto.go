package usecase

import "car-sell-buy-system/pkg/sqlutil"

type BasicListRequestDTO struct { // todo refactoring
	Filter     sqlutil.FiltersRequest
	Sort       sqlutil.SortsRequest
	Pagination sqlutil.Pagination
}

type AdListRequestDto struct {
	BasicListRequestDTO
}
