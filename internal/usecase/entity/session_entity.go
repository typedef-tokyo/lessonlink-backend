package entity

import (
	"time"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type SessionEntity struct {
	SessionID string
	UserID    vo.UserID
	RoleKey   vo.RoleKey
	ExpiresAt time.Time
}
