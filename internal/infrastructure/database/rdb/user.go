package rdb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/user"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type User struct {
	c *sql.DB
}

func NewUserRepository(c IMySQL) repository.UserRepository {
	return &User{c: c.GetConn()}
}

func (c *User) Save(ctx context.Context, tx *sql.Tx, user *user.RootUserModel, userID vo.UserID) error {

	dtoObject := c.toDTO(user, userID, ACTIVE)

	var err error

	// 新規登録の場合
	if user.ID() == vo.USER_ID_INITIAL {
		err = dtoObject.Insert(ctx, tx, boil.Infer())
	} else {

		// 更新の場合
		_, err = dtoObject.Update(ctx, tx, boil.Whitelist(
			dto.TBLUserColumns.RoleKey,
			dto.TBLUserColumns.Password,
			dto.TBLUserColumns.Name,
			dto.TBLUserColumns.UpdateUserID,
		))
	}
	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return nil
}

func (c *User) Delete(ctx context.Context, tx *sql.Tx, user *user.RootUserModel, deleteFromUserID vo.UserID) error {

	dtoObject := c.toDTO(user, deleteFromUserID, IN_ACTIVE)

	_, err := dtoObject.Update(ctx, tx, boil.Whitelist(
		dto.TBLUserColumns.UpdateUserID,
		dto.TBLUserColumns.DeleteFlag,
	))

	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return nil
}

func (c *User) FindAll(ctx context.Context) (user.RootUserModelSlice, error) {

	usersDTOs, err := dto.TBLUsers(
		dto.TBLUserWhere.DeleteFlag.EQ(ACTIVE),
		qm.OrderBy(dto.TBLUserColumns.ID+" desc"),
	).All(ctx, c.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	models := make([]*user.RootUserModel, 0, len(usersDTOs))

	for _, dto := range usersDTOs {

		model, err := c.toModel(dto)

		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		models = append(models, model)
	}

	return models, nil
}

func (c *User) FindByUserName(ctx context.Context, userName string) (*user.RootUserModel, error) {

	usersDTO, err := dto.TBLUsers(
		dto.TBLUserWhere.UserName.EQ(userName),
		dto.TBLUserWhere.DeleteFlag.EQ(ACTIVE),
	).All(ctx, c.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	if len(usersDTO) == 0 {
		return nil, nil
	}

	model, err := c.toModel(usersDTO[0])

	if err != nil {

		return nil, log.WrapErrorWithStackTrace(err)
	}

	return model, nil
}

func (c *User) FindByUserID(ctx context.Context, userId vo.UserID) (*user.RootUserModel, error) {

	usersDTOs, err := dto.TBLUsers(
		dto.TBLUserWhere.ID.EQ(userId.Value()),
		dto.TBLUserWhere.DeleteFlag.EQ(ACTIVE),
	).All(ctx, c.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	if len(usersDTOs) == 0 {
		return nil, nil
	}

	model, err := c.toModel(usersDTOs[0])

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return model, nil

}

func (c *User) FindByUserIDs(ctx context.Context, userIds []vo.UserID) (user.RootUserModelSlice, error) {

	ids := lo.Map(userIds, func(item vo.UserID, _ int) int {
		return item.Value()
	})

	usersDTOs, err := dto.TBLUsers(
		dto.TBLUserWhere.ID.IN(ids),
		dto.TBLUserWhere.DeleteFlag.EQ(ACTIVE),
	).All(ctx, c.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	userModels := make([]*user.RootUserModel, len(usersDTOs))
	for index, dto := range usersDTOs {

		model, err := c.toModel(dto)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		userModels[index] = model
	}

	return userModels, nil
}

func (d *User) toModel(record *dto.TBLUser) (*user.RootUserModel, error) {

	var errs error

	var id vo.UserID
	var roleKey vo.RoleKey
	var userName vo.UserName
	var password vo.UserPassword
	var displayName vo.UserDisplayName

	errs = errors.Join(errs, vo.SetVOConstructor(&id, vo.NewUserID, record.ID))
	errs = errors.Join(errs, vo.SetVOConstructor(&roleKey, vo.NewRoleKey, record.RoleKey))
	errs = errors.Join(errs, vo.SetVOConstructor(&userName, vo.NewUserName, record.UserName))
	errs = errors.Join(errs, vo.SetVOConstructor(&password, vo.ReconstructHashedPassword, record.Password))
	errs = errors.Join(errs, vo.SetVOConstructor(&displayName, vo.NewUserDisplayName, record.Name))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("%v", errs.Error()))
	}

	return user.NewRootUserModel(
		id,
		roleKey,
		userName,
		password,
		displayName,
	), nil
}

func (d *User) toDTO(user *user.RootUserModel, userID vo.UserID, deleteFlag int) dto.TBLUser {

	return dto.TBLUser{
		ID:           user.ID().Value(),
		RoleKey:      user.RoleKey().Value(),
		UserName:     user.UserName().Value(),
		Password:     user.Password().Value(),
		Name:         user.DisplayName().Value(),
		UpdateUserID: userID.Value(),
		DeleteFlag:   deleteFlag,
	}
}
