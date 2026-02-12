package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IUserListController interface {
		Execute(c echo.Context) error
	}

	UserListController struct {
		inputPort usecase.IUserListInputPort
		presenter presenter.IUserListPresenter
		logger    ILogWriter
	}
)

func NewUserListController(
	inputPort usecase.IUserListInputPort,
	presenter presenter.IUserListPresenter,
	logger ILogWriter,
) IUserListController {
	return &UserListController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

// @Summary ユーザー一覧取得
// @Description
// @Produce json
// @Success 200 {object} presenter.UserListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/list [get]
func (h *UserListController) Execute(c echo.Context) error {

	userID, roleKey, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
	}

	result, err := h.inputPort.Execute(c.Request().Context(), userID, roleKey)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
