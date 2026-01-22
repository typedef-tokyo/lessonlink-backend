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
	IScheduleItemShiftInputPort interface {
		Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputHistoryIndex int, inputRoomIndex int) (*port.ScheduleItemEditOutput, error)
	}
)

type (
	ScheduleItemShiftInteractor struct {
		txManager                     util.TxManager
		repositorySchedule            repository.ScheduleRepository
		repositoryUser                repository.UserRepository
		repositoryLesson              repository.LessonRepository
		mapperScheduleItemEditOutput  mapper.ScheduleItemEditOutputMapper
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleItemShiftInteractor(
	txManager util.TxManager,
	repositorySchedule repository.ScheduleRepository,
	repositoryUser repository.UserRepository,
	repositoryLesson repository.LessonRepository,
	mapperScheduleItemEditOutput mapper.ScheduleItemEditOutputMapper,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleItemShiftInputPort {
	return &ScheduleItemShiftInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		repositoryUser:                repositoryUser,
		repositoryLesson:              repositoryLesson,
		mapperScheduleItemEditOutput:  mapperScheduleItemEditOutput,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleItemShiftInteractor) Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputHistoryIndex int, inputRoomIndex int) (*port.ScheduleItemEditOutput, error) {

	scheduleID, historyIndex, roomIndex, err := r.createVO(inputScheduleID, inputHistoryIndex, inputRoomIndex)
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

		err = scheduleData.RoomItemShift(roomIndex)
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

func (ScheduleItemShiftInteractor) createVO(inputScheduleID int, inputHistoryIndex int, inputRoomIndex int) (vo.ScheduleID, vo.HistoryIndex, vo.RoomIndex, error) {

	var scheduleID vo.ScheduleID
	var historyIndex vo.HistoryIndex
	var roomIndex vo.RoomIndex

	var errs error
	errs = errors.Join(errs, vo.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))
	errs = errors.Join(errs, vo.SetVOConstructor(&historyIndex, vo.NewHistoryIndex, inputHistoryIndex))
	errs = errors.Join(errs, vo.SetVOConstructor(&roomIndex, vo.NewRoomIndex, inputRoomIndex))

	if errs != nil {
		return scheduleID, historyIndex, roomIndex, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return scheduleID, historyIndex, roomIndex, nil
}

func (r ScheduleItemShiftInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, user vo.UserID) (*schedule.RootScheduleModel, error) {

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
