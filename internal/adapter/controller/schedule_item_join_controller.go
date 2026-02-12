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
	IScheduleItemJoinController interface {
		Execute(c echo.Context) error
	}

	ScheduleItemJoinController struct {
		inputPort usecase.IScheduleItemJoinInputPort
		presenter presenter.IScheduleItemEditPresenter
		logger    ILogWriter
	}
)

func NewScheduleItemJoinController(
	inputPort usecase.IScheduleItemJoinInputPort,
	presenter presenter.IScheduleItemEditPresenter,
	logger ILogWriter,
) IScheduleItemJoinController {
	return &ScheduleItemJoinController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	ScheduleItemJoinRequestData struct {
		HistoryIndex       int    `json:"history_index"`
		JoinFromIdentifier string `json:"join_from_identifier"`
		JoinToIdentifier   string `json:"join_to_identifier"`
	}
)

// @Summary スケジュール編集アイテムリスト分割
// @Description
// @Produce json
// @Param schedule_id path int true "ScheduleID"
// @Param request body ScheduleItemJoinRequestData true "アイテム結合リクエスト"
// @Success 200 {object} presenter.ScheduleItemEditResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/{schedule_id}/item-join [post]
func (h *ScheduleItemJoinController) Execute(c echo.Context) error {

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

	var requestData ScheduleItemJoinRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	result, err := h.inputPort.Execute(c.Request().Context(), role, userID, scheduleID, requestData.HistoryIndex, requestData.JoinFromIdentifier, requestData.JoinToIdentifier)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
