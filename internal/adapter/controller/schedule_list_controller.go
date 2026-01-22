package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	schedulelist "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/schedule_list"
)

type (
	IScheduleListController interface {
		Execute(c echo.Context) error
	}

	ScheduleListController struct {
		inputPort schedulelist.IScheduleListQueryInputPort
		presenter presenter.IScheduleListPresenter
		logger    ILogWriter
	}
)

func NewScheduleListController(
	inputPort schedulelist.IScheduleListQueryInputPort,
	presenter presenter.IScheduleListPresenter,
	logger ILogWriter,
) IScheduleListController {
	return &ScheduleListController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

// @Summary スケジュールリスト取得
// @Description
// @Produce json
// @Param campus path string true "校舎"
// @Success 200 {object} presenter.ScheduleListResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /schedule/list/{campus} [get]
func (h *ScheduleListController) Execute(c echo.Context) error {

	campus := c.Param("campus")
	if campus == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	result, err := h.inputPort.Execute(c.Request().Context(), campus)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
