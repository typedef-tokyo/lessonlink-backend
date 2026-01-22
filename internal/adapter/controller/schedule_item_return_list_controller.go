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
	IScheduleItemReturnListController interface {
		Execute(c echo.Context) error
	}

	ScheduleItemReturnListController struct {
		inputPort usecase.IScheduleItemReturnListInputPort
		presenter presenter.IScheduleItemEditPresenter
		logger    ILogWriter
	}
)

func NewScheduleItemReturnListController(
	inputPort usecase.IScheduleItemReturnListInputPort,
	presenter presenter.IScheduleItemEditPresenter,
	logger ILogWriter,
) IScheduleItemReturnListController {
	return &ScheduleItemReturnListController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	ScheduleItemReturnListRequestData struct {
		HistoryIndex int    `json:"history_index"`
		LessonID     int    `json:"lesson_id"`
		Identifier   string `json:"identifier"`
		Duration     int    `json:"duration"`
	}
)

// @Summary スケジュール編集アイテムリスト移動
// @Description
// @Produce json
// @Param schedule_id path string true "ScheduleID"
// @Param request body ScheduleItemReturnListRequestData true "アイテムリスト移動リクエスト"
// @Success 200 {object} presenter.ScheduleItemEditResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /schedule/{schedule_id}/item-return-list [post]
func (h *ScheduleItemReturnListController) Execute(c echo.Context) error {

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

	var requestData ScheduleItemReturnListRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	result, err := h.inputPort.Execute(c.Request().Context(), role, userID, scheduleID, requestData.HistoryIndex, usecase.ScheduleItemReturnListInput{
		LessonID:   requestData.LessonID,
		Identifier: requestData.Identifier,
		Duration:   requestData.Duration,
	})

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
