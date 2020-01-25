package helpers

import (
	"../structs/fs"
	"../structs/paginator"
	"github.com/go-pg/pg/orm"
)

func SetFSAndPG(q *orm.Query, fs fs.FS, pg paginator.Paginator) *orm.Query {
	filter := fs.FilterBy
	v := fs.FilterValue
	sortBy := fs.SortBy
	orderBy := fs.OrderBy
	limit := pg.Limit
	offset := pg.Offset

	if filter != nil && v != nil {
		fstr := filter.(string)
		vstr := v.(string)
		if len(fstr) > 0 && len(vstr) > 0 {
			q.Where(Join(fstr, " = ?"), vstr)
		}
	}

	if sortBy != nil && orderBy != nil {
		q.Order(Join(fs.SortBy.(string), " ", fs.OrderBy.(string)))
	}

	if limit != nil {
		q.Limit(limit.(int))
	}

	if offset != nil {
		q.Offset(offset.(int))
	}

	return q
}
