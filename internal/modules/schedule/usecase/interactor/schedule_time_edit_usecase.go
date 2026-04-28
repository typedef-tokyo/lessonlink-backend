package interactor

import (
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/service"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/mapper"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type IScheduleTimeEditInputPort interface {
	Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputScheduleStartTime int, inputScheduleEndTime int) error
}

type (
	ScheduleTimeEditInteractor struct {
		txManager                     database.TxManager
		repositorySchedule            repository.ScheduleRepository
		facadeGetUser                 external.IUserGetFacade
		mapperScheduleItemEditOutput  mapper.ScheduleItemEditOutputMapper
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleTimeEditEditInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
	facadeGetUser external.IUserGetFacade,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
	mapperScheduleItemEditOutput mapper.ScheduleItemEditOutputMapper,
) IScheduleTimeEditInputPort {
	return &ScheduleTimeEditInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		facadeGetUser:                 facadeGetUser,
		mapperScheduleItemEditOutput:  mapperScheduleItemEditOutput,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleTimeEditInteractor) Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputScheduleStartTime int, inputScheduleEndTime int) error {

	userID, scheduleID, scheduleTime, err := r.createVO(inputUserID, inputScheduleID, inputScheduleStartTime, inputScheduleEndTime)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	var scheduleData *schedule.RootScheduleModel
	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		var err error
		scheduleData, err = r.getSchedule(ctx, tx, scheduleID, userID)
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

func (r ScheduleTimeEditInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, editUserID vo.UserID) (*schedule.RootScheduleModel, error) {

	scheduleData, err := r.repositorySchedule.FindByIDWithLock(ctx, tx, scheduleID)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if scheduleData == nil {
		return nil, log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したIDのスケジュールは存在しません:%d", scheduleID.Value()))
	}

	_, outEditUserRole, err := r.facadeGetUser.Execute(ctx, editUserID.Value())
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	editUserRole, err := vo.NewRoleKey(outEditUserRole)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	isEnable := r.serviceScheduleEditPermission.AllowsEditingBy(scheduleData, editUserID, editUserRole)
	if !isEnable {
		return nil, log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	scheduleData, err = r.repositorySchedule.FindByIDWithLockHistoryIndex(ctx, tx, scheduleData.ID(), scheduleData.HistoryIndex())
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return scheduleData, nil
}

func (ScheduleTimeEditInteractor) createVO(inputUserID int, inputScheduleID int, inputScheduleStartTime int, inputScheduleEndTime int) (vo.UserID, vo.ScheduleID, vo.ScheduleTime, error) {

	var userID vo.UserID
	var scheduleID vo.ScheduleID

	var errs error
	errs = errors.Join(errs, utility.SetVOConstructor(&userID, vo.NewUserID, inputUserID))
	errs = errors.Join(errs, utility.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))

	scheduleTime, err := vo.NewScheduleTime(inputScheduleStartTime, inputScheduleEndTime)
	errs = errors.Join(errs, err)

	if errs != nil {
		return userID, scheduleID, scheduleTime, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return userID, scheduleID, scheduleTime, nil
}
