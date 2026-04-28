package public

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	UserGetFacade struct {
		repositoryUser repository.UserRepository
	}
)

func NewUserGetFacade(
	repositoryUser repository.UserRepository,
) *UserGetFacade {
	return &UserGetFacade{
		repositoryUser: repositoryUser,
	}
}

func (r UserGetFacade) Execute(ctx context.Context, inputUserID int) (userName string, role string, err error) {

	userID, err := vo.NewUserID(inputUserID)
	if err != nil {
		return "", "", log.WrapErrorWithStackTraceBadRequest(err)
	}

	user, err := r.repositoryUser.FindByUserID(ctx, userID)
	if err != nil {
		return "", "", log.WrapErrorWithStackTrace(err)
	}

	return user.UserName().Value(), user.RoleKey().Value(), nil
}
