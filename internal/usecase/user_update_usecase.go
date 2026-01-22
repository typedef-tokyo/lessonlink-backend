package usecase

import (
	"cmp"
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/user"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/hash"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type IUserUpdateInputPort interface {
	Execute(ctx context.Context, user UserUpdateInput, updateFromUserID vo.UserID, roleKey vo.RoleKey) error
}

type UserUpdateInput struct {
	UserID      int
	RoleKey     string
	UserName    string
	Password    string
	DisplayName string
}

type UserUpdateInteractor struct {
	txManager      util.TxManager
	repositoryUser repository.UserRepository
}

func NewUserUpdateInteractor(
	txManager util.TxManager,
	repositoryUser repository.UserRepository,
) IUserUpdateInputPort {
	return &UserUpdateInteractor{
		txManager:      txManager,
		repositoryUser: repositoryUser,
	}
}
func (r UserUpdateInteractor) Execute(ctx context.Context, inputUser UserUpdateInput, updateFromUserID vo.UserID, roleKey vo.RoleKey) error {

	userID, err := vo.NewUserID(inputUser.UserID)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	// 対象のユーザーを取得
	userData, err := r.repositoryUser.FindByUserID(ctx, userID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if userData == nil {
		return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したユーザーは存在しません:%d", userID.Value()))
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		if err = r.updateUserModel(inputUser, userData, updateFromUserID, roleKey); err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		// 登録を実行
		if err = r.repositoryUser.Save(ctx, tx, userData, updateFromUserID); err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil
}

func (r UserUpdateInteractor) updateUserModel(
	inputUser UserUpdateInput,
	userModel *user.RootUserModel,
	updateFromUserID vo.UserID,
	updateUserRoleKey vo.RoleKey,
) error {

	var roleKey vo.RoleKey
	var userName vo.UserName
	var displayName vo.UserDisplayName

	var errs error

	errs = errors.Join(errs, vo.SetVOConstructor(&roleKey, vo.NewRoleKey, inputUser.RoleKey))
	errs = errors.Join(errs, vo.SetVOConstructor(&userName, vo.NewUserName, inputUser.UserName))

	var encryptPassword vo.UserPassword
	if inputUser.Password != "" {

		password, err := vo.NewPasswordForCreation(inputUser.Password)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		// パスワードをハッシュ化
		_encryptPassword, err := hash.HashPassword(password.Value())
		if err != nil {
			return log.WrapErrorWithStackTraceInternalServerError(err)
		}

		encryptPassword, _ = vo.ReconstructHashedPassword(_encryptPassword)
	}

	password := cmp.Or(encryptPassword, vo.USER_PASSWORD_NONE)

	errs = errors.Join(errs, vo.SetVOConstructor(&displayName, vo.NewUserDisplayName, inputUser.DisplayName))

	if errs != nil {
		return log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	err := userModel.UpdateUser(
		roleKey,
		password,
		displayName,
		updateFromUserID,
		updateUserRoleKey,
	)

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil
}
