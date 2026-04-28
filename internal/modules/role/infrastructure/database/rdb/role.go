package rdb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/domain/model/role"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type Role struct {
	c *sql.DB
}

func NewRoleRepository(c rdb.IMySQL) repository.RoleRepository {
	return &Role{c: c.GetConn()}
}

func (c *Role) FindAll(ctx context.Context) (role.RootRoleModelSlice, error) {

	roleDTOs, err := dto.DataRoles().All(ctx, c.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	models := make([]*role.RootRoleModel, 0, len(roleDTOs))

	for _, dto := range roleDTOs {

		model, err := c.toModel(dto)

		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		models = append(models, model)
	}

	return models, nil
}

func (d *Role) toModel(record *dto.DataRole) (*role.RootRoleModel, error) {

	var id vo.RoleID
	var roleKey vo.RoleKey
	var roleName vo.RoleName

	var errs error
	errs = errors.Join(errs, utility.SetVOConstructor(&id, vo.NewRoleID, record.ID))
	errs = errors.Join(errs, utility.SetVOConstructor(&roleKey, vo.NewRoleKey, record.RoleKey))
	errs = errors.Join(errs, utility.SetVOConstructor(&roleName, vo.NewRoleName, record.RoleName))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("%v", errs.Error()))
	}

	return role.NewRootRoleModel(
		id,
		roleKey,
		roleName,
	), nil
}
