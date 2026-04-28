package lessonlist

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	ILessonListInputPort interface {
		Execute(ctx context.Context, inputCampus string) (*LessonListQueryOutput, error)
	}
)

type (
	LessonListQueryOutput struct {
		LessonList []*QueryLessonDTO
	}
)

type LessonListQueryInteractor struct {
	repositroryLessonListQuery LessonListQueryRepository
}

func NewLessonListQueryInteractor(
	repositroryLessonListQuery LessonListQueryRepository,
) ILessonListInputPort {
	return &LessonListQueryInteractor{
		repositroryLessonListQuery: repositroryLessonListQuery,
	}
}

func (r LessonListQueryInteractor) Execute(ctx context.Context, inputCampus string) (*LessonListQueryOutput, error) {

	lessonData, err := r.repositroryLessonListQuery.GetListByCampus(ctx, inputCampus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &LessonListQueryOutput{
		LessonList: lessonData,
	}, nil
}
