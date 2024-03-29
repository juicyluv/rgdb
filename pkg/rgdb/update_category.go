package rgdb

import (
	"context"
	"fmt"
	"github.com/juicyluv/rgdb/pkg/rgdberr"
	"github.com/juicyluv/rgdb/pkg/rgdbmsg"
)

//language=PostgreSQL
const updateCategoryQuery = `
	select 
		update_category
	from core.update_category(
	  _invoker_id := $1,
	  _category_id := $2,
	  _name := $3
	)
`

func (c *Client) UpdateCategory(ctx context.Context, request *rgdbmsg.UpdateCategoryRequest) error {
	row, err := c.Driver.Query(
		ctx,
		updateCategoryQuery,

		request.InvokerId,
		request.CategoryId,
		request.Name,
	)

	if err != nil {
		return fmt.Errorf(`%w: %v`, rgdberr.ErrInternal, err)
	}

	defer row.Close()

	if !row.Next() {
		if err = row.Err(); err != nil {
			return fmt.Errorf(`%w: %v`, rgdberr.ErrInternal, err)
		}

		return rgdberr.ErrInternal
	}

	var status []byte

	err = row.Scan(&status)

	if err != nil {
		return fmt.Errorf(`%w: %v`, rgdberr.ErrInternal, err)
	}

	if err = rgdberr.AnalyzeQueryStatus(status); err != nil {
		return err
	}

	return nil
}
