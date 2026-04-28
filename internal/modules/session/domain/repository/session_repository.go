package repository

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/model"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/vo"
)

type SessionRepository interface {
	Save(ctx context.Context, tx *sql.Tx, session *model.SessionModel) error
	Update(ctx context.Context, tx *sql.Tx, session *model.SessionModel) error
	Delete(ctx context.Context, tx *sql.Tx, userID vo.UserID) error
	Find(ctx context.Context, sessionID string) (*model.SessionModel, error)
}
