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
	IScheduleSaveTitleController interface {
		Execute(c echo.Context) error
	}

	ScheduleSaveTitleController struct {
		inputPort usecase.IScheduleSaveTitleInputPort
		presenter presenter.IScheduleSaveTitlePresenter
		logger    ILogWriter
	}
)

func NewScheduleSaveTitleController(
	inputPort usecase.IScheduleSaveTitleInputPort,
	presenter presenter.IScheduleSaveTitlePresenter,
	logger ILogWriter,
) IScheduleSaveTitleController {
	return &ScheduleSaveTitleController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	ScheduleSaveTitleRequestData struct {
		Title string `json:"title"`
	}
)

// @Summary スケジュールタイトル保存
// @Description
// @Produce json
// @Param schedule_id path int true "ScheduleID"
// @Param request body ScheduleSaveTitleRequestData true "タイトル保存リクエスト"
// @Success 200 {object} presenter.ScheduleSaveTitleResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/{schedule_id}/title [patch]
func (h *ScheduleSaveTitleController) Execute(c echo.Context) error {

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

	var requestData ScheduleSaveTitleRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	err = h.inputPort.Execute(c.Request().Context(), role, userID, scheduleID, requestData.Title)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present())
}
