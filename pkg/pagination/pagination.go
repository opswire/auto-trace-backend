package pagination

import "github.com/Masterminds/squirrel"

const (
	DefaultPage    = 0
	DefaultPerPage = 10
	MinPerPage     = 1
	MaxPerPage     = 100
)

type Params struct {
	PerPage uint64
	Page    uint64
}

func (p *Params) ApplyPaginationToBuilder(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return builder.
		Limit(p.PerPage).
		Offset(p.Page)
}

type ListRange struct {
	PerPage uint64 `json:"per_page"`
	Page    uint64 `json:"page"`
	Count   uint64 `json:"count"`
}

func NewListRange(params Params, count uint64) ListRange {
	return ListRange{
		PerPage: params.PerPage,
		Page:    params.Page,
		Count:   count,
	}
}
