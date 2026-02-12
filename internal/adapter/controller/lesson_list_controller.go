package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/lessonlist"
)

type (
	ILessonListController interface {
		Execute(c echo.Context) error
	}

	LessonListController struct {
		inputPort lessonlist.ILessonListInputPort
		presenter presenter.ILessonListPresenter
		logger    ILogWriter
	}
)

func NewLessonListController(
	inputPort lessonlist.ILessonListInputPort,
	presenter presenter.ILessonListPresenter,
	logger ILogWriter,
) ILessonListController {
	return &LessonListController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

// @Summary 講座一覧取得
// @Description
// @Produce json
// @Param campus path string true "校舎"
// @Success 200 {object} presenter.LessonListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lesson/{campus}/list [get]
func (h *LessonListController) Execute(c echo.Context) error {

	campus := c.Param("campus")
	if campus == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "校舎識別子が不正です",
		})
	}

	result, err := h.inputPort.Execute(c.Request().Context(), campus)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusOK, h.presenter.Present(result))

}
