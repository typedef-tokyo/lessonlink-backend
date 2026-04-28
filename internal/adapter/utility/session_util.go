package utility

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/constants"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

const (
	INVALID_USER_ID  = 0
	INVALID_ROLE_KEY = ""
)

func GetSessionData(c echo.Context) (int, string, error) {

	userID, ok := c.Get(constants.USER_IDENTIFIER).(int)
	if !ok || userID == INVALID_USER_ID {
		return INVALID_USER_ID, INVALID_ROLE_KEY, log.WrapErrorWithStackTrace(errors.New("ユーザーIDが取得できません"))
	}

	roleKey, ok := c.Get(constants.ROLE_IDENTIFIER).(string)
	if !ok || roleKey == INVALID_ROLE_KEY {
		return INVALID_USER_ID, INVALID_ROLE_KEY, log.WrapErrorWithStackTrace(errors.New("ロールIDが取得できません"))
	}

	return userID, roleKey, nil
}
