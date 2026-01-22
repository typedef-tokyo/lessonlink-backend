package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/lesson"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/service"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/mapper"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/port"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/util"
)

type (
	IScheduleItemReturnListInputPort interface {
		Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputHistoryIndex int, inputData ScheduleItemReturnListInput) (*port.ScheduleItemEditOutput, error)
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
		txManager                     util.TxManager
		repositorySchedule            repository.ScheduleRepository
		repositoryUser                repository.UserRepository
		repositoryLesson              repository.LessonRepository
		mapperScheduleItemEditOutput  mapper.ScheduleItemEditOutputMapper
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleItemReturnListInteractor(
	txManager util.TxManager,
	repositorySchedule repository.ScheduleRepository,
	repositoryUser repository.UserRepository,
	repositoryLesson repository.LessonRepository,
	mapperScheduleItemEditOutput mapper.ScheduleItemEditOutputMapper,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,

) IScheduleItemReturnListInputPort {
	return &ScheduleItemReturnListInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		repositoryUser:                repositoryUser,
		repositoryLesson:              repositoryLesson,
		mapperScheduleItemEditOutput:  mapperScheduleItemEditOutput,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleItemReturnListInteractor) Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputHistoryIndex int, inputData ScheduleItemReturnListInput) (*port.ScheduleItemEditOutput, error) {

	scheduleID, historyIndex, err := r.createVO(inputScheduleID, inputHistoryIndex)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	var scheduleData *schedule.RootScheduleModel
	var lessons lesson.RootLessonModelSlice
	err = r.txManager.Do(ctx, func(tx *sql.Tx) error {

		var err error
		scheduleData, err = r.getSchedule(ctx, tx, scheduleID, historyIndex, user)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		lessons, err = r.repositoryLesson.FindByCampus(ctx, scheduleData.Campus())
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

		scheduleData.ModifyEditing(historyIndex, user)

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
		ScheduleItem: r.mapperScheduleItemEditOutput.ToScheduleItemEditOutput(scheduleData, lessons),
	}, nil
}

func (r ScheduleItemReturnListInteractor) createReturnItem(inputData ScheduleItemReturnListInput) (*schedule.ScheduleItemModel, error) {

	var lessonID vo.LessonID
	var identifier vo.Identifier
	var duration vo.LessonDuration

	var errs error
	errs = errors.Join(errs, vo.SetVOConstructor(&lessonID, vo.NewLessonID, inputData.LessonID))
	errs = errors.Join(errs, vo.SetVOConstructor(&identifier, vo.NewIdentifier, inputData.Identifier))
	errs = errors.Join(errs, vo.SetVOConstructor(&duration, vo.NewLessonDuration, inputData.Duration))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return schedule.NewScheduleItemModel(
		lessonID,
		identifier,
		duration,
	), nil

}

func (r ScheduleItemReturnListInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, user vo.UserID) (*schedule.RootScheduleModel, error) {

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

func (ScheduleItemReturnListInteractor) createVO(inputScheduleID int, inputHistoryIndex int) (vo.ScheduleID, vo.HistoryIndex, error) {

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
