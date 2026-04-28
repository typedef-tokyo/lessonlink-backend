package interactor

import (
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/domain/model/lesson"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type (
	ILessonAddInputPort interface {
		Execute(ctx context.Context, inputRole string, inputCampus string, input LessonAddInputDTO) error
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
		txManager        database.TxManager
		facadeCampus     external.ICampusExistFacade
		repositoryLesson repository.LessonRepository
	}
)

func NewLessonAddInteractor(
	txManager database.TxManager,
	facadeCampus external.ICampusExistFacade,
	repositoryLesson repository.LessonRepository,
) ILessonAddInputPort {
	return &LessonAddInteractor{
		txManager:        txManager,
		facadeCampus:     facadeCampus,
		repositoryLesson: repositoryLesson,
	}
}

func (r LessonAddInteractor) Execute(ctx context.Context, inputRole string, inputCampus string, input LessonAddInputDTO) error {

	role, err := vo.NewRoleKey(inputRole)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	if !role.IsOwner() {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	campus, err := vo.NewCampus(inputCampus)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	campusExist, err := r.facadeCampus.Execute(ctx, inputCampus)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if !campusExist {
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

	errs = errors.Join(errs, utility.SetVOConstructor(&lessonName, vo.NewLessonName, input.LessonName))
	errs = errors.Join(errs, utility.SetVOConstructor(&duration, vo.NewLessonDuration, input.Duration))

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
