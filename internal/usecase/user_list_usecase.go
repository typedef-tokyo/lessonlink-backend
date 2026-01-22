package usecase

import (
	"context"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/user"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type IUserListInputPort interface {
	Execute(ctx context.Context, userID vo.UserID, roleKey vo.RoleKey) (*UserListOutput, error)
}

type UserListInteractor struct {
	repositoryRole repository.RoleRepository
	repositoryUser repository.UserRepository
}

type (
	UserListOutput struct {
		UserList []*UserListOutputDTO
	}

	UserListOutputDTO struct {
		ID          int
		RoleName    string
		UserName    string
		DisplayName string
	}
)

func NewUserListInteractor(
	repositoryRole repository.RoleRepository,
	repositoryUser repository.UserRepository,
) IUserListInputPort {
	return &UserListInteractor{
		repositoryRole: repositoryRole,
		repositoryUser: repositoryUser,
	}
}

func (r UserListInteractor) Execute(ctx context.Context, userID vo.UserID, roleKey vo.RoleKey) (*UserListOutput, error) {

	var err error
	roles, err := r.repositoryRole.FindAll(ctx)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	users, err := r.repositoryUser.FindAll(ctx)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	// オーナー以外は自分の情報のみを取得する
	if !roleKey.IsOwner() {

		mySelf := users.FindByUserID(userID)
		if mySelf == nil {
			return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("ユーザーが見つかりません"))
		}

		users = user.RootUserModelSlice{mySelf}
	}

	return &UserListOutput{
		UserList: lo.Map(users, func(item *user.RootUserModel, _ int) *UserListOutputDTO {
			return &UserListOutputDTO{
				ID:          item.ID().Value(),
				RoleName:    roles.FindNameByKey(item.RoleKey()).Value(),
				UserName:    item.UserName().Value(),
				DisplayName: item.DisplayName().Value(),
			}
		}),
	}, nil
}
