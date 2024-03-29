package rgdb

import (
	"context"
	"fmt"
	"github.com/juicyluv/rgdb/pkg/rgdberr"
	"github.com/juicyluv/rgdb/pkg/rgdbmsg"
)

//language=PostgreSQL
const getAuthorQuery = `
	select 
		name,
		creator_id,
		creator_username,
		description,
		created_at,
		updated_at,
		error
	from core.get_author(
	  _author_id := $1
	)
`

func (c *Client) GetAuthor(ctx context.Context, request *rgdbmsg.GetAuthorRequest) (*rgdbmsg.Author, error) {
	row, err := c.Driver.Query(ctx, getAuthorQuery, request.AuthorId)

	if err != nil {
		return nil, fmt.Errorf(`%w: %v`, rgdberr.ErrInternal, err)
	}

	defer row.Close()

	if !row.Next() {
		if err = row.Err(); err != nil {
			return nil, fmt.Errorf(`%w: %v`, rgdberr.ErrInternal, err)
		}

		return nil, rgdberr.ErrInternal
	}

	var (
		status []byte
		author rgdbmsg.Author
	)

	err = row.Scan(
		&author.Name,
		&author.CreatorId,
		&author.CreatorUsername,
		&author.Description,
		&author.CreatedAt,
		&author.UpdatedAt,
		&status,
	)

	if err != nil {
		return nil, fmt.Errorf(`%w: %v`, rgdberr.ErrInternal, err)
	}

	if err = rgdberr.AnalyzeQueryStatus(status); err != nil {
		return nil, err
	}

	author.AuthorId = &request.AuthorId

	return &author, nil
}
