package usecase

import (
	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	IUserGetInputPort interface {
		Execute(c echo.Context, userID int) (*UserGetOutput, error)
	}
)

type (
	UserGetOutput struct {
		ID       int
		RoleKey  string
		UserName string
		Name     string
	}
)

type UserGetInteractor struct {
	repositoryRole repository.RoleRepository
	repositoryUser repository.UserRepository
}

func NewUserGetInteractor(
	repositoryRole repository.RoleRepository,
	repositoryUser repository.UserRepository,
) IUserGetInputPort {
	return &UserGetInteractor{
		repositoryRole: repositoryRole,
		repositoryUser: repositoryUser,
	}
}

func (r UserGetInteractor) Execute(c echo.Context, _userID int) (*UserGetOutput, error) {

	userID, err := vo.NewUserID(_userID)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	ctx := c.Request().Context()
	user, err := r.repositoryUser.FindByUserID(ctx, userID)

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if user == nil {
		return nil, log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したユーザーは存在しません:%d", _userID))
	}

	return &UserGetOutput{
		ID:       user.ID().Value(),
		RoleKey:  user.RoleKey().Value(),
		UserName: user.UserName().Value(),
		Name:     user.DisplayName().Value(),
	}, nil
}
