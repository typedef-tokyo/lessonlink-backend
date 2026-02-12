package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IScheduleDuplicateController interface {
		Execute(c echo.Context) error
	}

	ScheduleDuplicateController struct {
		inputPort usecase.IScheduleDuplicatePort
		logger    ILogWriter
	}
)

func NewScheduleDuplicateController(
	inputPort usecase.IScheduleDuplicatePort,
	logger ILogWriter,
) IScheduleDuplicateController {
	return &ScheduleDuplicateController{
		inputPort: inputPort,
		logger:    logger,
	}
}

// @Summary スケジュール複製
// @Description
// @Produce json
// @Param schedule_id path int true "ScheduleID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/{schedule_id}/duplicate [post]
func (h *ScheduleDuplicateController) Execute(c echo.Context) error {

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

	err = h.inputPort.Execute(c.Request().Context(), role, userID, scheduleID)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.NoContent(http.StatusNoContent)
}
