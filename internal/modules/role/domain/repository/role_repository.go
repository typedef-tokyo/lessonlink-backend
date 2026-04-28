package repository

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/domain/model/role"
)

type RoleRepository interface {
	FindAll(ctx context.Context) (role.RootRoleModelSlice, error)
}
