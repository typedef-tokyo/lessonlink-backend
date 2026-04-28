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
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/port"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type (
	IScheduleItemJoinInputPort interface {
		Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int, inputJoinFromIdentifier string, inputJoinToIdentifier string) (*port.ScheduleItemEditOutput, error)
	}
)

type (
	ScheduleItemJoinInteractor struct {
		txManager                     database.TxManager
		repositorySchedule            repository.ScheduleRepository
		facadeGetUser                 external.IUserGetFacade
		facadeLessonGet               external.ILessonGetFacade
		mapperScheduleItemEditOutput  mapper.ScheduleItemEditOutputMapper
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleItemJoinInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
	facadeGetUser external.IUserGetFacade,
	facadeLessonGet external.ILessonGetFacade,
	mapperScheduleItemEditOutput mapper.ScheduleItemEditOutputMapper,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleItemJoinInputPort {
	return &ScheduleItemJoinInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		facadeGetUser:                 facadeGetUser,
		facadeLessonGet:               facadeLessonGet,
		mapperScheduleItemEditOutput:  mapperScheduleItemEditOutput,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleItemJoinInteractor) Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int, inputJoinFromIdentifier string, inputJoinToIdentifier string) (*port.ScheduleItemEditOutput, error) {

	userID, scheduleID, historyIndex, joinFromIdentifier, joinToIdentifier, err := r.createVO(inputUserID, inputScheduleID, inputHistoryIndex, inputJoinFromIdentifier, inputJoinToIdentifier)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	var scheduleData *schedule.RootScheduleModel
	var lessonsDTOSlice mapper.LessonsDTOSlice
	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		var err error
		scheduleData, err = r.getSchedule(ctx, tx, scheduleID, historyIndex, userID)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		if scheduleData == nil {
			return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したIDのスケジュールは存在しません:%d", scheduleID.Value()))
		}

		lessons, err := r.facadeLessonGet.Execute(ctx, scheduleData.Campus().Value())
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		lessonsDTOSlice, err = mapper.NewLessonDTOSlice(lessons.Lessons)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		err = scheduleData.ItemJoin(joinFromIdentifier, joinToIdentifier)
		if err != nil {
			return log.WrapErrorWithStackTraceBadRequest(err)
		}

		scheduleData.ModifyEditing(historyIndex, userID)

		_, err = r.repositorySchedule.Save(ctx, tx, scheduleData)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	})

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &port.ScheduleItemEditOutput{
		ScheduleItem: r.mapperScheduleItemEditOutput.ToScheduleItemEditOutput(scheduleData, lessonsDTOSlice),
	}, nil
}

func (r ScheduleItemJoinInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, editUserID vo.UserID) (*schedule.RootScheduleModel, error) {

	scheduleData, err := r.repositorySchedule.FindByIDWithLockHistoryIndex(ctx, tx, scheduleID, historyIndex)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
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

func (ScheduleItemJoinInteractor) createVO(inputUserID int, inputScheduleID int, inputHistoryIndex int, inputJoinFromIdentifier string, inputJoinToIdentifier string) (vo.UserID, vo.ScheduleID, vo.HistoryIndex, vo.Identifier, vo.Identifier, error) {

	var userID vo.UserID
	var scheduleID vo.ScheduleID
	var historyIndex vo.HistoryIndex
	var joinFromIdentifier vo.Identifier
	var joinToIdentifier vo.Identifier

	var errs error
	errs = errors.Join(errs, utility.SetVOConstructor(&userID, vo.NewUserID, inputUserID))
	errs = errors.Join(errs, utility.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))
	errs = errors.Join(errs, utility.SetVOConstructor(&historyIndex, vo.NewHistoryIndex, inputHistoryIndex))
	errs = errors.Join(errs, utility.SetVOConstructor(&joinFromIdentifier, vo.NewIdentifier, inputJoinFromIdentifier))
	errs = errors.Join(errs, utility.SetVOConstructor(&joinToIdentifier, vo.NewIdentifier, inputJoinToIdentifier))

	if errs != nil {
		return userID, scheduleID, historyIndex, joinFromIdentifier, joinToIdentifier, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return userID, scheduleID, historyIndex, joinFromIdentifier, joinToIdentifier, nil
}
