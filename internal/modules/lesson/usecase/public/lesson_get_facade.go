package public

import (
	"context"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/domain/model/lesson"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	LessonsGetOutDTO struct {
		Lessons []LessonGetOutDTO
	}

	LessonGetOutDTO struct {
		ID       int
		Name     string
		Duration int
	}
)

///////////////

type (
	LessonGetFacade struct {
		repositoryLesson repository.LessonRepository
	}
)

func NewLessonGetFacade(
	repositoryLesson repository.LessonRepository,
) *LessonGetFacade {
	return &LessonGetFacade{
		repositoryLesson: repositoryLesson,
	}
}

func (r LessonGetFacade) Execute(ctx context.Context, inputCampus string) (*LessonsGetOutDTO, error) {

	campus, err := vo.NewCampus(inputCampus)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	lessons, err := r.repositoryLesson.FindByCampus(ctx, campus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &LessonsGetOutDTO{
		Lessons: lo.Map(lessons, func(item *lesson.RootLessonModel, _ int) LessonGetOutDTO {
			return LessonGetOutDTO{
				ID:       item.ID().Value(),
				Name:     item.Name().Value(),
				Duration: item.Duration().Value(),
			}
		}),
	}, nil
}
