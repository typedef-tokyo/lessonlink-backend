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
	ILoginUserGetController interface {
		Execute(c echo.Context) error
	}

	LoginUserGetController struct {
		inputPort interactor.IUserGetInputPort
		presenter presenter.IUserGetPresenter
		logger    logger.ILogWriter
	}
)

func NewLoginUserGetController(
	inputPort interactor.IUserGetInputPort,
	presenter presenter.IUserGetPresenter,
	logger logger.ILogWriter,
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
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/self [get]
func (h *LoginUserGetController) Execute(c echo.Context) error {

	userID, _, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
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
