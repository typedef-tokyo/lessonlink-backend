package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IScheduleGetController interface {
		Execute(c echo.Context) error
	}

	ScheduleGetController struct {
		inputPort usecase.IScheduleGetInputPort
		presenter presenter.IScheduleGetPresenter
		logger    ILogWriter
	}
)

func NewScheduleGetController(
	inputPort usecase.IScheduleGetInputPort,
	presenter presenter.IScheduleGetPresenter,
	logger ILogWriter,
) IScheduleGetController {
	return &ScheduleGetController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

// @Summary スケジュールリスト取得
// @Description
// @Produce json
// @Param schedule_id path int true "ScheduleID"
// @Param history query int false "履歴番号"
// @Success 200 {object} presenter.ScheduleGetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/{schedule_id} [get]
func (h *ScheduleGetController) Execute(c echo.Context) error {

	scheduleID, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "スケジュールIDが不正です",
		})
	}

	historyIndex := 0
	paramHistoryIndex := c.QueryParam("history")
	if paramHistoryIndex != "" {

		inputHistoryIndex, err := strconv.Atoi(paramHistoryIndex)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"msg": "履歴番号が不正です",
			})
		}

		historyIndex = inputHistoryIndex
	}

	result, err := h.inputPort.Execute(c.Request().Context(), scheduleID, historyIndex)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
