package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IScheduleCreateController interface {
		Execute(c echo.Context) error
	}

	ScheduleCreateController struct {
		inputPort usecase.IScheduleCreateInputPort
		presenter presenter.IScheduleCreatePresenter
		logger    ILogWriter
	}
)

func NewScheduleCreateController(
	inputPort usecase.IScheduleCreateInputPort,
	presenter presenter.IScheduleCreatePresenter,
	logger ILogWriter,
) IScheduleCreateController {
	return &ScheduleCreateController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	ScheduleCreateRequestData struct {
		StartTime int `json:"start_time"`
		EndTime   int `json:"end_time"`
	}
)

// @Summary スケジュール作成
// @Description
// @Produce json
// @Param campus path string true "校舎"
// @Param request body ScheduleCreateRequestData true "スケジュール作成リクエスト"
// @Success 200 {object} presenter.ScheduleCreateResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/create/{campus} [post]
func (h *ScheduleCreateController) Execute(c echo.Context) error {

	var err error

	// セッション情報を取得
	userID, roleKey, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
	}

	campus := c.Param("campus")

	var requestData ScheduleCreateRequestData
	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	result, err := h.inputPort.Execute(c.Request().Context(), roleKey, userID, campus, requestData.StartTime, requestData.EndTime)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
