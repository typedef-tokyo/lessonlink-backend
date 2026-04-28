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

type IScheduleSaveTitleInputPort interface {
	Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputTitle string) error
}

type (
	ScheduleSaveTitleInteractor struct {
		txManager                     database.TxManager
		repositorySchedule            repository.ScheduleRepository
		facadeGetUser                 external.IUserGetFacade
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleSaveTitleInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
	facadeGetUser external.IUserGetFacade,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleSaveTitleInputPort {
	return &ScheduleSaveTitleInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		facadeGetUser:                 facadeGetUser,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}
func (r ScheduleSaveTitleInteractor) Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputTitle string) error {

	userID, scheduleID, title, err := r.createVO(inputUserID, inputScheduleID, inputTitle)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		scheduleData, err := r.getSchedule(ctx, tx, scheduleID, userID)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		scheduleData.ChangeTitle(title)

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

func (r ScheduleSaveTitleInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, editUserID vo.UserID) (*schedule.RootScheduleModel, error) {

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

func (ScheduleSaveTitleInteractor) createVO(inputUserID int, inputScheduleID int, inputTitle string) (vo.UserID, vo.ScheduleID, vo.ScheduleTitle, error) {

	var userID vo.UserID
	var scheduleID vo.ScheduleID
	var title vo.ScheduleTitle

	var errs error
	errs = errors.Join(errs, utility.SetVOConstructor(&userID, vo.NewUserID, inputUserID))
	errs = errors.Join(errs, utility.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))
	errs = errors.Join(errs, utility.SetVOConstructor(&title, vo.NewScheduleTitle, inputTitle))

	if errs != nil {
		return userID, scheduleID, title, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return userID, scheduleID, title, nil
}
