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
	IScheduleItemReturnListInputPort interface {
		Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int, inputData ScheduleItemReturnListInput) (*port.ScheduleItemEditOutput, error)
	}
)

type (
	ScheduleItemReturnListInput struct {
		LessonID   int
		Identifier string
		Duration   int
	}
)

type (
	ScheduleItemReturnListInteractor struct {
		txManager                     database.TxManager
		repositorySchedule            repository.ScheduleRepository
		facadeGetUser                 external.IUserGetFacade
		facadeLessonGet               external.ILessonGetFacade
		mapperScheduleItemEditOutput  mapper.ScheduleItemEditOutputMapper
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleItemReturnListInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
	facadeGetUser external.IUserGetFacade,
	facadeLessonGet external.ILessonGetFacade,
	mapperScheduleItemEditOutput mapper.ScheduleItemEditOutputMapper,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,

) IScheduleItemReturnListInputPort {
	return &ScheduleItemReturnListInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		facadeGetUser:                 facadeGetUser,
		facadeLessonGet:               facadeLessonGet,
		mapperScheduleItemEditOutput:  mapperScheduleItemEditOutput,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleItemReturnListInteractor) Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int, inputData ScheduleItemReturnListInput) (*port.ScheduleItemEditOutput, error) {

	userID, scheduleID, historyIndex, err := r.createVO(inputUserID, inputScheduleID, inputHistoryIndex)
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

		returnItem, err := r.createReturnItem(inputData)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		err = scheduleData.ItemReturnList(returnItem)
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

func (r ScheduleItemReturnListInteractor) createReturnItem(inputData ScheduleItemReturnListInput) (*schedule.ScheduleItemModel, error) {

	var lessonID vo.LessonID
	var identifier vo.Identifier
	var duration vo.LessonDuration

	var errs error
	errs = errors.Join(errs, utility.SetVOConstructor(&lessonID, vo.NewLessonID, inputData.LessonID))
	errs = errors.Join(errs, utility.SetVOConstructor(&identifier, vo.NewIdentifier, inputData.Identifier))
	errs = errors.Join(errs, utility.SetVOConstructor(&duration, vo.NewLessonDuration, inputData.Duration))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return schedule.NewScheduleItemModel(
		lessonID,
		identifier,
		duration,
	), nil

}

func (r ScheduleItemReturnListInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, editUserID vo.UserID) (*schedule.RootScheduleModel, error) {

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

func (ScheduleItemReturnListInteractor) createVO(inputUserID int, inputScheduleID int, inputHistoryIndex int) (vo.UserID, vo.ScheduleID, vo.HistoryIndex, error) {

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
