package repo

import (
	"github.com/Masterminds/squirrel"
	"github.com/golanguzb70/system-design-openbudget/internal/entity"
)

func PrepareFilter(filters []entity.Filter) squirrel.And {
	where := squirrel.And{}

	for _, e := range filters {
		switch e.Type {
		case "eq":
			where = append(where, squirrel.Eq{e.Column: e.Value})
		case "neq":
			where = append(where, squirrel.NotEq{e.Column: e.Value})
		case "gt":
			where = append(where, squirrel.Gt{e.Column: e.Value})
		case "gte":
			where = append(where, squirrel.GtOrEq{e.Column: e.Value})
		case "lt":
			where = append(where, squirrel.Lt{e.Column: e.Value})
		case "lte":
			where = append(where, squirrel.LtOrEq{e.Column: e.Value})
		case "search":
			where = append(where, squirrel.ILike{e.Column: "%" + e.Value + "%"})
		}
	}

	return where
}

func PrepareGetListQuery(selectQuery squirrel.SelectBuilder, filterRequest entity.GetListFilter) (query squirrel.SelectBuilder, where squirrel.And) {
	where = PrepareFilter(filterRequest.Filters)

	selectQuery = selectQuery.Where(where)

	for _, e := range filterRequest.OrderBy {
		selectQuery = selectQuery.OrderBy(e.Column + " " + e.Order)
	}

	if filterRequest.Limit <= 0 {
		filterRequest.Limit = 10
	}

	if filterRequest.Page <= 0 {
		filterRequest.Page = 1
	}

	selectQuery = selectQuery.Limit(uint64(filterRequest.Limit)).Offset(uint64((filterRequest.Page - 1) * filterRequest.Limit))

	return selectQuery, where
}
