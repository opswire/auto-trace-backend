package sqlutil

import "github.com/Masterminds/squirrel"

const (
	_defaultPage    = 0
	_defaultPerPage = 10
)

type Pagination struct {
	PerPage int
	Page    int
}

func (p *Pagination) ApplyPagination(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	page := p.Page - 1
	perPage := p.PerPage

	if page < 0 {
		page = _defaultPage
	}

	if perPage == 0 {
		perPage = _defaultPerPage
	}

	return builder.
		Limit(uint64(perPage)).
		Offset(uint64(page * perPage))
}
