package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IUserUpdateController interface {
		Execute(c echo.Context) error
	}

	UserUpdateController struct {
		inputPort usecase.IUserUpdateInputPort
		presenter presenter.IUpdateUserPresenter
		logger    ILogWriter
	}
)

func NewUserUpdateController(
	inputPort usecase.IUserUpdateInputPort,
	presenter presenter.IUpdateUserPresenter,
	logger ILogWriter,
) IUserUpdateController {
	return &UserUpdateController{
		inputPort: inputPort,
		presenter: presenter,
		logger:    logger,
	}
}

type (
	UserUpdateRequestData struct {
		ID              int    `json:"id"`
		DisplayName     string `json:"display_name"`
		UserName        string `json:"user_name"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
		RoleKey         string `json:"role_key"`
	}
)

// @Summary ユーザー更新
// @Description
// @Produce json
// @Param request body UserUpdateRequestData true "ユーザー更新リクエスト"
// @Success 201 {object} presenter.UserUpdateResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user [put]
func (h *UserUpdateController) Execute(c echo.Context) error {

	// セッション情報を取得
	userID, roleKey, err := session_util.GetSessionData(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
	}

	var requestData UserUpdateRequestData

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

	err = h.inputPort.Execute(c.Request().Context(), usecase.UserUpdateInput{
		UserID:      requestData.ID,
		RoleKey:     requestData.RoleKey,
		UserName:    requestData.UserName,
		Password:    requestData.Password,
		DisplayName: requestData.DisplayName,
	}, userID, roleKey)

	if err != nil {
		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, map[string]any{
			"msg": msg,
		})
	}

	return c.JSON(http.StatusCreated, h.presenter.Present())
}
