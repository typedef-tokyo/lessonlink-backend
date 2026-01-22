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

type IScheduleSaveTitleInputPort interface {
	Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputTitle string) error
}

type (
	ScheduleSaveTitleInteractor struct {
		txManager                     util.TxManager
		repositorySchedule            repository.ScheduleRepository
		repositoryUser                repository.UserRepository
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleSaveTitleInteractor(
	txManager util.TxManager,
	repositorySchedule repository.ScheduleRepository,
	repositoryUser repository.UserRepository,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleSaveTitleInputPort {
	return &ScheduleSaveTitleInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		repositoryUser:                repositoryUser,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}
func (r ScheduleSaveTitleInteractor) Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputTitle string) error {

	scheduleID, title, err := r.createVO(inputScheduleID, inputTitle)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		scheduleData, err := r.getSchedule(ctx, tx, scheduleID, user)
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

func (r ScheduleSaveTitleInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, user vo.UserID) (*schedule.RootScheduleModel, error) {

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

func (ScheduleSaveTitleInteractor) createVO(inputScheduleID int, inputTitle string) (vo.ScheduleID, vo.ScheduleTitle, error) {

	var scheduleID vo.ScheduleID
	var title vo.ScheduleTitle

	var errs error
	errs = errors.Join(errs, vo.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))
	errs = errors.Join(errs, vo.SetVOConstructor(&title, vo.NewScheduleTitle, inputTitle))

	if errs != nil {
		return scheduleID, title, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return scheduleID, title, nil
}
