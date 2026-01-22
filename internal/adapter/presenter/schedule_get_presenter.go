package presenter

import (
	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/port"
)

type IScheduleGetPresenter interface {
	Present(result *usecase.ScheduleGetOutput) *ScheduleGetResponse
}

type ScheduleGetPresenter struct {
}

func NewScheduleGet(
	env configs.EnvConfig,
) IScheduleGetPresenter {
	return &ScheduleGetPresenter{}
}

type (
	ScheduleGetResponse struct {
		ScheduleID        int                  `json:"schedule_id"`
		Campus            string               `json:"campus"`
		Title             string               `json:"title"`
		ScheduleStartTime int                  `json:"schedule_start_time"`
		ScheduleEndTime   int                  `json:"schedule_end_time"`
		HistoryIndex      int                  `json:"history_index"`
		Rooms             []ScheduleRoomDTO    `json:"rooms"`
		LessonItemList    []ScheduleLessonItem `json:"lesson_item_list"`
		RoomLessonList    []ScheduleRoomLesson `json:"room_lesson_list"`
		CreatedUserID     int                  `json:"created_user_id"`
	}

	ScheduleRoomDTO struct {
		RoomIndex int    `json:"room_index"`
		RoomName  string `json:"room_name"`
		Visible   bool   `json:"visible"`
	}

	ScheduleLessonItem struct {
		LessonID   int    `json:"lesson_id"`
		Identifier string `json:"identifier"`
		LessonName string `json:"lesson_name"`
		Duration   int    `json:"duration"`
	}

	ScheduleRoomLesson struct {
		ItemTag          string `json:"item_tag"`
		LessonID         int    `json:"lesson_id"`
		Identifier       string `json:"identifier"`
		LessonName       string `json:"lesson_name"`
		Duration         int    `json:"duration"`
		StartTimeHour    int    `json:"start_time_hour"`
		StartTimeMinutes int    `json:"start_time_minutes"`
		EndTimeHour      int    `json:"end_time_hour"`
		EndTimeMinutes   int    `json:"end_time_minutes"`
		RoomIndex        int    `json:"room_index"`
	}
)

func (h *ScheduleGetPresenter) Present(result *usecase.ScheduleGetOutput) *ScheduleGetResponse {

	return &ScheduleGetResponse{
		ScheduleID:        result.ScheduleID,
		Campus:            result.Campus,
		Title:             result.Title,
		ScheduleStartTime: result.ScheduleTime.StartTime,
		ScheduleEndTime:   result.ScheduleTime.EndTime,
		HistoryIndex:      result.HistoryIndex,
		Rooms: lo.Map(result.Rooms, func(item usecase.ScheduleRoomDTO, _ int) ScheduleRoomDTO {
			return ScheduleRoomDTO{
				RoomIndex: item.RoomIndex,
				RoomName:  item.RoomName,
				Visible:   item.Visible,
			}
		}),
		LessonItemList: lo.Map(result.LessonItemList, func(item port.ScheduleLessonItem, _ int) ScheduleLessonItem {
			return ScheduleLessonItem{
				LessonID:   item.LessonID,
				Identifier: item.Identifier,
				LessonName: item.LessonName,
				Duration:   item.Duration,
			}
		}),
		RoomLessonList: lo.Map(result.RoomLessonList, func(item port.ScheduleRoomLesson, _ int) ScheduleRoomLesson {
			return ScheduleRoomLesson{
				ItemTag:          item.ItemTag,
				LessonID:         item.LessonID,
				Identifier:       item.Identifier,
				LessonName:       item.LessonName,
				Duration:         item.Duration,
				StartTimeHour:    item.StartTime.ScheduleItemTimeHour,
				StartTimeMinutes: item.StartTime.ScheduleItemTimeMinutes,
				EndTimeHour:      item.EndTime.ScheduleItemTimeHour,
				EndTimeMinutes:   item.EndTime.ScheduleItemTimeMinutes,
				RoomIndex:        item.RoomIndex,
			}
		}),
		CreatedUserID: result.CreatedUserID,
	}
}
