package postgres

import (
	"strings"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

func integerFilter(q *bun.SelectQuery, name string, fltr storage.IntegerFilter) *bun.SelectQuery {
	switch {
	case fltr.Between != nil:
		q.Where("? BETWEEN ? AND ?", bun.Safe(name), fltr.Between.From, fltr.Between.To)
	case fltr.Eq > 0:
		q.Where("? = ?", bun.Safe(name), fltr.Eq)
	case fltr.Neq > 0:
		q.Where("? != ?", bun.Safe(name), fltr.Neq)
	default:
		if fltr.Lte > 0 {
			q.Where("? <= ?", bun.Safe(name), fltr.Lte)
		}
		if fltr.Lt > 0 {
			q.Where("? < ?", bun.Safe(name), fltr.Lt)
		}
		if fltr.Gte > 0 {
			q.Where("? >= ?", bun.Safe(name), fltr.Gte)
		}
		if fltr.Gt > 0 {
			q.Where("? > ?", bun.Safe(name), fltr.Gt)
		}
	}

	return q
}

func timeFilter(q *bun.SelectQuery, name string, fltr storage.TimeFilter) *bun.SelectQuery {
	switch {
	case fltr.Between != nil:
		q.Where("extract(epoch from ?) BETWEEN ? AND ?", bun.Safe(name), fltr.Between.From, fltr.Between.To)
	default:
		if fltr.Lte > 0 {
			q.Where("extract(epoch from ?) <= ?", bun.Safe(name), fltr.Lte)
		}
		if fltr.Lt > 0 {
			q.Where("extract(epoch from ?) < ?", bun.Safe(name), fltr.Lt)
		}
		if fltr.Gte > 0 {
			q.Where("extract(epoch from ?) >= ?", bun.Safe(name), fltr.Gte)
		}
		if fltr.Gt > 0 {
			q.Where("extract(epoch from ?) > ?", bun.Safe(name), fltr.Gt)
		}
	}

	return q
}

func enumFilter(q *bun.SelectQuery, name string, fltr storage.EnumFilter) *bun.SelectQuery {
	switch {
	case fltr.Eq > 0:
		q.Where("? = ?", bun.Safe(name), fltr.Eq)
	case fltr.Neq > 0:
		q.Where("? != ?", bun.Safe(name), fltr.Neq)
	case len(fltr.In) > 0:
		q.Where("? IN (?)", bun.Safe(name), bun.In(fltr.In))
	case len(fltr.Notin) > 0:
		q.Where("? NOT IN (?)", bun.Safe(name), bun.In(fltr.Notin))
	}
	return q
}

func enumStringFilter(q *bun.SelectQuery, name string, fltr storage.EnumStringFilter) *bun.SelectQuery {
	switch {
	case fltr.Eq != "":
		q.Where("? = ?", bun.Safe(name), fltr.Eq)
	case fltr.Neq != "":
		q.Where("? != ?", bun.Safe(name), fltr.Neq)
	case len(fltr.In) > 0:
		q.Where("? IN (?)", bun.Safe(name), bun.In(fltr.In))
	case len(fltr.Notin) > 0:
		q.Where("? NOT IN (?)", bun.Safe(name), bun.In(fltr.Notin))
	}
	return q
}

func stringFilter(q *bun.SelectQuery, name string, fltr storage.StringFilter) *bun.SelectQuery {
	switch {
	case fltr.Eq != "":
		q.Where("? = ?", bun.Safe(name), fltr.Eq)
	case len(fltr.In) > 0:
		q.Where("? IN (?)", bun.Safe(name), bun.In(fltr.In))
	}

	return q
}

func equalityFilter(q *bun.SelectQuery, name string, fltr storage.EqualityFilter) *bun.SelectQuery {
	switch {
	case fltr.Eq != "":
		q.Where("? = ?", bun.Safe(name), fltr.Eq)
	case fltr.Neq != "":
		q.Where("? != ?", bun.Safe(name), fltr.Neq)
	}
	return q
}

func addressFilter(q *bun.SelectQuery, name string, fltr storage.BytesFilter, joinColumn string) *bun.SelectQuery {
	if name == "" || joinColumn == "" {
		return q
	}

	switch {
	case len(fltr.Eq) > 0:
		q = q.Relation(joinColumn)
		q = q.Where("?.? = ?", bun.Safe(joinColumn), bun.Safe(name), fltr.Eq)
	case len(fltr.In) > 0:
		q = q.Relation(joinColumn)
		q = q.Where("?.? IN (?)", bun.Safe(joinColumn), bun.Safe(name), bun.In(fltr.In))
	}

	return q
}

func idFilter(q *bun.SelectQuery, name string, fltr storage.IdFilter, joinColumn string) *bun.SelectQuery {
	if name == "" || joinColumn == "" {
		return q
	}

	switch {
	case fltr.Eq > 0:
		q = q.Relation(joinColumn)
		q = q.Where("? = ?", bun.Safe(name), fltr.Eq)
	case len(fltr.In) > 0:
		q = q.Relation(joinColumn)
		q = q.Where("? IN (?)", bun.Safe(name), bun.In(fltr.In))
	}

	return q
}

func jsonFilter(q *bun.SelectQuery, name string, fltr map[string]string) *bun.SelectQuery {
	if len(fltr) > 0 {
		q.Where("? is not null", bun.Ident(name))
	}

	for key, value := range fltr {
		builder := new(strings.Builder)
		builder.WriteString(name)
		path := strings.Split(key, ".")
		for i := range path {
			if i == len(path)-1 {
				builder.WriteString("->>")
			} else {
				builder.WriteString("->")
			}
			builder.WriteByte('\'')
			builder.WriteString(path[i])
			builder.WriteByte('\'')
		}
		builder.WriteByte('=')
		builder.WriteByte('\'')
		builder.WriteString(value)
		builder.WriteByte('\'')
		q.Where(builder.String())
	}
	return q
}

func addLimit(q *bun.SelectQuery, limit int) *bun.SelectQuery {
	if limit == 0 {
		return q
	}
	return q.Limit(limit)
}

func addOffset(q *bun.SelectQuery, offset int) *bun.SelectQuery {
	if offset == 0 {
		return q
	}
	return q.Offset(offset)
}

func addSort(q *bun.SelectQuery, field string, order sdk.SortOrder) *bun.SelectQuery {
	if field == "" {
		return q
	}
	if order == sdk.SortOrderAsc {
		return q.OrderExpr("(?) asc", bun.Ident(field))
	}
	return q.OrderExpr("(?) desc", bun.Ident(field))
}

func optionsFilter(q *bun.SelectQuery, tableName string, opts ...storage.FilterOption) *bun.SelectQuery {
	var opt storage.FilterOptions
	for i := range opts {
		opts[i](&opt)
	}
	q = addLimit(q, opt.Limit)
	q = addOffset(q, opt.Offset)

	if len(opt.SortFields) > 0 {
		q.Order(opt.SortFields...)
	} else {
		q = addSort(q, opt.SortField, opt.SortOrder)
	}

	if opt.MaxHeight > 0 {
		q = q.Where("?.? <= ?", bun.Ident(tableName), bun.Safe(opt.HeightColumnName), opt.MaxHeight)
	}
	if opt.Cursor > 0 {
		q = q.Where("?.id > ?", bun.Ident(tableName), opt.Cursor)
	}

	return q
}
