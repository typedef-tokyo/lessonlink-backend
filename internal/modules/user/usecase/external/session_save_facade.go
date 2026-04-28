package external

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/usecase/public"
)

type (
	ISessionSaveFacade interface {
		Execute(ctx context.Context, sessionDTO public.SessionSaveInputDTO) error
	}
)
