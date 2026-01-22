package usecase

import (
	"context"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/entity"
	session "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type IUserLoginInputPort interface {
	Execute(ctx context.Context, input UserLoginInput) (*UserLoginOutput, error)
}

type UserLoginInput struct {
	UserName        string
	UserRawPassword string
}

type (
	UserLoginOutput struct {
		User      *UserLoginOutputDTO
		Err       error
		SessionID string
	}

	UserLoginOutputDTO struct {
		ID       int
		Name     string
		UserName string
		RoleKey  string
	}
)

type UserLoginInteractor struct {
	txManager         util.TxManager
	repositoryUser    repository.UserRepository
	repositorySession session.SessionRepository
}

func NewUserLoginInteractor(
	txManager util.TxManager,
	repositoryUser repository.UserRepository,
	repositorySession session.SessionRepository,
) IUserLoginInputPort {
	return &UserLoginInteractor{
		txManager:         txManager,
		repositoryUser:    repositoryUser,
		repositorySession: repositorySession,
	}
}

func (r *UserLoginInteractor) Execute(ctx context.Context, input UserLoginInput) (*UserLoginOutput, error) {

	user, err := r.repositoryUser.FindByUserName(ctx, input.UserName)

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	password, err := vo.ReconstructHashedPassword(input.UserRawPassword)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	if user == nil || !user.AuthenticatePassword(password) {
		return nil, log.WrapErrorWithStackTraceUnauthorized(log.Errorf("ユーザー名またはパスワードが違います"))
	}

	sessionID := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(securecookie.GenerateRandomKey(32))
	entity := entity.SessionEntity{
		SessionID: sessionID,
		UserID:    user.ID(),
		RoleKey:   user.RoleKey(),
		ExpiresAt: time.Now().Add(60 * time.Minute),
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		if err = r.repositorySession.Save(ctx, tx, entity); err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &UserLoginOutput{
		User: &UserLoginOutputDTO{
			ID:       user.ID().Value(),
			Name:     user.DisplayName().Value(),
			UserName: user.UserName().Value(),
			RoleKey:  user.RoleKey().Value(),
		},
		SessionID: sessionID,
	}, nil
}
