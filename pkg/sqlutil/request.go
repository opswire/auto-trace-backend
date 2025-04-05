package sqlutil

type BasicListRequestDTO struct {
	Filter     FiltersRequest
	Sort       SortsRequest
	Pagination Pagination
}
