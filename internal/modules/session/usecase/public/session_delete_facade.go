package public

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	SessionDeleteFacade struct {
		repositorySession repository.SessionRepository
	}
)

func NewSessionDeleteFacade(
	repositorySession repository.SessionRepository,
) *SessionDeleteFacade {
	return &SessionDeleteFacade{
		repositorySession: repositorySession,
	}
}

func (r SessionDeleteFacade) Execute(ctx context.Context, inputDeleteUserID int) error {

	deleteUserID, err := vo.NewUserID(inputDeleteUserID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	err = r.repositorySession.Delete(ctx, nil, deleteUserID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil
}
