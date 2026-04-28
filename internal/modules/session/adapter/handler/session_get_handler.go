package handler

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/usecase/interactor"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/logger"
)

type (
	ISessionGetHandler interface {
		Execute(ctx context.Context, sessionID string) (userID int, roleKey string, err error)
	}

	SessionGetHandler struct {
		interactor interactor.ISessionGetInputPort
		logger     logger.ILogWriter
	}
)

func NewSessionGetHandler(
	interactor interactor.ISessionGetInputPort,
	logger logger.ILogWriter,
) ISessionGetHandler {
	return &SessionGetHandler{
		interactor: interactor,
		logger:     logger,
	}
}

func (h *SessionGetHandler) Execute(ctx context.Context, sessionID string) (userID int, roleKey string, err error) {

	result, err := h.interactor.Execute(ctx, sessionID)

	if err != nil {
		return vo.USER_ID_INVALID.Value(), vo.ROLE_KEY_INVALID.Value(), err
	}

	return result.UserID, result.RoleKey, err
}
