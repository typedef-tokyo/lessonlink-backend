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
	IScheduleItemDivideInputPort interface {
		Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputHistoryIndex int, inputDivide ScheduleItemDivideInput) (*port.ScheduleItemEditOutput, error)
	}
)

type (
	ScheduleItemDivideInput struct {
		LessonID      int
		Identifier    string
		DivideMinutes int
	}
)

type (
	ScheduleItemDivideInteractor struct {
		txManager                     util.TxManager
		repositorySchedule            repository.ScheduleRepository
		repositoryUser                repository.UserRepository
		repositoryLesson              repository.LessonRepository
		mapperScheduleItemEditOutput  mapper.ScheduleItemEditOutputMapper
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleItemDivideInteractor(
	txManager util.TxManager,
	repositorySchedule repository.ScheduleRepository,
	repositoryUser repository.UserRepository,
	repositoryLesson repository.LessonRepository,
	mapperScheduleItemEditOutput mapper.ScheduleItemEditOutputMapper,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleItemDivideInputPort {
	return &ScheduleItemDivideInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		repositoryUser:                repositoryUser,
		repositoryLesson:              repositoryLesson,
		mapperScheduleItemEditOutput:  mapperScheduleItemEditOutput,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleItemDivideInteractor) Execute(ctx context.Context, role vo.RoleKey, user vo.UserID, inputScheduleID int, inputHistoryIndex int, inputDivide ScheduleItemDivideInput) (*port.ScheduleItemEditOutput, error) {

	scheduleID, historyIndex, lessonID, identifier, divideMinutes, err := r.createVO(inputScheduleID, inputHistoryIndex, inputDivide)
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

		if scheduleData == nil {
			return log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したIDのスケジュールは存在しません:%d", scheduleID.Value()))
		}

		lessons, err = r.repositoryLesson.FindByCampus(ctx, scheduleData.Campus())
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		lesson := lessons.FindByID(lessonID)
		if lesson == nil {
			return log.WrapErrorWithStackTraceBadRequest(errors.New("分割対象の講座は登録されていません"))
		}

		err = scheduleData.ItemDivide(lessonID, lesson.Duration(), identifier, divideMinutes)
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

func (ScheduleItemDivideInteractor) createVO(inputScheduleID int, inputHistoryIndex int, inputData ScheduleItemDivideInput) (vo.ScheduleID, vo.HistoryIndex, vo.LessonID, vo.Identifier, vo.ItemDivideMinutes, error) {

	var scheduleID vo.ScheduleID
	var historyIndex vo.HistoryIndex
	var lessonID vo.LessonID
	var identifier vo.Identifier
	var divideMinutes vo.ItemDivideMinutes

	var errs error
	errs = errors.Join(errs, vo.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))
	errs = errors.Join(errs, vo.SetVOConstructor(&historyIndex, vo.NewHistoryIndex, inputHistoryIndex))
	errs = errors.Join(errs, vo.SetVOConstructor(&lessonID, vo.NewLessonID, inputData.LessonID))
	errs = errors.Join(errs, vo.SetVOConstructor(&identifier, vo.NewIdentifier, inputData.Identifier))
	errs = errors.Join(errs, vo.SetVOConstructor(&divideMinutes, vo.NewItemDivideMinutes, inputData.DivideMinutes))

	if errs != nil {
		return scheduleID, historyIndex, lessonID, identifier, divideMinutes, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return scheduleID, historyIndex, lessonID, identifier, divideMinutes, nil
}

func (r ScheduleItemDivideInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, user vo.UserID) (*schedule.RootScheduleModel, error) {

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
