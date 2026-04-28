package controller

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/logger"
)

type (
	IUserLogoutController interface {
		Execute(c echo.Context) error
	}

	UserLogoutController struct {
		env                 configs.EnvConfig
		dbConnection        *sql.DB
		facadeSessionDelete external.ISessionDeleteFacade
		logger              logger.ILogWriter
	}
)

func NewUserLogoutController(
	env configs.EnvConfig,
	dbConnection rdb.IMySQL,
	facadeSessionDelete external.ISessionDeleteFacade,
	logger logger.ILogWriter,
) IUserLogoutController {
	return &UserLogoutController{
		env:                 env,
		dbConnection:        dbConnection.GetConn(),
		facadeSessionDelete: facadeSessionDelete,
		logger:              logger,
	}
}

// @Summary ユーザーログアウト
// @Description
// @Produce json
// @Success 204 {string} string ""
// @Failure 401 {object} map[string]string
// @Failure 500 {object} string
// @Router /user/logout [post]
func (h *UserLogoutController) Execute(c echo.Context) error {

	userID, _, err := session_util.GetSessionData(c)
	if err != nil {
		err := log.WrapErrorWithStackTraceInternalServerError(log.Errorf("ログイン情報が確認できません"))
		h.logger.WriteErrLog(c, err)
		return c.NoContent(http.StatusNoContent)
	}

	err = h.facadeSessionDelete.Execute(c.Request().Context(), userID)
	if err != nil {
		h.logger.WriteErrLog(c, log.WrapErrorWithStackTrace(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	sameSite := http.SameSiteStrictMode
	if os.Getenv("LOCAL_TEST") == "true" {
		sameSite = http.SameSiteNoneMode
	}

	clearCookie := &http.Cookie{
		Name:     h.env.SessionName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Domain:   "",
		Secure:   os.Getenv("LOCAL_TEST") != "true",
		HttpOnly: true,
		SameSite: sameSite,
	}

	c.SetCookie(clearCookie)
	return c.NoContent(http.StatusNoContent)

}
