package public

import (
	"context"
	"time"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/model"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type SessionSaveInputDTO struct {
	SessionID string
	UserID    int
	RoleKey   string
	ExpiresAt time.Time
}

type (
	SessionSaveFacade struct {
		repositorySession repository.SessionRepository
	}
)

func NewSessionSaveFacade(
	repositorySession repository.SessionRepository,
) *SessionSaveFacade {
	return &SessionSaveFacade{
		repositorySession: repositorySession,
	}
}

func (r SessionSaveFacade) Execute(ctx context.Context, inputDTO SessionSaveInputDTO) error {

	userID, err := vo.NewUserID(inputDTO.UserID)
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	roleKey, err := vo.NewRoleKey(inputDTO.RoleKey)
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	model := model.NewSessionModel(
		inputDTO.SessionID,
		userID,
		roleKey,
		inputDTO.ExpiresAt,
	)

	if err := r.repositorySession.Save(ctx, nil, model); err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil
}
