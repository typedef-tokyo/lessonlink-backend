package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	ILessonAddController interface {
		Execute(c echo.Context) error
	}

	LessonAddController struct {
		inputPort usecase.ILessonAddInputPort
		presenter presenter.ILessonAddPresenter
		logger    ILogWriter
	}
)

func NewLessonAddController(
	inputPort usecase.ILessonAddInputPort,
	presenter presenter.ILessonAddPresenter,
	logger ILogWriter,
) ILessonAddController {
	return &LessonAddController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	LessonAddRequestData struct {
		LessonName string `json:"lesson_name"`
		Duration   int    `json:"duration"`
	}
)

// @Summary 講座追加
// @Description
// @Produce json
// @Param campus path string true "校舎"
// @Param request body LessonAddRequestData true "講座追加リクエスト"
// @Success 200 {object} presenter.LessonAddResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lesson/{campus} [post]
func (h *LessonAddController) Execute(c echo.Context) error {

	_, role, err := session_util.GetSessionData(c)
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

	var requestData LessonAddRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	err = h.inputPort.Execute(c.Request().Context(), role, campus, usecase.LessonAddInputDTO{
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
