package postgres

import (
	"strings"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func integerFilter(q *orm.Query, name string, fltr storage.IntegerFilter) *orm.Query {
	switch {
	case fltr.Between != nil:
		q.Where("? BETWEEN ? AND ?", pg.Safe(name), fltr.Between.From, fltr.Between.To)
	case fltr.Eq > 0:
		q.Where("? = ?", pg.Safe(name), fltr.Eq)
	case fltr.Neq > 0:
		q.Where("? != ?", pg.Safe(name), fltr.Neq)
	default:
		if fltr.Lte > 0 {
			q.Where("? <= ?", pg.Safe(name), fltr.Lte)
		}
		if fltr.Lt > 0 {
			q.Where("? < ?", pg.Safe(name), fltr.Lt)
		}
		if fltr.Gte > 0 {
			q.Where("? >= ?", pg.Safe(name), fltr.Gte)
		}
		if fltr.Gt > 0 {
			q.Where("? > ?", pg.Safe(name), fltr.Gt)
		}
	}

	return q
}

func timeFilter(q *orm.Query, name string, fltr storage.TimeFilter) *orm.Query {
	switch {
	case fltr.Between != nil:
		q.Where("? BETWEEN ? AND ?", pg.Safe(name), fltr.Between.From, fltr.Between.To)
	default:
		if fltr.Lte > 0 {
			q.Where("? <= ?", pg.Safe(name), fltr.Lte)
		}
		if fltr.Lt > 0 {
			q.Where("? < ?", pg.Safe(name), fltr.Lt)
		}
		if fltr.Gte > 0 {
			q.Where("? >= ?", pg.Safe(name), fltr.Gte)
		}
		if fltr.Gt > 0 {
			q.Where("? > ?", pg.Safe(name), fltr.Gt)
		}
	}

	return q
}

func enumFilter(q *orm.Query, name string, fltr storage.EnumFilter) *orm.Query {
	switch {
	case fltr.Eq > 0:
		q.Where("? = ?", pg.Safe(name), fltr.Eq)
	case fltr.Neq > 0:
		q.Where("? != ?", pg.Safe(name), fltr.Neq)
	case len(fltr.In) > 0:
		q.Where("? IN (?)", pg.Safe(name), pg.In(fltr.In))
	case len(fltr.Notin) > 0:
		q.Where("? NOT IN (?)", pg.Safe(name), pg.In(fltr.Notin))
	}
	return q
}

func enumStringFilter(q *orm.Query, name string, fltr storage.EnumStringFilter) *orm.Query {
	switch {
	case fltr.Eq != "":
		q.Where("? = ?", pg.Safe(name), fltr.Eq)
	case fltr.Neq != "":
		q.Where("? != ?", pg.Safe(name), fltr.Neq)
	case len(fltr.In) > 0:
		q.Where("? IN (?)", pg.Safe(name), pg.In(fltr.In))
	case len(fltr.Notin) > 0:
		q.Where("? NOT IN (?)", pg.Safe(name), pg.In(fltr.Notin))
	}
	return q
}

func stringFilter(q *orm.Query, name string, fltr storage.StringFilter) *orm.Query {
	switch {
	case fltr.Eq != "":
		q.Where("? = ?", pg.Safe(name), fltr.Eq)
	case len(fltr.In) > 0:
		q.Where("? IN (?)", pg.Safe(name), pg.In(fltr.In))
	}

	return q
}

func equalityFilter(q *orm.Query, name string, fltr storage.EqualityFilter) *orm.Query {
	switch {
	case fltr.Eq != "":
		q.Where("? = ?", pg.Safe(name), fltr.Eq)
	case fltr.Neq != "":
		q.Where("? != ?", pg.Safe(name), fltr.Neq)
	}
	return q
}

func addressFilter(q *orm.Query, name string, fltr storage.BytesFilter, joinColumn string) *orm.Query {
	if name == "" || joinColumn == "" {
		return q
	}

	switch {
	case len(fltr.Eq) > 0:
		q = q.Relation(joinColumn)
		q = q.Where("?.? = ?", pg.Safe(joinColumn), pg.Safe(name), fltr.Eq)
	case len(fltr.In) > 0:
		q = q.Relation(joinColumn)
		q = q.Where("?.? IN (?)", pg.Safe(joinColumn), pg.Safe(name), pg.In(fltr.In))
	}

	return q
}

func idFilter(q *orm.Query, name string, fltr storage.IdFilter, joinColumn string) *orm.Query {
	if name == "" || joinColumn == "" {
		return q
	}

	switch {
	case fltr.Eq > 0:
		q = q.Relation(joinColumn)
		q = q.Where("? = ?", pg.Safe(name), fltr.Eq)
	case len(fltr.In) > 0:
		q = q.Relation(joinColumn)
		q = q.Where("? IN (?)", pg.Safe(name), pg.In(fltr.In))
	}

	return q
}

func jsonFilter(q *orm.Query, name string, fltr map[string]string) *orm.Query {
	if len(fltr) > 0 {
		q.Where("? is not null", pg.Ident(name))
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

func addLimit(q *orm.Query, limit int) *orm.Query {
	if limit == 0 {
		return q
	}
	return q.Limit(limit)
}

func addOffset(q *orm.Query, offset int) *orm.Query {
	if offset == 0 {
		return q
	}
	return q.Offset(offset)
}

func addSort(q *orm.Query, field string, order sdk.SortOrder) *orm.Query {
	if field == "" {
		return q
	}
	if order == sdk.SortOrderAsc {
		return q.OrderExpr("? asc", pg.Ident(field))
	}
	return q.OrderExpr("? desc", pg.Ident(field))
}

func optionsFilter(q *orm.Query, tableName string, opts ...storage.FilterOption) *orm.Query {
	var opt storage.FilterOptions
	for i := range opts {
		opts[i](&opt)
	}
	q = addLimit(q, opt.Limit)
	q = addOffset(q, opt.Offset)
	q = addSort(q, opt.SortField, opt.SortOrder)

	if opt.MaxHeight > 0 {
		q = q.Where("?.? <= ?", pg.Ident(tableName), pg.Safe(opt.HeightColumnName), opt.MaxHeight)
	}
	if opt.Cursor > 0 {
		q = q.Where("?.id > ?", pg.Ident(tableName), opt.Cursor)
	}

	return q
}
