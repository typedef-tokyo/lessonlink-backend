package rdb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/lesson"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type Lesson struct {
	c *sql.DB
}

func NewLessonRepository(c IMySQL) repository.LessonRepository {
	return &Lesson{c: c.GetConn()}
}

func (f *Lesson) FindByCampus(ctx context.Context, campus vo.Campus) (lesson.RootLessonModelSlice, error) {

	records, err := dto.DataLessons(
		dto.DataLessonWhere.Campus.EQ(campus.Value()),
	).All(ctx, f.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	models := make([]*lesson.RootLessonModel, 0, len(records))
	for _, record := range records {

		model, err := f.toModel(record)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		models = append(models, model)
	}

	return models, nil
}

func (f *Lesson) FindByID(ctx context.Context, id vo.LessonID) (*lesson.RootLessonModel, error) {

	record, err := dto.DataLessons(
		dto.DataLessonWhere.ID.EQ(id.Value()),
	).One(ctx, f.c)

	if err != nil && err != sql.ErrNoRows {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	if record == nil {
		return nil, nil
	}

	model, err := f.toModel(record)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return model, nil
}

func (f *Lesson) Save(ctx context.Context, tx *sql.Tx, lesson *lesson.RootLessonModel) error {

	lessonDTO := f.toDTO(lesson)

	var err error
	if lessonDTO.ID == ID_INITIAL {
		err = lessonDTO.Insert(ctx, tx, boil.Infer())
	} else {
		err = lessonDTO.Upsert(ctx, tx, boil.Infer(), boil.Infer())
	}

	if err != nil {
		return log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return nil
}

func (f *Lesson) toModel(record *dto.DataLesson) (*lesson.RootLessonModel, error) {

	var id vo.LessonID
	var campus vo.Campus
	var name vo.LessonName
	var duration vo.LessonDuration

	var errs error
	errs = errors.Join(errs, vo.SetVOConstructor(&id, vo.NewLessonID, record.ID))
	errs = errors.Join(errs, vo.SetVOConstructor(&campus, vo.NewCampus, record.Campus))
	errs = errors.Join(errs, vo.SetVOConstructor(&name, vo.NewLessonName, record.Name))
	errs = errors.Join(errs, vo.SetVOConstructor(&duration, vo.NewLessonDuration, record.Duration))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("%v", errs.Error()))
	}

	return lesson.NewRootLessonModel(
		id,
		campus,
		name,
		duration,
	), nil

}

func (f *Lesson) toDTO(model *lesson.RootLessonModel) *dto.DataLesson {

	return &dto.DataLesson{
		ID:       model.ID().Value(),
		Campus:   model.Campus().Value(),
		Name:     model.Name().Value(),
		Duration: model.Duration().Value(),
	}
}
