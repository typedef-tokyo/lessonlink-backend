package repository

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/entity"
)

type SessionRepository interface {
	Save(ctx context.Context, tx *sql.Tx, session entity.SessionEntity) error
	Update(ctx context.Context, tx *sql.Tx, session entity.SessionEntity) error
	Delete(ctx context.Context, tx *sql.Tx, userID vo.UserID) error
	Find(ctx context.Context, sessionID string) (*entity.SessionEntity, error)
}
