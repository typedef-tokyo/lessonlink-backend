package public

import (
	"context"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/domain/model/role"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	RolesGetOutDTO struct {
		Roles []RoleGetOutDTO
	}

	RoleGetOutDTO struct {
		RoleKey  string
		RoleName string
	}
)

///////////////

type (
	RoleGetFacade struct {
		repositoryRole repository.RoleRepository
	}
)

func NewRoleGetFacade(
	repositoryRole repository.RoleRepository,
) *RoleGetFacade {
	return &RoleGetFacade{
		repositoryRole: repositoryRole,
	}
}

func (r RoleGetFacade) Execute(ctx context.Context) (*RolesGetOutDTO, error) {

	roles, err := r.repositoryRole.FindAll(ctx)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &RolesGetOutDTO{
		Roles: lo.Map(roles, func(item *role.RootRoleModel, _ int) RoleGetOutDTO {
			return RoleGetOutDTO{
				RoleKey:  item.RoleKey().Value(),
				RoleName: item.RoleName().Value(),
			}
		}),
	}, nil
}
