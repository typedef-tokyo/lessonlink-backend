package utility

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/constants"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

func GetSessionData(c echo.Context) (vo.UserID, vo.RoleKey, error) {

	userID, ok := c.Get(constants.USER_IDENTIFIER).(vo.UserID)
	if !ok || !userID.IsValid() {
		return vo.USER_ID_INVALID, vo.ROLE_KEY_INVALID, log.WrapErrorWithStackTrace(errors.New("ユーザーIDが取得できません"))
	}

	roleKey, ok := c.Get(constants.ROLE_IDENTIFIER).(vo.RoleKey)
	if !ok || !roleKey.IsValid() {
		return vo.USER_ID_INVALID, vo.ROLE_KEY_INVALID, log.WrapErrorWithStackTrace(errors.New("ロールIDが取得できません"))
	}

	return userID, roleKey, nil
}
