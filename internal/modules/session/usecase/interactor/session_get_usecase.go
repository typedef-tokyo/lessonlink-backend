package interactor

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	ISessionGetInputPort interface {
		Execute(ctx context.Context, sessionID string) (*SessionGetOutput, error)
	}
)

type (
	SessionGetOutput struct {
		UserID  int
		RoleKey string
	}
)

type SessionGetInteractor struct {
	repositorySession repository.SessionRepository
}

func NewSessionGetInteractor(
	repositorySession repository.SessionRepository,
) ISessionGetInputPort {
	return &SessionGetInteractor{
		repositorySession: repositorySession,
	}
}

func (r *SessionGetInteractor) Execute(ctx context.Context, sessionID string) (*SessionGetOutput, error) {

	sessionModel, err := r.repositorySession.Find(ctx, sessionID)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if sessionModel == nil {
		return nil, log.WrapErrorWithStackTraceUnauthorized(log.Errorf("ログインしてください"))
	}

	if sessionModel.IsSessionExpired() {
		return nil, log.WrapErrorWithStackTraceUnauthorized(log.Errorf("ログインしてください"))
	}

	// セッションを延長
	sessionModel.KeepAliveSession()

	err = r.repositorySession.Update(ctx, nil, sessionModel)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &SessionGetOutput{
		UserID:  sessionModel.UserID().Value(),
		RoleKey: sessionModel.RoleKey().Value(),
	}, nil
}
