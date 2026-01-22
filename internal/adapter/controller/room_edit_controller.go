package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IRoomEditController interface {
		Execute(c echo.Context) error
	}

	RoomEditController struct {
		inputPort usecase.IRoomEditInputPort
		presenter presenter.IRoomEditPresenter
		logger    ILogWriter
	}
)

func NewRoomEditController(
	inputPort usecase.IRoomEditInputPort,
	presenter presenter.IRoomEditPresenter,
	logger ILogWriter,
) IRoomEditController {
	return &RoomEditController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	RoomEditRequestData struct {
		RoomList []RoomEditData `json:"room_list"`
	}

	RoomEditData struct {
		RoomIndex int    `json:"room_index"`
		Name      string `json:"room_name"`
	}
)

// @Summary 教室編集
// @Description
// @Produce json
// @Param campus path string true "校舎"
// @Param request body RoomEditRequestData true "教室編集リクエスト"
// @Success 200 {object} presenter.RoomEditResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /room/{campus}/edit [post]
func (h *RoomEditController) Execute(c echo.Context) error {

	// セッション情報を取得
	_, roleKey, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
	}

	campus := c.Param("campus")
	if campus == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "校舎識別子が不正です",
		})
	}

	var requestData RoomEditRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	err = h.inputPort.Execute(c.Request().Context(), roleKey, campus, usecase.RoomsEditInputDTO{
		Rooms: lo.Map(requestData.RoomList, func(item RoomEditData, _ int) usecase.RoomEditInputDTO {
			return usecase.RoomEditInputDTO{
				Index: item.RoomIndex,
				Name:  item.Name,
			}
		}),
	})

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present())

}
