package controller

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/repository"
)

type (
	IUserLogoutController interface {
		Execute(c echo.Context) error
	}

	UserLogoutController struct {
		env               configs.EnvConfig
		dbConnection      *sql.DB
		repositorySession repository.SessionRepository
		logger            ILogWriter
	}
)

func NewUserLogoutController(
	env configs.EnvConfig,
	dbConnection rdb.IMySQL,
	repositorySession repository.SessionRepository,
	logger ILogWriter,
) IUserLogoutController {
	return &UserLogoutController{
		env:               env,
		dbConnection:      dbConnection.GetConn(),
		repositorySession: repositorySession,
		logger:            logger,
	}
}

// @Summary ユーザーログアウト
// @Description
// @Produce json
// @Success 204 {string} string ""
// @Failure 500 {object} string
// @Router /user/logout [post]
func (h *UserLogoutController) Execute(c echo.Context) error {

	userID, _, err := session_util.GetSessionData(c)
	if err != nil {
		err := log.WrapErrorWithStackTraceInternalServerError(log.Errorf("ログイン情報が確認できません"))
		h.logger.WriteErrLog(c, err)
		return c.NoContent(http.StatusNoContent)
	}

	tx, err := h.dbConnection.BeginTx(c.Request().Context(), nil)
	if err != nil {
		h.logger.WriteErrLog(c, log.WrapErrorWithStackTraceInternalServerError(err))
		return c.NoContent(http.StatusNoContent)
	}
	defer tx.Rollback()

	err = h.repositorySession.Delete(c.Request().Context(), tx, userID)
	if err != nil {
		h.logger.WriteErrLog(c, log.WrapErrorWithStackTrace(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	tx.Commit()

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
