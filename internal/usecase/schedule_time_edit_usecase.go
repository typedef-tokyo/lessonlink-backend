package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/service"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/mapper"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type IScheduleTimeEditInputPort interface {
	Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputScheduleStartTime int, inputScheduleEndTime int) error
}

type (
	ScheduleTimeEditInteractor struct {
		txManager                     util.TxManager
		repositorySchedule            repository.ScheduleRepository
		repositoryUser                repository.UserRepository
		repositoryLesson              repository.LessonRepository
		mapperScheduleItemEditOutput  mapper.ScheduleItemEditOutputMapper
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleTimeEditEditInteractor(
	txManager util.TxManager,
	repositorySchedule repository.ScheduleRepository,
	repositoryUser repository.UserRepository,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
	repositoryLesson repository.LessonRepository,
	mapperScheduleItemEditOutput mapper.ScheduleItemEditOutputMapper,
) IScheduleTimeEditInputPort {
	return &ScheduleTimeEditInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		repositoryUser:                repositoryUser,
		repositoryLesson:              repositoryLesson,
		mapperScheduleItemEditOutput:  mapperScheduleItemEditOutput,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleTimeEditInteractor) Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputScheduleStartTime int, inputScheduleEndTime int) error {

	scheduleID, scheduleTime, err := r.createVO(inputScheduleID, inputScheduleStartTime, inputScheduleEndTime)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	var scheduleData *schedule.RootScheduleModel
	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		var err error
		scheduleData, err = r.getSchedule(ctx, tx, scheduleID, user)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		err = scheduleData.ChangeScheduleTime(scheduleTime)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		_, err = r.repositorySchedule.Save(ctx, tx, scheduleData)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	return nil

}

func (r ScheduleTimeEditInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, user vo.UserID) (*schedule.RootScheduleModel, error) {

	scheduleData, err := r.repositorySchedule.FindByIDWithLock(ctx, tx, scheduleID)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	editUser, err := r.repositoryUser.FindByUserID(ctx, user)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	isEnable := r.serviceScheduleEditPermission.AllowsEditingBy(scheduleData, editUser)
	if !isEnable {
		return nil, log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	scheduleData, err = r.repositorySchedule.FindByIDWithLockHistoryIndex(ctx, tx, scheduleData.ID(), scheduleData.HistoryIndex())
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return scheduleData, nil
}

func (ScheduleTimeEditInteractor) createVO(inputScheduleID int, inputScheduleStartTime int, inputScheduleEndTime int) (vo.ScheduleID, vo.ScheduleTime, error) {

	var scheduleID vo.ScheduleID

	var errs error
	errs = errors.Join(errs, vo.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))

	scheduleTime, err := vo.NewScheduleTime(inputScheduleStartTime, inputScheduleEndTime)
	errs = errors.Join(errs, err)

	if errs != nil {
		return scheduleID, scheduleTime, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return scheduleID, scheduleTime, nil
}
