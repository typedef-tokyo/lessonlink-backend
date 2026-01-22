package lessonlist

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	campusRepository "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/campus"
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
	repositroryCampusQuery     campusRepository.CampusQueryRepository
	repositroryLessonListQuery LessonListQueryRepository
}

func NewLessonListQueryInteractor(
	repositroryCampusQuery campusRepository.CampusQueryRepository,
	repositroryLessonListQuery LessonListQueryRepository,
) ILessonListInputPort {
	return &LessonListQueryInteractor{
		repositroryCampusQuery:     repositroryCampusQuery,
		repositroryLessonListQuery: repositroryLessonListQuery,
	}
}

func (r LessonListQueryInteractor) Execute(ctx context.Context, inputCampus string) (*LessonListQueryOutput, error) {

	campus, err := r.repositroryCampusQuery.GetByCampus(ctx, inputCampus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if campus == nil {
		return nil, log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したキャンパスはありません:%s", inputCampus))
	}

	lessonData, err := r.repositroryLessonListQuery.GetListByCampus(ctx, inputCampus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &LessonListQueryOutput{
		LessonList: lessonData,
	}, nil
}
