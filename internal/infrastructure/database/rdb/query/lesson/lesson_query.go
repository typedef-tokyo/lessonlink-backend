package lesson

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/lessonlist"
)

type LessonQuery struct {
	c *sql.DB
}

func NewLessonQueryRepository(c rdb.IMySQL) lessonlist.LessonListQueryRepository {
	return &LessonQuery{c: c.GetConn()}
}

func (f *LessonQuery) GetListByCampus(ctx context.Context, campus string) ([]*lessonlist.QueryLessonDTO, error) {

	lessonRecords, err := dto.DataLessons(
		dto.DataLessonWhere.Campus.EQ(campus),
	).All(ctx, f.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return f.toList(lessonRecords), nil
}

func (f *LessonQuery) toList(LessonDTOs []*dto.DataLesson) []*lessonlist.QueryLessonDTO {

	lessonList := make([]*lessonlist.QueryLessonDTO, 0, len(LessonDTOs))

	for _, LessonDTO := range LessonDTOs {

		Lesson := &lessonlist.QueryLessonDTO{
			ID:             LessonDTO.ID,
			LessonName:     LessonDTO.Name,
			LessonDuration: LessonDTO.Duration,
		}

		lessonList = append(lessonList, Lesson)
	}

	return lessonList
}
