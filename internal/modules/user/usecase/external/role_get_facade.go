package external

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/usecase/public"
)

type (
	IRoleGetFacade interface {
		Execute(ctx context.Context) (*public.RolesGetOutDTO, error)
	}
)
