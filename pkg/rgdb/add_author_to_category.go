package rgdb

import (
	"context"
	"fmt"
	rgdberr2 "github.com/juicyluv/rgdb/pkg/rgdberr"
	"github.com/juicyluv/rgdb/pkg/rgdbmsg"
)

//language=PostgreSQL
const addAuthorToCategoryQuery = `
	select 
		add_author_to_category
	from core.add_author_to_category(
	  _invoker_id := $1, 
	  _author_id := $2,
	  _category_id := $3
	)
`

func (d *driver) AddAuthorToCategory(ctx context.Context, request *rgdbmsg.AddAuthorToCategoryRequest) error {
	row, err := d.pool.Query(
		ctx,
		addAuthorToCategoryQuery,

		request.InvokerId,
		request.AuthorId,
		request.CategoryId,
	)

	if err != nil {
		return fmt.Errorf(`%w: %v`, rgdberr2.ErrInternal, err)
	}

	defer row.Close()

	if !row.Next() {
		if err = row.Err(); err != nil {
			return fmt.Errorf(`%w: %v`, rgdberr2.ErrInternal, err)
		}

		return rgdberr2.ErrInternal
	}

	var status []byte

	err = row.Scan(&status)

	if err != nil {
		return fmt.Errorf(`%w: %v`, rgdberr2.ErrInternal, err)
	}

	if err = rgdberr2.AnalyzeQueryStatus(status); err != nil {
		return err
	}

	return nil
}
