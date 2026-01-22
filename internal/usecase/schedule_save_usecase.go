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
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type IScheduleSaveInputPort interface {
	Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputHistoryIndex int) (*ScheduleSaveOutput, error)
}

type ScheduleSaveOutput struct {
	HistoryIndex int
}

type (
	ScheduleSaveInteractor struct {
		txManager                     util.TxManager
		repositorySchedule            repository.ScheduleRepository
		repositoryUser                repository.UserRepository
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleSaveInteractor(
	txManager util.TxManager,
	repositorySchedule repository.ScheduleRepository,
	repositoryUser repository.UserRepository,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleSaveInputPort {
	return &ScheduleSaveInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		repositoryUser:                repositoryUser,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}
func (r ScheduleSaveInteractor) Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputHistoryIndex int) (*ScheduleSaveOutput, error) {

	scheduleID, historyIndex, err := r.createVO(inputScheduleID, inputHistoryIndex)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	var scheduleData *schedule.RootScheduleModel
	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		var err error
		scheduleData, err = r.getSchedule(ctx, tx, scheduleID, historyIndex, user)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		scheduleData.ModifySaving(historyIndex, user)

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

func (r ScheduleSaveInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, user vo.UserID) (*schedule.RootScheduleModel, error) {

	scheduleData, err := r.repositorySchedule.FindByIDWithLockHistoryIndex(ctx, tx, scheduleID, historyIndex)
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

	return scheduleData, nil
}

func (ScheduleSaveInteractor) createVO(inputScheduleID int, inputHistoryIndex int) (vo.ScheduleID, vo.HistoryIndex, error) {

	var scheduleID vo.ScheduleID
	var historyIndex vo.HistoryIndex

	var errs error
	errs = errors.Join(errs, vo.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))
	errs = errors.Join(errs, vo.SetVOConstructor(&historyIndex, vo.NewHistoryIndex, inputHistoryIndex))

	if errs != nil {
		return scheduleID, historyIndex, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return scheduleID, historyIndex, nil
}
