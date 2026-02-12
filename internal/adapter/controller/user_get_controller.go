package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IUserGetController interface {
		Execute(c echo.Context) error
	}

	UserGetController struct {
		inputPort usecase.IUserGetInputPort
		presenter presenter.IUserGetPresenter
		logger    ILogWriter
	}
)

func NewUserGetController(
	inputPort usecase.IUserGetInputPort,
	presenter presenter.IUserGetPresenter,
	logger ILogWriter,
) IUserGetController {
	return &UserGetController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

// @Summary ユーザー取得
// @Description
// @Produce json
// @Param userid path string true "UserID"
// @Success 200 {object} presenter.UserGetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/{userid} [get]
func (h *UserGetController) Execute(c echo.Context) error {

	_userID := c.Param("userid")
	userID, err := strconv.Atoi(_userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	result, err := h.inputPort.Execute(c, userID)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
