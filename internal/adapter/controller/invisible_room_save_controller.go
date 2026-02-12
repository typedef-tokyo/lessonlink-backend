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
	IInvisibleRoomController interface {
		Execute(c echo.Context) error
	}

	InvisibleRoomController struct {
		inputPort usecase.IInvisibleRoomSaveInputPort
		presenter presenter.IInvisibleRoomPresenter
		logger    ILogWriter
	}
)

func NewInvisibleRoomController(
	inputPort usecase.IInvisibleRoomSaveInputPort,
	presenter presenter.IInvisibleRoomPresenter,
	logger ILogWriter,
) IInvisibleRoomController {
	return &InvisibleRoomController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	InvisibleRoomSaveRequestData struct {
		InvisibleRooms []int `json:"invisible_rooms"`
	}
)

// @Summary 非表示ルーム登録
// @Description
// @Produce json
// @Param schedule_id path int true "ScheduleID"
// @Param request body InvisibleRoomSaveRequestData true "非表示ルームリクエスト"
// @Success 200 {object} presenter.InvisibleRoomSaveResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/{schedule_id}/room/invisible [put]
func (h *InvisibleRoomController) Execute(c echo.Context) error {

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

	var requestData InvisibleRoomSaveRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	err = h.inputPort.Execute(c.Request().Context(), roleKey, userID, scheduleID, requestData.InvisibleRooms)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present())
}
