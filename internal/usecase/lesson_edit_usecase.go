package usecase

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type (
	ILessonEditInputPort interface {
		Execute(c echo.Context, role vo.RoleKey, userID vo.UserID, input LessonEditInputDTO) error
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
		txManager        util.TxManager
		repositoryLesson repository.LessonRepository
	}
)

func NewLessonEditInteractor(
	txManager util.TxManager,
	repositoryLesson repository.LessonRepository,
) ILessonEditInputPort {
	return &LessonEditInteractor{
		txManager:        txManager,
		repositoryLesson: repositoryLesson,
	}
}

func (r LessonEditInteractor) Execute(c echo.Context, role vo.RoleKey, userID vo.UserID, input LessonEditInputDTO) error {

	if !role.IsOwner() {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	var errs error

	var lessonID vo.LessonID
	var lessonName vo.LessonName
	var duration vo.LessonDuration

	errs = errors.Join(errs, vo.SetVOConstructor(&lessonID, vo.NewLessonID, input.ID))
	errs = errors.Join(errs, vo.SetVOConstructor(&lessonName, vo.NewLessonName, input.LessonName))
	errs = errors.Join(errs, vo.SetVOConstructor(&duration, vo.NewLessonDuration, input.Duration))

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
