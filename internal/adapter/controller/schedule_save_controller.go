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
	IScheduleSaveController interface {
		Execute(c echo.Context) error
	}

	ScheduleSaveController struct {
		inputPort usecase.IScheduleSaveInputPort
		presenter presenter.IScheduleSavePresenter
		logger    ILogWriter
	}
)

func NewScheduleSaveController(
	inputPort usecase.IScheduleSaveInputPort,
	presenter presenter.IScheduleSavePresenter,
	logger ILogWriter,
) IScheduleSaveController {
	return &ScheduleSaveController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	ScheduleSaveRequestData struct {
		HistoryIndex int `json:"history_index"`
	}
)

// @Summary スケジュール保存
// @Description
// @Produce json
// @Param schedule_id path string true "ScheduleID"
// @Param request body ScheduleSaveRequestData true "スケジュール保存リクエスト"
// @Success 200 {object} presenter.ScheduleSaveResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /schedule/{schedule_id} [post]
func (h *ScheduleSaveController) Execute(c echo.Context) error {

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

	var requestData ScheduleSaveRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	result, err := h.inputPort.Execute(c.Request().Context(), role, userID, scheduleID, requestData.HistoryIndex)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
