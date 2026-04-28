package interactor

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/service"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
)

type (
	IScheduleDeleteInputPort interface {
		Execute(ctx context.Context, scheduleID int, inputDeleteUserID int) error
	}
)

type (
	ScheduleDeleteInteractor struct {
		txManager                     database.TxManager
		repositorySchedule            repository.ScheduleRepository
		facadeGetUser                 external.IUserGetFacade
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleDeleteInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
	facadeGetUser external.IUserGetFacade,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleDeleteInputPort {
	return &ScheduleDeleteInteractor{
		txManager:                     txManager,
		facadeGetUser:                 facadeGetUser,
		repositorySchedule:            repositorySchedule,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleDeleteInteractor) Execute(ctx context.Context, inputScheduleID int, inputDeleteUserID int) error {

	deleteUserID, err := vo.NewUserID(inputDeleteUserID)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	scheduleID, err := vo.NewScheduleID(inputScheduleID)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	schedule, err := r.repositorySchedule.FindByID(ctx, scheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	if schedule == nil {
		return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したIDのスケジュールは存在しません:%d", scheduleID.Value()))
	}

	_, outDeleteUserRole, err := r.facadeGetUser.Execute(ctx, deleteUserID.Value())
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	deleteUserRole, err := vo.NewRoleKey(outDeleteUserRole)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	isEnable := r.serviceScheduleEditPermission.AllowsEditingBy(schedule, deleteUserID, deleteUserRole)
	if !isEnable {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		err = r.repositorySchedule.Delete(ctx, tx, scheduleID, deleteUserID)
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
