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
	IScheduleItemMoveController interface {
		Execute(c echo.Context) error
	}

	ScheduleItemMoveController struct {
		inputPort usecase.IScheduleItemMoveInputPort
		presenter presenter.IScheduleItemEditPresenter
		logger    ILogWriter
	}
)

func NewScheduleItemMoveController(
	inputPort usecase.IScheduleItemMoveInputPort,
	presenter presenter.IScheduleItemEditPresenter,
	logger ILogWriter,
) IScheduleItemMoveController {
	return &ScheduleItemMoveController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	ScheduleItemMoveRequestData struct {
		HistoryIndex    int    `json:"history_index"`
		LessonID        int    `json:"lesson_id"`
		ItemTag         string `json:"item_tag"`
		Identifier      string `json:"identifier"`
		Duration        int    `json:"duration"`
		StartTimeHour   int    `json:"start_time_hour"`
		StartTimeMinute int    `json:"start_time_minute"`
		EndTimeHour     int    `json:"end_time_hour"`
		EndTimeMinutes  int    `json:"end_time_minutes"`
		RoomIndex       int    `json:"room_index"`
	}
)

// @Summary スケジュール編集アイテム移動
// @Description
// @Produce json
// @Param schedule_id path string true "ScheduleID"
// @Param request body ScheduleItemMoveRequestData true "アイテム移動リクエスト"
// @Success 200 {object} presenter.ScheduleItemEditResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /schedule/{schedule_id}/item-move [post]
func (h *ScheduleItemMoveController) Execute(c echo.Context) error {

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

	result, err := h.inputPort.Execute(c.Request().Context(), role, userID, scheduleID, requestData.HistoryIndex, usecase.ScheduleItemMoveInput{
		LessonID:        requestData.LessonID,
		ItemTag:         requestData.ItemTag,
		Identifier:      requestData.Identifier,
		Duration:        requestData.Duration,
		StartTimeHour:   requestData.StartTimeHour,
		StartTimeMinute: requestData.StartTimeMinute,
		EndTimeHour:     requestData.EndTimeHour,
		EndTimeMinutes:  requestData.EndTimeMinutes,
		RoomIndex:       requestData.RoomIndex,
	})

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))

}
