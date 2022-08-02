package rgdb

import (
	"context"
	"fmt"
	"github.com/juicyluv/rgdb/pkg/rgdberr"
	"github.com/juicyluv/rgdb/pkg/rgdbmsg"
)

//language=PostgreSQL
const addAchievementToUserQuery = `
	select 
		add_achievement_to_user
	from core.add_achievement_to_user(
	  _invoker_id := $1, 
	  _achievement_id := $2,
	  _user_id := $3
	)
`

func (c *Client) AddAchievementToUser(ctx context.Context, request *rgdbmsg.AddAchievementToUserRequest) error {
	row, err := c.Driver.Query(
		ctx,
		addAchievementToUserQuery,

		request.InvokerId,
		request.AchievementId,
		request.UserId,
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
