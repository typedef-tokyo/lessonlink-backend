package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/lesson"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type (
	ILessonAddInputPort interface {
		Execute(ctx context.Context, role vo.RoleKey, inputCampus string, input LessonAddInputDTO) error
	}
)

type (
	LessonAddInputDTO struct {
		LessonName string
		Duration   int
	}
)

type (
	LessonAddInteractor struct {
		txManager        util.TxManager
		repositoryCampus repository.CampusRepository
		repositoryLesson repository.LessonRepository
	}
)

func NewLessonAddInteractor(
	txManager util.TxManager,
	repositoryCampus repository.CampusRepository,
	repositoryLesson repository.LessonRepository,
) ILessonAddInputPort {
	return &LessonAddInteractor{
		txManager:        txManager,
		repositoryCampus: repositoryCampus,
		repositoryLesson: repositoryLesson,
	}
}

func (r LessonAddInteractor) Execute(ctx context.Context, role vo.RoleKey, inputCampus string, input LessonAddInputDTO) error {

	if !role.IsOwner() {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	campus, err := vo.NewCampus(inputCampus)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	campuses, err := r.repositoryCampus.FindAll(ctx)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if !campuses.IsExist(campus) {
		return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定した校舎はありません:%s", campus.Value()))
	}

	lessonModel, err := r.createModel(campus, input)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	lessons, err := r.repositoryLesson.FindByCampus(ctx, campus)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if lessons.CheckDuplicateEntry(lessonModel) {
		return log.WrapErrorWithStackTraceBadRequest(log.Errorf("同名の講座がすでに登録されています"))
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		if err := r.repositoryLesson.Save(ctx, tx, lessonModel); err != nil {
			return log.WrapErrorWithStackTraceInternalServerError(err)
		}

		return nil
	})

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil
}

func (r LessonAddInteractor) createModel(campus vo.Campus, input LessonAddInputDTO) (*lesson.RootLessonModel, error) {

	var errs error

	var lessonName vo.LessonName
	var duration vo.LessonDuration

	errs = errors.Join(errs, vo.SetVOConstructor(&lessonName, vo.NewLessonName, input.LessonName))
	errs = errors.Join(errs, vo.SetVOConstructor(&duration, vo.NewLessonDuration, input.Duration))

	if errs != nil {
		return nil, log.WrapErrorWithStackTrace(log.Errorf("%v", errs.Error()))
	}

	return lesson.NewRootLessonModel(
		vo.LESSON_ID_INITIAL,
		campus,
		lessonName,
		duration,
	), nil
}
