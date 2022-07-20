package rgdb

import (
	"context"
	"fmt"
	"rgdb/rgdberr"
	"rgdb/rgdbmsg"
)

//language=PostgreSQL
const removeBookReviewQuery = `
	select 
		remove_book_review
	from core.remove_book_review(
	  _invoker_id := $1, 
	  _review_id := $2
	)
`

func (d *driver) RemoveBookReview(ctx context.Context, request *rgdbmsg.RemoveBookReviewRequest) error {
	row, err := d.pool.Query(
		ctx,
		removeBookReviewQuery,

		request.InvokerId,
		request.ReviewId,
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
