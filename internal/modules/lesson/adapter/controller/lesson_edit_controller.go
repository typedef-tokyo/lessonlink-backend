package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/usecase/interactor"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/logger"
)

type (
	ILessonEditController interface {
		Execute(c echo.Context) error
	}

	LessonEditController struct {
		inputPort interactor.ILessonEditInputPort
		presenter presenter.ILessonEditPresenter
		logger    logger.ILogWriter
	}
)

func NewLessonEditController(
	inputPort interactor.ILessonEditInputPort,
	presenter presenter.ILessonEditPresenter,
	logger logger.ILogWriter,
) ILessonEditController {
	return &LessonEditController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	LessonEditRequestData struct {
		LessonName string `json:"lesson_name"`
		Duration   int    `json:"duration"`
	}
)

// @Summary 講座編集
// @Description
// @Produce json
// @Param lessonid path string true "講座"
// @Param request body LessonEditRequestData true "講座編集リクエスト"
// @Success 200 {object} presenter.LessonEditResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lesson/{lessonid} [patch]
func (h *LessonEditController) Execute(c echo.Context) error {

	// セッション情報を取得
	userID, roleKey, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
	}

	inputLessonid := c.Param("lessonid")
	lessonid, err := strconv.Atoi(inputLessonid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{ // not foundな気もする
			"msg": "講座IDの形式が不正です",
		})
	}

	var requestData LessonEditRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	err = h.inputPort.Execute(c, roleKey, userID, interactor.LessonEditInputDTO{
		ID:         lessonid,
		LessonName: requestData.LessonName,
		Duration:   requestData.Duration,
	})

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present())

}
