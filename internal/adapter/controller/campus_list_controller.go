package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	ICampusListController interface {
		Execute(c echo.Context) error
	}

	CampusListController struct {
		usecase   usecase.ICampusListInputPort
		presenter presenter.ICampusListPresenter
		logger    ILogWriter
	}
)

func NewCampusListController(
	usecase usecase.ICampusListInputPort,
	presenter presenter.ICampusListPresenter,
	logger ILogWriter,
) ICampusListController {
	return &CampusListController{
		usecase:   usecase,
		presenter: presenter,
		logger:    logger,
	}
}

// @Summary キャンパスリスト取得
// @Description
// @Produce json
// @Success 200 {object} presenter.CampusListResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /campus/list [get]
func (h *CampusListController) Execute(c echo.Context) error {

	result, err := h.usecase.Execute(c.Request().Context())

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
