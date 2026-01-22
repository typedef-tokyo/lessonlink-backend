package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IUserDeleteController interface {
		Execute(c echo.Context) error
	}

	UserDeleteController struct {
		inputPort usecase.IUserDeleteInputPort
		presenter presenter.IUserDeletePresenter
		logger    ILogWriter
	}
)

func NewUserDeleteController(
	inputPort usecase.IUserDeleteInputPort,
	presenter presenter.IUserDeletePresenter,
	logger ILogWriter,
) IUserDeleteController {
	return &UserDeleteController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

// @Summary ユーザー削除
// @Description
// @Produce json
// @Param userid path string true "UserID"
// @Success 200 {object} presenter.UserDeleteResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /user/{userid} [delete]
func (h *UserDeleteController) Execute(c echo.Context) error {

	// セッション情報を取得
	userID, roleKey, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": err.Error(),
		})
	}

	userIDParam := c.Param("userid")
	userIDint, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "ユーザーIDの形式が不正です",
		})
	}

	err = h.inputPort.Execute(c.Request().Context(), roleKey, userIDint, userID)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present())
}
