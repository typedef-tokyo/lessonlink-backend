package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/usecase/interactor"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/logger"
)

type (
	IRoomEditController interface {
		Execute(c echo.Context) error
	}

	RoomEditController struct {
		inputPort interactor.IRoomEditInputPort
		presenter presenter.IRoomEditPresenter
		logger    logger.ILogWriter
	}
)

func NewRoomEditController(
	inputPort interactor.IRoomEditInputPort,
	presenter presenter.IRoomEditPresenter,
	logger logger.ILogWriter,
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
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
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

	err = h.inputPort.Execute(c.Request().Context(), roleKey, campus, interactor.RoomsEditInputDTO{
		Rooms: lo.Map(requestData.RoomList, func(item RoomEditData, _ int) interactor.RoomEditInputDTO {
			return interactor.RoomEditInputDTO{
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
