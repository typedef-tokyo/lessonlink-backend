package interactor

import (
	"context"
	"database/sql"
	"errors"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/usecase/public"
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
	IScheduleItemDivideInputPort interface {
		Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int, inputDivide ScheduleItemDivideInput) (*port.ScheduleItemEditOutput, error)
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
		txManager                     database.TxManager
		repositorySchedule            repository.ScheduleRepository
		facadeGetUser                 external.IUserGetFacade
		facadeLessonGet               external.ILessonGetFacade
		mapperScheduleItemEditOutput  mapper.ScheduleItemEditOutputMapper
		serviceScheduleEditPermission service.IScheduleEditPermissionService
	}
)

func NewScheduleItemDivideInteractor(
	txManager database.TxManager,
	repositorySchedule repository.ScheduleRepository,
	facadeGetUser external.IUserGetFacade,
	facadeLessonGet external.ILessonGetFacade,
	mapperScheduleItemEditOutput mapper.ScheduleItemEditOutputMapper,
	serviceScheduleEditPermission service.IScheduleEditPermissionService,
) IScheduleItemDivideInputPort {
	return &ScheduleItemDivideInteractor{
		txManager:                     txManager,
		repositorySchedule:            repositorySchedule,
		facadeGetUser:                 facadeGetUser,
		facadeLessonGet:               facadeLessonGet,
		mapperScheduleItemEditOutput:  mapperScheduleItemEditOutput,
		serviceScheduleEditPermission: serviceScheduleEditPermission,
	}
}

func (r ScheduleItemDivideInteractor) Execute(ctx context.Context, inputUserID int, inputScheduleID int, inputHistoryIndex int, inputDivide ScheduleItemDivideInput) (*port.ScheduleItemEditOutput, error) {

	userID, scheduleID, historyIndex, lessonID, identifier, divideMinutes, err := r.createVO(inputUserID, inputScheduleID, inputHistoryIndex, inputDivide)
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

		lesson, found := lo.Find(lessons.Lessons, func(item public.LessonGetOutDTO) bool {
			return item.ID == lessonID.Value()
		})

		if !found {
			return log.WrapErrorWithStackTraceBadRequest(errors.New("分割対象の講座は登録されていません"))
		}

		lessonDuration, err := vo.NewLessonDuration(lesson.Duration)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		lessonsDTOSlice, err = mapper.NewLessonDTOSlice(lessons.Lessons)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		err = scheduleData.ItemDivide(lessonID, lessonDuration, identifier, divideMinutes)
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

func (ScheduleItemDivideInteractor) createVO(inputUserID int, inputScheduleID int, inputHistoryIndex int, inputData ScheduleItemDivideInput) (vo.UserID, vo.ScheduleID, vo.HistoryIndex, vo.LessonID, vo.Identifier, vo.ItemDivideMinutes, error) {

	var userID vo.UserID
	var scheduleID vo.ScheduleID
	var historyIndex vo.HistoryIndex
	var lessonID vo.LessonID
	var identifier vo.Identifier
	var divideMinutes vo.ItemDivideMinutes

	var errs error
	errs = errors.Join(errs, utility.SetVOConstructor(&userID, vo.NewUserID, inputUserID))
	errs = errors.Join(errs, utility.SetVOConstructor(&scheduleID, vo.NewScheduleID, inputScheduleID))
	errs = errors.Join(errs, utility.SetVOConstructor(&historyIndex, vo.NewHistoryIndex, inputHistoryIndex))
	errs = errors.Join(errs, utility.SetVOConstructor(&lessonID, vo.NewLessonID, inputData.LessonID))
	errs = errors.Join(errs, utility.SetVOConstructor(&identifier, vo.NewIdentifier, inputData.Identifier))
	errs = errors.Join(errs, utility.SetVOConstructor(&divideMinutes, vo.NewItemDivideMinutes, inputData.DivideMinutes))

	if errs != nil {
		return userID, scheduleID, historyIndex, lessonID, identifier, divideMinutes, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
	}

	return userID, scheduleID, historyIndex, lessonID, identifier, divideMinutes, nil
}

func (r ScheduleItemDivideInteractor) getSchedule(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex, editUserID vo.UserID) (*schedule.RootScheduleModel, error) {

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
