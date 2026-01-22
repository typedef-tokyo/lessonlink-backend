package rdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/entity"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/repository"
)

type value struct {
	UserID    int       `json:"user_id"`
	RoleKey   string    `json:"role_key"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Session struct {
	c *sql.DB
}

func NewSessionRepository(c IMySQL) repository.SessionRepository {
	return &Session{c: c.GetConn()}
}

func (f *Session) Save(ctx context.Context, tx *sql.Tx, session entity.SessionEntity) error {

	sessionDTO, err := f.toDTO(session)
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	err = sessionDTO.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return nil
}

func (f *Session) Update(ctx context.Context, tx *sql.Tx, session entity.SessionEntity) error {

	sessionDTO, err := f.toDTO(session)
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	if tx != nil {
		_, err = sessionDTO.Update(ctx, tx, boil.Whitelist(
			dto.SysSessionColumns.Value,
			dto.SysSessionColumns.UpdatedAt,
		))
	} else {
		_, err = sessionDTO.Update(ctx, f.c, boil.Infer())
	}

	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return nil
}

func (f *Session) Delete(ctx context.Context, tx *sql.Tx, userID vo.UserID) error {

	sessionDTOs, err := dto.SysSessions(
		dto.SysSessionWhere.UserID.EQ(userID.Value()),
	).All(ctx, f.c)

	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	var errs error
	for _, session := range sessionDTOs {

		if tx != nil {
			_, err = session.Delete(ctx, tx)
		} else {
			_, err = session.Delete(ctx, f.c)
		}

		if err != nil {
			errs = errors.Join(errs, log.WrapErrorWithStackTraceInternalServerError(err))
		}
	}

	if errs != nil {
		return log.WrapErrorWithStackTrace(log.Errorf("%s", errs.Error()))
	}

	return nil
}

func (f *Session) Find(ctx context.Context, sessionID string) (*entity.SessionEntity, error) {

	sessionDTO, err := dto.SysSessions(
		dto.SysSessionWhere.SessionID.EQ(sessionID),
	).All(ctx, f.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	var entity *entity.SessionEntity
	for _, dto := range sessionDTO {

		entity, err = f.toEntity(dto)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

	}

	return entity, nil
}

func (f *Session) toEntity(session *dto.SysSession) (*entity.SessionEntity, error) {

	var v value
	err := json.Unmarshal([]byte(session.Value), &v)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	_userID, err := vo.NewUserID(v.UserID)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	_roleKey, err := vo.NewRoleKey(v.RoleKey)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	entity := &entity.SessionEntity{
		SessionID: session.SessionID,
		UserID:    _userID,
		RoleKey:   _roleKey,
		ExpiresAt: v.ExpiresAt,
	}

	return entity, nil
}

func (f *Session) toDTO(session entity.SessionEntity) (*dto.SysSession, error) {

	sessionJson, err := json.Marshal(value{
		UserID:    session.UserID.Value(),
		RoleKey:   session.RoleKey.Value(),
		ExpiresAt: session.ExpiresAt,
	})

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return &dto.SysSession{
		SessionID: session.SessionID,
		UserID:    session.UserID.Value(),
		Value:     string(sessionJson),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}, nil
}
