package repository

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/lesson"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type LessonRepository interface {
	FindByCampus(ctx context.Context, campus vo.Campus) (lesson.RootLessonModelSlice, error)
	FindByID(ctx context.Context, id vo.LessonID) (*lesson.RootLessonModel, error)
	Save(ctx context.Context, tx *sql.Tx, lesson *lesson.RootLessonModel) error
}
