package presenter

import (
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	IUserGetPresenter interface {
		Present(result *usecase.UserGetOutput) *UserGetResponse
	}

	UserGetPresenter struct {
	}
)

func NewUserGetPresenter() IUserGetPresenter {
	return &UserGetPresenter{}
}

type (
	UserGetResponse struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		UserName string `json:"user_name"`
		RoleKey  string `json:"role_key"`
	}
)

func (h *UserGetPresenter) Present(result *usecase.UserGetOutput) *UserGetResponse {

	return &UserGetResponse{
		ID:       result.ID,
		Name:     result.Name,
		UserName: result.UserName,
		RoleKey:  result.RoleKey,
	}
}
