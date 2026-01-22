package usecase

import (
	"context"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/room"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/mapper"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/port"
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
		repositoryRoom                  repository.RoomRepository
		repositoryScheduleInvisibleRoom repository.ScheduleInvisibleRoomRepository
		repositoryLesson                repository.LessonRepository
		mapperScheduleItemOutput        mapper.ScheduleItemEditOutputMapper
	}
)

func NewScheduleGetInteractor(
	repositorySchedule repository.ScheduleRepository,
	repositoryRoom repository.RoomRepository,
	repositoryScheduleInvisibleRoom repository.ScheduleInvisibleRoomRepository,
	repositoryLesson repository.LessonRepository,
	mapperScheduleItemOutput mapper.ScheduleItemEditOutputMapper,
) IScheduleGetInputPort {
	return &ScheduleGetInteractor{
		repositorySchedule:              repositorySchedule,
		repositoryRoom:                  repositoryRoom,
		repositoryScheduleInvisibleRoom: repositoryScheduleInvisibleRoom,
		repositoryLesson:                repositoryLesson,
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
		setHistoryIndex = scheduleData.HistoryIndex()
	} else {
		scheduleData, err = r.repositorySchedule.FindByIDWithHistoryIndex(ctx, scheduleID, historyIndex)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}
		setHistoryIndex = historyIndex
	}

	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	lessons, err := r.repositoryLesson.FindByCampus(ctx, scheduleData.Campus())
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

func (r ScheduleGetInteractor) getRooms(ctx context.Context, scheduelData *schedule.RootScheduleModel) ([]ScheduleRoomDTO, error) {

	rooms, err := r.repositoryRoom.FindByCampus(ctx, scheduelData.Campus())
	if err != nil {
		return []ScheduleRoomDTO{}, log.WrapErrorWithStackTrace(err)
	}

	invisibleRooms, err := r.repositoryScheduleInvisibleRoom.FindBySheduleID(ctx, scheduelData.ID())
	if err != nil {
		return []ScheduleRoomDTO{}, log.WrapErrorWithStackTrace(err)
	}

	return lo.Map(rooms, func(item *room.RootRoomModel, _ int) ScheduleRoomDTO {
		return ScheduleRoomDTO{
			RoomIndex: item.RoomIndex().Value(),
			RoomName:  item.RoomName().Value(),
			Visible:   !invisibleRooms.IsInvisible(item.RoomIndex()),
		}
	}), nil
}
