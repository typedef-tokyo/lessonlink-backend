package lessonlist

import (
	"context"
)

type QueryLessonDTO struct {
	ID             int
	LessonName     string
	LessonDuration int
}

type LessonListQueryRepository interface {
	GetListByCampus(ctx context.Context, campus string) ([]*QueryLessonDTO, error)
}
