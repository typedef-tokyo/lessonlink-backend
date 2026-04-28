package external

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/usecase/public"
)

type (
	ILessonGetFacade interface {
		Execute(ctx context.Context, campus string) (*public.LessonsGetOutDTO, error)
	}
)
