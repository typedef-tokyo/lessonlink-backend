package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IScheduleTimeEditController interface {
		Execute(c echo.Context) error
	}

	ScheduleTimeEditController struct {
		inputPort usecase.IScheduleTimeEditInputPort
		logger    ILogWriter
	}
)

func NewScheduleTimeEditController(
	inputPort usecase.IScheduleTimeEditInputPort,
	logger ILogWriter,
) IScheduleTimeEditController {
	return &ScheduleTimeEditController{
		inputPort: inputPort,
		logger:    logger,
	}
}

type (
	ScheduleTimeEditRequestData struct {
		StartTime int `json:"start_time"`
		EndTime   int `json:"end_time"`
	}
)

// @Summary スケジュール時間変更
// @Description
// @Produce json
// @Param schedule_id path string true "ScheduleID"
// @Param request body ScheduleTimeEditRequestData true "スケジュール時間変更リクエスト"
// @Success 200 {object} presenter.ScheduleItemEditResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /schedule/{schedule_id}/time [patch]
func (h *ScheduleTimeEditController) Execute(c echo.Context) error {

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

	var requestData ScheduleTimeEditRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	err = h.inputPort.Execute(c.Request().Context(), role, userID, scheduleID, requestData.StartTime, requestData.EndTime)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.NoContent(http.StatusOK)
}
