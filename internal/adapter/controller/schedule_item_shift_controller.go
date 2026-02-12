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
	IScheduleItemShiftController interface {
		Execute(c echo.Context) error
	}

	ScheduleItemShiftController struct {
		inputPort usecase.IScheduleItemShiftInputPort
		presenter presenter.IScheduleItemEditPresenter
		logger    ILogWriter
	}
)

func NewScheduleItemShiftController(
	inputPort usecase.IScheduleItemShiftInputPort,
	presenter presenter.IScheduleItemEditPresenter,
	logger ILogWriter,
) IScheduleItemShiftController {
	return &ScheduleItemShiftController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	ScheduleItemShiftRequestData struct {
		HistoryIndex int `json:"history_index"`
		RoomIndex    int `json:"room_index"`
	}
)

// @Summary スケジュール編集アイテムシフト
// @Description
// @Produce json
// @Param schedule_id path int true "ScheduleID"
// @Param request body ScheduleItemShiftRequestData true "アイテムシフトリクエスト"
// @Success 200 {object} presenter.ScheduleItemEditResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/{schedule_id}/item-shift [post]
func (h *ScheduleItemShiftController) Execute(c echo.Context) error {

	var err error

	// セッション情報を取得
	userID, role, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
	}

	scheduleID, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "スケジュールIDが不正です",
		})
	}

	var requestData ScheduleItemMoveRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	result, err := h.inputPort.Execute(c.Request().Context(), role, userID, scheduleID, requestData.HistoryIndex, requestData.RoomIndex)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
