package usecase

import (
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

type IUserAddPort interface {
	Execute(ctx context.Context, inputAddUser UserAddInput, roleKey vo.RoleKey, addFromUserID vo.UserID) error
}

type UserAddInput struct {
	RoleKey     string
	UserName    string
	Password    string
	DisplayName string
}

type (
	UserAddInteractor struct {
		txManager      util.TxManager
		repositoryUser repository.UserRepository
	}
)

func NewUserAddInteractor(
	txManager util.TxManager,
	repositoryUser repository.UserRepository,
) IUserAddPort {
	return &UserAddInteractor{
		txManager:      txManager,
		repositoryUser: repositoryUser,
	}
}

func (r UserAddInteractor) Execute(ctx context.Context, inputAddUser UserAddInput, roleKey vo.RoleKey, addFromUserID vo.UserID) error {

	// ロールを確認
	if !roleKey.IsOwner() {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	// ユーザーモデルを作成する
	addUser, err := r.createUserModel(inputAddUser)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	// 既に登録済みでないか確認
	{
		user, err := r.repositoryUser.FindByUserName(ctx, addUser.UserName().Value())
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		// 既に登録されている場合
		if user != nil {
			return log.WrapErrorWithStackTraceConflict(log.Errorf("登録済みのユーザーです:%s", addUser.UserName()))
		}
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		err = r.repositoryUser.Save(ctx, tx, addUser, addFromUserID)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil
}

func (a UserAddInteractor) createUserModel(inputAddUser UserAddInput) (*user.RootUserModel, error) {

	var roleKey vo.RoleKey
	var userName vo.UserName
	var displayName vo.UserDisplayName

	var errs error

	errs = errors.Join(errs, vo.SetVOConstructor(&roleKey, vo.NewRoleKey, inputAddUser.RoleKey))
	errs = errors.Join(errs, vo.SetVOConstructor(&userName, vo.NewUserName, inputAddUser.UserName))

	password, err := vo.NewPasswordForCreation(inputAddUser.Password)
	errs = errors.Join(errs, err)

	errs = errors.Join(errs, vo.SetVOConstructor(&displayName, vo.NewUserDisplayName, inputAddUser.DisplayName))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	// パスワードをハッシュ化
	rawEncryptPassword, err := hash.HashPassword(password.Value())
	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	encryptPassword, _ := vo.ReconstructHashedPassword(rawEncryptPassword)

	return user.NewCreateUserModel(
		roleKey,
		userName,
		encryptPassword,
		displayName,
	), nil
}
