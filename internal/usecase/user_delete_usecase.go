package usecase

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	session "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type IUserDeleteInputPort interface {
	Execute(ctx context.Context, role vo.RoleKey, userID int, deleteFromUserID vo.UserID) error
}

type (
	UserDeleteInteractor struct {
		txManager         util.TxManager
		repositoryUser    repository.UserRepository
		repositorySession session.SessionRepository
		env               configs.EnvConfig
	}
)

func NewUserDeleteInteractor(
	txManager util.TxManager,
	repositoryUser repository.UserRepository,
	repositorySession session.SessionRepository,
	env configs.EnvConfig,
) IUserDeleteInputPort {
	return &UserDeleteInteractor{
		txManager:         txManager,
		repositoryUser:    repositoryUser,
		repositorySession: repositorySession,
		env:               env,
	}
}

func (r UserDeleteInteractor) Execute(ctx context.Context, role vo.RoleKey, userID int, deleteFromUserID vo.UserID) error {

	deleteUserID, err := vo.NewUserID(userID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if role != vo.ROLE_KEY_OWNER {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	deleteUserData, err := r.repositoryUser.FindByUserID(ctx, deleteUserID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if deleteUserData == nil {
		return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したユーザーは存在しません:%d", deleteUserID))
	}

	if deleteUserData.IsEnableDelete(deleteFromUserID) {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	var sessionDeleteErr error
	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		if err = r.repositoryUser.Delete(ctx, tx, deleteUserData, deleteFromUserID); err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		sessionDeleteErr = r.repositorySession.Delete(ctx, nil, deleteUserID)
		if sessionDeleteErr != nil {
			sessionDeleteErr = log.WrapErrorWithStackTrace(sessionDeleteErr)
		}

		return nil
	})

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return sessionDeleteErr
}
