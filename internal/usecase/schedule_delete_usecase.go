package usecase

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/service"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type (
	IScheduleDeleteInputPort interface {
		Execute(ctx context.Context, role vo.RoleKey, scheduleID int, inputDeleteUserID vo.UserID) error
	}
)

type (
	ScheduleDeleteInteractor struct {
		txManager                     util.TxManager
		repositorySchedule            repository.ScheduleRepository
		repositoryUser                repository.UserRepository
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleDeleteInteractor(
	txManager util.TxManager,
	repositorySchedule repository.ScheduleRepository,
	repositoryUser repository.UserRepository,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleDeleteInputPort {
	return &ScheduleDeleteInteractor{
		txManager:                     txManager,
		repositoryUser:                repositoryUser,
		repositorySchedule:            repositorySchedule,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleDeleteInteractor) Execute(ctx context.Context, role vo.RoleKey, inputScheduleID int, inputDeleteUserID vo.UserID) error {

	scheduleID, err := vo.NewScheduleID(inputScheduleID)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	schedule, err := r.repositorySchedule.FindByID(ctx, scheduleID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	deleteUser, err := r.repositoryUser.FindByUserID(ctx, inputDeleteUserID)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	isEnable := r.serviceScheduleEditPermission.AllowsEditingBy(schedule, deleteUser)
	if !isEnable {
		return log.WrapErrorWithStackTraceForbidden(log.Errorf("許可されていない操作です"))
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		err = r.repositorySchedule.Delete(ctx, tx, scheduleID, inputDeleteUserID)
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
