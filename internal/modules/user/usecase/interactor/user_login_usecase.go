package interactor

import (
	"context"
	"encoding/base32"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/usecase/public"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
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
	txManager         database.TxManager
	repositoryUser    repository.UserRepository
	FacadeSessionSave external.ISessionSaveFacade
}

func NewUserLoginInteractor(
	txManager database.TxManager,
	repositoryUser repository.UserRepository,
	FacadeSessionSave external.ISessionSaveFacade,
) IUserLoginInputPort {
	return &UserLoginInteractor{
		txManager:         txManager,
		repositoryUser:    repositoryUser,
		FacadeSessionSave: FacadeSessionSave,
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
	sessionDTO := public.SessionSaveInputDTO{
		SessionID: sessionID,
		UserID:    user.ID().Value(),
		RoleKey:   user.RoleKey().Value(),
		ExpiresAt: time.Now().Add(60 * time.Minute),
	}

	if err = r.FacadeSessionSave.Execute(ctx, sessionDTO); err != nil {
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
