package rgdb

import (
	"context"
	"fmt"
	"github.com/juicyluv/rgdb/pkg/rgdberr"
	"github.com/juicyluv/rgdb/pkg/rgdbmsg"
)

//language=PostgreSQL
const removeAchievementFromUserQuery = `
	select 
		remove_achievement_from_user
	from core.remove_achievement_from_user(
	  _invoker_id := $1, 
	  _achievement_id := $2,
	  _user_id := $3
	)
`

func (c *Client) RemoveAchievementFromUser(ctx context.Context, request *rgdbmsg.RemoveAchievementFromUserRequest) error {
	row, err := c.Driver.Query(
		ctx,
		removeAchievementFromUserQuery,

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
