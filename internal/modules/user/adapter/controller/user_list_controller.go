package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/usecase/interactor"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/logger"
)

type (
	IUserListController interface {
		Execute(c echo.Context) error
	}

	UserListController struct {
		inputPort interactor.IUserListInputPort
		presenter presenter.IUserListPresenter
		logger    logger.ILogWriter
	}
)

func NewUserListController(
	inputPort interactor.IUserListInputPort,
	presenter presenter.IUserListPresenter,
	logger logger.ILogWriter,
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

	result, err := h.inputPort.Execute(c.Request().Context(), roleKey, userID)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
