package presenter

import (
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type IUserLoginPresenter interface {
	Present(result *usecase.UserLoginOutput) *UserLoginResponse
}

type UserLoginPresenter struct {
}

func NewUserLoginPresenter() IUserLoginPresenter {
	return &UserLoginPresenter{}
}

type (
	LoginUserDTO struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		UserName string `json:"user_name"`
		RoleKey  string `json:"role_key"`
	}

	UserLoginResponse struct {
		Msg       string       `json:"msg"`
		LoginUser LoginUserDTO `json:"login_user"`
	}
)

func (h *UserLoginPresenter) Present(result *usecase.UserLoginOutput) *UserLoginResponse {

	return &UserLoginResponse{
		Msg: "ログインが完了しました",
		LoginUser: LoginUserDTO{
			Id:       result.User.ID,
			Name:     result.User.Name,
			UserName: result.User.UserName,
			RoleKey:  result.User.RoleKey,
		},
	}
}
