package interactor

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
)

type IUserDeleteInputPort interface {
	Execute(ctx context.Context, inputRole string, inputUserID int, inputDeleteFromUserID int) error
}

type (
	UserDeleteInteractor struct {
		txManager           database.TxManager
		repositoryUser      repository.UserRepository
		facadeSessionDelete external.ISessionDeleteFacade
		env                 configs.EnvConfig
	}
)

func NewUserDeleteInteractor(
	txManager database.TxManager,
	repositoryUser repository.UserRepository,
	facadeSessionDelete external.ISessionDeleteFacade,
	env configs.EnvConfig,
) IUserDeleteInputPort {
	return &UserDeleteInteractor{
		txManager:           txManager,
		repositoryUser:      repositoryUser,
		facadeSessionDelete: facadeSessionDelete,
		env:                 env,
	}
}

func (r UserDeleteInteractor) Execute(ctx context.Context, inputRole string, inputUserID int, inputDeleteFromUserID int) error {

	deleteUserID, err := vo.NewUserID(inputUserID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	role, err := vo.NewRoleKey(inputRole)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	if !role.IsOwner() {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	deleteUserData, err := r.repositoryUser.FindByUserID(ctx, deleteUserID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if deleteUserData == nil {
		return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したユーザーは存在しません:%d", deleteUserID))
	}

	deleteFromUserID, err := vo.NewUserID(inputDeleteFromUserID)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	if deleteUserData.IsEnableDelete(deleteFromUserID) {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		if err = r.repositoryUser.Delete(ctx, tx, deleteUserData, deleteFromUserID); err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	err = r.facadeSessionDelete.Execute(ctx, deleteUserID.Value())
	if err != nil {
		err = log.WrapErrorWithStackTrace(err)
	}

	return err
}
