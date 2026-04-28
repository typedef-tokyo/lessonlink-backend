package interactor

import (
	"context"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/usecase/public"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/domain/model/user"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type IUserListInputPort interface {
	Execute(ctx context.Context, inputRole string, inputUserID int) (*UserListOutput, error)
}

type UserListInteractor struct {
	facadeRole     external.IRoleGetFacade
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
	facadeRole external.IRoleGetFacade,
	repositoryUser repository.UserRepository,
) IUserListInputPort {
	return &UserListInteractor{
		facadeRole:     facadeRole,
		repositoryUser: repositoryUser,
	}
}

func (r UserListInteractor) Execute(ctx context.Context, inputRole string, inputUserID int) (*UserListOutput, error) {

	var err error

	role, err := vo.NewRoleKey(inputRole)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	userID, err := vo.NewUserID(inputUserID)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	roles, err := r.facadeRole.Execute(ctx)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	users, err := r.repositoryUser.FindAll(ctx)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	// オーナー以外は自分の情報のみを取得する
	if !role.IsOwner() {

		mySelf := users.FindByUserID(userID)
		if mySelf == nil {
			return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("ユーザーが見つかりません"))
		}

		users = user.RootUserModelSlice{mySelf}
	}

	return &UserListOutput{

		UserList: lo.Map(users, func(item *user.RootUserModel, _ int) *UserListOutputDTO {

			roleItem, found := lo.Find(roles.Roles, func(role public.RoleGetOutDTO) bool {
				return role.RoleKey == item.RoleKey().Value()
			})

			roleName := roleItem.RoleKey
			if !found {
				roleName = "role_name_unknown"
			}

			return &UserListOutputDTO{
				ID:          item.ID().Value(),
				RoleName:    roleName,
				UserName:    item.UserName().Value(),
				DisplayName: item.DisplayName().Value(),
			}
		}),
	}, nil
}
