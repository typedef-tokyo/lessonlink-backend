package model

import (
	"time"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/domain/vo"
)

type SessionModel struct {
	sessionID string
	userID    vo.UserID
	roleKey   vo.RoleKey
	expiresAt time.Time
}

func NewSessionModel(
	sessionID string,
	userID vo.UserID,
	roleKey vo.RoleKey,
	expiresAt time.Time,
) *SessionModel {

	return &SessionModel{
		sessionID: sessionID,
		userID:    userID,
		roleKey:   roleKey,
		expiresAt: expiresAt,
	}
}

func (r SessionModel) SessionID() string {
	return r.sessionID
}

func (r SessionModel) UserID() vo.UserID {
	return r.userID
}

func (r SessionModel) RoleKey() vo.RoleKey {
	return r.roleKey
}

func (r SessionModel) ExpiresAt() time.Time {
	return r.expiresAt
}

func (r SessionModel) IsSessionExpired() bool {
	return r.expiresAt.Before(time.Now())
}

func (r *SessionModel) KeepAliveSession() {
	const LIMITED = 3600 * time.Second
	r.expiresAt = time.Now().Add(LIMITED)
}
