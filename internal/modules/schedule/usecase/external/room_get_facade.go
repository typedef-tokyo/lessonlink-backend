package external

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/usecase/public"
)

type (
	IRoomGetFacade interface {
		Execute(ctx context.Context, campus string) (*public.RoomsGetOutDTO, error)
	}
)
