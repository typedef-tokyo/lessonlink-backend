package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/usecase/interactor"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/logger"
)

type (
	IUserAddController interface {
		Execute(c echo.Context) error
	}

	UserAddController struct {
		inputPort interactor.IUserAddPort
		presenter presenter.IUserAddPresenter
		logger    logger.ILogWriter
	}
)

func NewUserAddController(
	inputPort interactor.IUserAddPort,
	presenter presenter.IUserAddPresenter,
	logger logger.ILogWriter,
) IUserAddController {
	return &UserAddController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	UserAddRequestData struct {
		Name            string `json:"name"`
		UserName        string `json:"user_name"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
		RoleKey         string `json:"role_key"`
	}
)

// @Summary ユーザー追加
// @Description
// @Produce json
// @Param request body UserAddRequestData true "ユーザー作成リクエスト"
// @Success 201 {object} presenter.UserAddResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user [post]
func (h *UserAddController) Execute(c echo.Context) error {

	// セッション情報を取得
	userID, roleKey, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
	}

	var requestData UserAddRequestData

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	// パスワードと確認用パスワードが一致しているか確認
	if requestData.Password != requestData.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "パスワードと確認用パスワードが異なります",
		})
	}

	err = h.inputPort.Execute(c.Request().Context(), interactor.UserAddInput{
		RoleKey:     requestData.RoleKey,
		UserName:    requestData.UserName,
		Password:    requestData.Password,
		DisplayName: requestData.Name,
	}, roleKey, userID)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusCreated, h.presenter.Present())

}
