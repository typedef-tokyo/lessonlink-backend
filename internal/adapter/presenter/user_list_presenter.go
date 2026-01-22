package presenter

import (
	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type IUserListPresenter interface {
	Present(result *usecase.UserListOutput) *UserListResponse
}

type UserListPresenter struct {
}

func NewUserListPresenter() IUserListPresenter {
	return &UserListPresenter{}
}

type (
	UserListResponse struct {
		Users []UserDTO `json:"users"`
	}

	UserDTO struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		UserName string `json:"user_name"`
		RoleName string `json:"role_name"`
	}
)

func (h *UserListPresenter) Present(result *usecase.UserListOutput) *UserListResponse {

	return &UserListResponse{
		Users: lo.Map(result.UserList, func(item *usecase.UserListOutputDTO, _ int) UserDTO {
			return UserDTO{
				ID:       item.ID,
				Name:     item.DisplayName,
				UserName: item.UserName,
				RoleName: item.RoleName,
			}
		}),
	}
}
