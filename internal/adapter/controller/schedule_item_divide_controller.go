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
	IScheduleItemDivideController interface {
		Execute(c echo.Context) error
	}

	ScheduleItemDivideController struct {
		inputPort usecase.IScheduleItemDivideInputPort
		presenter presenter.IScheduleItemEditPresenter
		logger    ILogWriter
	}
)

func NewScheduleItemDivideController(
	inputPort usecase.IScheduleItemDivideInputPort,
	presenter presenter.IScheduleItemEditPresenter,
	logger ILogWriter,
) IScheduleItemDivideController {
	return &ScheduleItemDivideController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	ScheduleItemDivideRequestData struct {
		HistoryIndex  int    `json:"history_index"`
		LessonID      int    `json:"lesson_id"`
		Identifier    string `json:"identifier"`
		DivideMinutes int    `json:"divide_minutes"`
	}
)

// @Summary スケジュール編集アイテムリスト分割
// @Description
// @Produce json
// @Param schedule_id path int true "ScheduleID"
// @Param request body ScheduleItemDivideRequestData true "スケジュール保存リクエスト"
// @Success 200 {object} presenter.ScheduleItemEditResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/{schedule_id}/item-divide [post]
func (h *ScheduleItemDivideController) Execute(c echo.Context) error {

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

	var requestData ScheduleItemDivideRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	result, err := h.inputPort.Execute(c.Request().Context(), role, userID, scheduleID, requestData.HistoryIndex, usecase.ScheduleItemDivideInput{
		LessonID:      requestData.LessonID,
		Identifier:    requestData.Identifier,
		DivideMinutes: requestData.DivideMinutes,
	})

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
