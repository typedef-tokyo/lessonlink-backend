package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IScheduleDeleteController interface {
		Execute(c echo.Context) error
	}

	ScheduleDeleteController struct {
		inputPort usecase.IScheduleDeleteInputPort
		logger    ILogWriter
	}
)

func NewScheduleDeleteController(
	inputPort usecase.IScheduleDeleteInputPort,
	logger ILogWriter,
) IScheduleDeleteController {
	return &ScheduleDeleteController{
		inputPort: inputPort,
		logger:    logger,
	}
}

// @Summary スケジュール削除
// @Description
// @Produce json
// @Param schedule_id path string true "ScheduleID"
// @Success 204
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /schedule/{schedule_id} [delete]
func (h *ScheduleDeleteController) Execute(c echo.Context) error {

	var err error

	// セッション情報を取得
	userID, roleKey, err := session_util.GetSessionData(c)
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

	err = h.inputPort.Execute(c.Request().Context(), roleKey, scheduleID, userID)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.NoContent(http.StatusNoContent)

}
