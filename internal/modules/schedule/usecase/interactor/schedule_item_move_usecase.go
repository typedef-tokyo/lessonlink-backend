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
	IScheduleItemMoveInputPort interface {
		Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int, inputData ScheduleItemMoveInput) (*port.ScheduleItemEditOutput, error)
	}
)

type (
	ScheduleItemMoveInput struct {
		LessonID        int
		ItemTag         string
		Identifier      string
		Duration        int
		StartTimeHour   int
		StartTimeMinute int
		EndTimeHour     int
		EndTimeMinutes  int
		RoomIndex       int
	}
)

type (
	ScheduleItemMoveInteractor struct {
		txManager                     database.TxManager
		repositorySchedule            repository.ScheduleRepository
		facadeGetUser                 external.IUserGetFacade
		facadeLessonGet               external.ILessonGetFacade
		mapperScheduleItemEditOutput  mapper.ScheduleItemEditOutputMapper
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleItemMoveInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
	facadeLessonGet external.ILessonGetFacade,
	facadeGetUser external.IUserGetFacade,
	mapperScheduleItemEditOutput mapper.ScheduleItemEditOutputMapper,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleItemMoveInputPort {
	return &ScheduleItemMoveInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		facadeGetUser:                 facadeGetUser,
		facadeLessonGet:               facadeLessonGet,
		mapperScheduleItemEditOutput:  mapperScheduleItemEditOutput,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleItemMoveInteractor) Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int, inputData ScheduleItemMoveInput) (*port.ScheduleItemEditOutput, error) {

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

		moveItem, err := r.createNewMoveItem(inputData)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		err = scheduleData.RoomItemMove(moveItem)
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

func (r ScheduleItemMoveInteractor) createNewMoveItem(inputData ScheduleItemMoveInput) (*schedule.ScheduleRoomItemModel, error) {

	var itemTag vo.RoomItemTag
	var lessonID vo.LessonID
	var identifier vo.Identifier
	var duration vo.LessonDuration
	var roomIndex vo.RoomIndex

	var errs error
	errs = errors.Join(errs, utility.SetVOConstructor(&itemTag, vo.NewRoomItemTag, inputData.ItemTag))
	errs = errors.Join(errs, utility.SetVOConstructor(&lessonID, vo.NewLessonID, inputData.LessonID))
	errs = errors.Join(errs, utility.SetVOConstructor(&identifier, vo.NewIdentifier, inputData.Identifier))
	errs = errors.Join(errs, utility.SetVOConstructor(&duration, vo.NewLessonDuration, inputData.Duration))

	startTime, err := vo.NewScheduleLessonTime(inputData.StartTimeHour, inputData.StartTimeMinute)
	errs = errors.Join(errs, err)

	endTime, err := vo.NewScheduleLessonTime(inputData.EndTimeHour, inputData.EndTimeMinutes)
	errs = errors.Join(errs, err)

	errs = errors.Join(errs, utility.SetVOConstructor(&roomIndex, vo.NewRoomIndex, inputData.RoomIndex))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return schedule.NewScheduleRoomItemModel(
		itemTag,
		lessonID,
		identifier,
		duration,
		startTime,
		endTime,
		roomIndex,
	), nil

}

func (r ScheduleItemMoveInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, editUserID vo.UserID) (*schedule.RootScheduleModel, error) {

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

func (ScheduleItemMoveInteractor) createVO(inputUserID int, inputScheduleID int, inputHistoryIndex int) (vo.UserID, vo.ScheduleID, vo.HistoryIndex, error) {

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
