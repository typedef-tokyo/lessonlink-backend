package interactor

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/external"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/mapper"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/port"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	IScheduleGetInputPort interface {
		Execute(ctx context.Context, inputScheduleID int, intputHistoryIndex int) (*ScheduleGetOutput, error)
	}
)

type (
	ScheduleGetOutput struct {
		ScheduleID     int
		Campus         string
		Title          string
		ScheduleTime   ScheduleTimeDTO
		HistoryIndex   int
		Rooms          []ScheduleRoomDTO
		LessonItemList []port.ScheduleLessonItem
		RoomLessonList []port.ScheduleRoomLesson
		CreatedUserID  int
	}

	ScheduleTimeDTO struct {
		StartTime int
		EndTime   int
	}

	ScheduleRoomDTO struct {
		RoomIndex int
		RoomName  string
		Visible   bool
	}

	ScheduleLessonItem struct {
		LessonID   int
		Identifier string
		LessonName string
		Duration   int
	}
)

type (
	ScheduleGetInteractor struct {
		repositorySchedule              repository.ScheduleRepository
		repositoryScheduleInvisibleRoom repository.ScheduleInvisibleRoomRepository
		facadeRoomGet                   external.IRoomGetFacade
		facadeLessonGet                 external.ILessonGetFacade
		mapperScheduleItemOutput        mapper.ScheduleItemEditOutputMapper
	}
)

func NewScheduleGetInteractor(
	repositorySchedule repository.ScheduleRepository,
	repositoryScheduleInvisibleRoom repository.ScheduleInvisibleRoomRepository,
	facadeRoomGet external.IRoomGetFacade,
	facadeLessonGet external.ILessonGetFacade,
	mapperScheduleItemOutput mapper.ScheduleItemEditOutputMapper,
) IScheduleGetInputPort {
	return &ScheduleGetInteractor{
		repositorySchedule:              repositorySchedule,
		repositoryScheduleInvisibleRoom: repositoryScheduleInvisibleRoom,
		facadeRoomGet:                   facadeRoomGet,
		facadeLessonGet:                 facadeLessonGet,
		mapperScheduleItemOutput:        mapperScheduleItemOutput,
	}
}

func (r ScheduleGetInteractor) Execute(ctx context.Context, inputScheduleID int, intputHistoryIndex int) (*ScheduleGetOutput, error) {

	scheduleID, err := vo.NewScheduleID(inputScheduleID)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	historyIndex, err := vo.NewHistoryIndex(intputHistoryIndex)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	var scheduleData *schedule.RootScheduleModel

	var setHistoryIndex vo.HistoryIndex
	if historyIndex.IsUseLatest() {
		scheduleData, err = r.repositorySchedule.FindByID(ctx, scheduleID)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		if scheduleData == nil {
			return nil, log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したIDのスケジュールは存在しません:%d", scheduleID.Value()))
		}

		setHistoryIndex = scheduleData.HistoryIndex()
	} else {
		scheduleData, err = r.repositorySchedule.FindByIDWithHistoryIndex(ctx, scheduleID, historyIndex)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		if scheduleData == nil {
			return nil, log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したIDのスケジュールは存在しません:%d", scheduleID.Value()))
		}
		setHistoryIndex = historyIndex
	}

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	lessons, err := r.getLessons(ctx, scheduleData.Campus())
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	roomsDTO, err := r.getRooms(ctx, scheduleData)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	startTime, endTime := scheduleData.ScheduleTime().Value()
	scheduleTIme := ScheduleTimeDTO{
		StartTime: startTime,
		EndTime:   endTime,
	}

	return &ScheduleGetOutput{
		ScheduleID:     scheduleData.ID().Value(),
		Campus:         scheduleData.Campus().Value(),
		Title:          scheduleData.Title().Value(),
		ScheduleTime:   scheduleTIme,
		HistoryIndex:   setHistoryIndex.Value(),
		Rooms:          roomsDTO,
		LessonItemList: r.mapperScheduleItemOutput.BuildScheduleLessonItems(scheduleData, lessons),
		RoomLessonList: r.mapperScheduleItemOutput.BuildScheduleRoomLessonItems(scheduleData, lessons),
		CreatedUserID:  scheduleData.CreateUser().Value(),
	}, nil
}

func (r ScheduleGetInteractor) getLessons(ctx context.Context, campus vo.Campus) (mapper.LessonsDTOSlice, error) {

	lessons, err := r.facadeLessonGet.Execute(ctx, campus.Value())
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	lessonsDTOSlice, err := mapper.NewLessonDTOSlice(lessons.Lessons)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return lessonsDTOSlice, nil
}

func (r ScheduleGetInteractor) getRooms(ctx context.Context, scheduelData *schedule.RootScheduleModel) ([]ScheduleRoomDTO, error) {

	room, err := r.facadeRoomGet.Execute(ctx, scheduelData.Campus().Value())
	if err != nil {
		return []ScheduleRoomDTO{}, log.WrapErrorWithStackTrace(err)
	}

	invisibleRooms, err := r.repositoryScheduleInvisibleRoom.FindBySheduleID(ctx, scheduelData.ID())
	if err != nil {
		return []ScheduleRoomDTO{}, log.WrapErrorWithStackTrace(err)
	}

	scheduleRoomDTO := make([]ScheduleRoomDTO, 0, len(room.Rooms))
	for _, item := range room.Rooms {
		roomIndex, err := vo.NewRoomIndex(item.Index)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		scheduleRoomDTO = append(scheduleRoomDTO, ScheduleRoomDTO{
			RoomIndex: item.Index,
			RoomName:  item.Name,
			Visible:   !invisibleRooms.IsInvisible(roomIndex),
		})
	}

	return scheduleRoomDTO, nil
}
