package controller

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IUserLoginController interface {
		Execute(c echo.Context) error
	}

	UserLoginController struct {
		inputPort usecase.IUserLoginInputPort
		presenter presenter.IUserLoginPresenter
		env       configs.EnvConfig
		logger    ILogWriter
	}
)

func NewUserLoginController(
	inputPort usecase.IUserLoginInputPort,
	presenter presenter.IUserLoginPresenter,
	env configs.EnvConfig,
	logger ILogWriter,
) IUserLoginController {
	return &UserLoginController{
		inputPort: inputPort,
		presenter: presenter,
		env:       env,
		logger:    logger,
	}
}

type UserLoginParams struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// @Summary ユーザーログイン
// @Description
// @Produce json
// @Param request body controller.UserLoginParams true "ユーザーログイン情報"
// @Success 200 {object} presenter.UserLoginResponse
// @Failure 400 {object} presenter.UserLoginResponse
// @Failure 401 {object} presenter.UserLoginResponse
// @Failure 500 {object} presenter.UserLoginResponse
// @Router /user/login [post]
func (h *UserLoginController) Execute(c echo.Context) error {

	var bodyParams UserLoginParams

	err := c.Bind(&bodyParams)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "リクエスト形式が不正です",
		})
	}

	// ユーザー名を取得
	userName := utility.Trimmer(bodyParams.UserName)

	// パスワードを取得
	password := utility.Trimmer(bodyParams.Password)

	// ユースケースを実行
	result, err := h.inputPort.Execute(c.Request().Context(), usecase.UserLoginInput{UserName: userName, UserRawPassword: password})

	if err != nil {

		type LoginUserDTO struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			UserName string `json:"user_name"`
			RoleKey  string `json:"role_key"`
		}

		type UserLoginResponse struct {
			Msg       string       `json:"msg"`
			LoginUser LoginUserDTO `json:"login_user"`
		}

		status, msg := h.logger.WriteErrLog(c, err)
		return c.JSON(status, UserLoginResponse{Msg: msg,

			LoginUser: LoginUserDTO{
				Id:       -1,
				Name:     "",
				UserName: "",
				RoleKey:  "",
			},
		})
	}

	sameSite := http.SameSiteStrictMode
	if os.Getenv("LOCAL_TEST") == "true" {
		sameSite = http.SameSiteNoneMode
	}

	cookie := &http.Cookie{
		Name:     h.env.SessionName,
		Value:    result.SessionID,
		Path:     "/",
		MaxAge:   3600,
		Domain:   "",
		Secure:   os.Getenv("LOCAL_TEST") != "true",
		HttpOnly: true,
		SameSite: sameSite,
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, h.presenter.Present(result))
}
