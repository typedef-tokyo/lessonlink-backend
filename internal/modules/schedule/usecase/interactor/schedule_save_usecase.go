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
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type IScheduleSaveInputPort interface {
	Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int) (*ScheduleSaveOutput, error)
}

type ScheduleSaveOutput struct {
	HistoryIndex int
}

type (
	ScheduleSaveInteractor struct {
		txManager                     database.TxManager
		repositorySchedule            repository.ScheduleRepository
		facadeGetUser                 external.IUserGetFacade
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleSaveInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
	facadeGetUser external.IUserGetFacade,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleSaveInputPort {
	return &ScheduleSaveInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		facadeGetUser:                 facadeGetUser,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}
func (r ScheduleSaveInteractor) Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int) (*ScheduleSaveOutput, error) {

	userID, scheduleID, historyIndex, err := r.createVO(inputUserID, inputScheduleID, inputHistoryIndex)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	var scheduleData *schedule.RootScheduleModel
	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		var err error
		scheduleData, err = r.getSchedule(ctx, tx, scheduleID, historyIndex, userID)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		scheduleData.ModifySaving(historyIndex, userID)

		_, err = r.repositorySchedule.Save(ctx, tx, scheduleData)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &ScheduleSaveOutput{HistoryIndex: scheduleData.HistoryIndex().Value()}, nil
}

func (r ScheduleSaveInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, editUserID vo.UserID) (*schedule.RootScheduleModel, error) {

	scheduleData, err := r.repositorySchedule.FindByIDWithLockHistoryIndex(ctx, tx, scheduleID, historyIndex)
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

	return scheduleData, nil
}

func (ScheduleSaveInteractor) createVO(inputUserID int, inputScheduleID int, inputHistoryIndex int) (vo.UserID, vo.ScheduleID, vo.HistoryIndex, error) {

	var userID vo.UserID
	var scheduleID vo.ScheduleID
	var historyIndex vo.HistoryIndex

	var errs error
	errs = errors.Join(errs, utility.SetVOConstructor(&userID, vo.NewUserID, inputUserID))
	errs = errors.Join(errs, utility.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))
	errs = errors.Join(errs, utility.SetVOConstructor(&historyIndex, vo.NewHistoryIndex, inputHistoryIndex))

	if errs != nil {
		return userID, scheduleID, historyIndex, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return userID, scheduleID, historyIndex, nil
}
