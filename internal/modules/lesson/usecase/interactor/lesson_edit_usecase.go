package interactor

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type (
	ILessonEditInputPort interface {
		Execute(c echo.Context, inputRole string, inputUserID int, input LessonEditInputDTO) error
	}
)

type (
	LessonEditInputDTO struct {
		ID         int
		LessonName string
		Duration   int
	}
)

type (
	LessonEditInteractor struct {
		txManager        database.TxManager
		repositoryLesson repository.LessonRepository
	}
)

func NewLessonEditInteractor(
	txManager database.TxManager,
	repositoryLesson repository.LessonRepository,
) ILessonEditInputPort {
	return &LessonEditInteractor{
		txManager:        txManager,
		repositoryLesson: repositoryLesson,
	}
}

func (r LessonEditInteractor) Execute(c echo.Context, inputRole string, inputUserID int, input LessonEditInputDTO) error {

	role, err := vo.NewRoleKey(inputRole)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	if !role.IsOwner() {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	var errs error

	var lessonID vo.LessonID
	var lessonName vo.LessonName
	var duration vo.LessonDuration

	errs = errors.Join(errs, utility.SetVOConstructor(&lessonID, vo.NewLessonID, input.ID))
	errs = errors.Join(errs, utility.SetVOConstructor(&lessonName, vo.NewLessonName, input.LessonName))
	errs = errors.Join(errs, utility.SetVOConstructor(&duration, vo.NewLessonDuration, input.Duration))

	if errs != nil {
		return log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	ctx := c.Request().Context()
	lessonModel, err := r.repositoryLesson.FindByID(ctx, lessonID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if lessonModel == nil {
		return log.WrapErrorWithStackTraceNotFound(log.Errorf("講座が見つかりません: %d", lessonID.Value()))
	}

	lessonModel.Revise(lessonName, duration)

	lessons, err := r.repositoryLesson.FindByCampus(ctx, lessonModel.Campus())
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
