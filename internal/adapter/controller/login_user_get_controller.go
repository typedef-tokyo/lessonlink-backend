package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	ILoginUserGetController interface {
		Execute(c echo.Context) error
	}

	LoginUserGetController struct {
		inputPort usecase.IUserGetInputPort
		presenter presenter.IUserGetPresenter
		logger    ILogWriter
	}
)

func NewLoginUserGetController(
	inputPort usecase.IUserGetInputPort,
	presenter presenter.IUserGetPresenter,
	logger ILogWriter,
) ILoginUserGetController {
	return &LoginUserGetController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

// @Summary ログインユーザー取得
// @Description
// @Produce json
// @Success 200 {object} presenter.UserGetResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /user/self [get]
func (h *LoginUserGetController) Execute(c echo.Context) error {

	userID, _, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
	}

	result, err := h.inputPort.Execute(c, userID.Value())

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
