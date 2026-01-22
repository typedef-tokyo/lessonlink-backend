package presenter

import (
	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/port"
)

type IScheduleItemEditPresenter interface {
	Present(result *port.ScheduleItemEditOutput) *ScheduleItemEditResponse
}

type ScheduleItemEditPresenter struct {
}

func NewScheduleItemEditPresenter() IScheduleItemEditPresenter {
	return &ScheduleItemEditPresenter{}
}

type (
	ScheduleItemEditResponse struct {
		HistoryIndex   int                          `json:"history_index"`
		LessonItemList []ScheduleItemEditLessonItem `json:"lesson_item_list"`
		RoomLessonList []ScheduleItemEditRoomLesson `json:"room_lesson_list"`
	}

	ScheduleItemEditLessonItem struct {
		LessonID   int    `json:"lesson_id"`
		Identifier string `json:"identifier"`
		LessonName string `json:"lesson_name"`
		Duration   int    `json:"duration"`
	}

	ScheduleItemEditRoomLesson struct {
		ItemTag         string `json:"item_tag"`
		LessonID        int    `json:"lesson_id"`
		Identifier      string `json:"identifier"`
		LessonName      string `json:"lesson_name"`
		Duration        int    `json:"duration"`
		StartTimeHour   int    `json:"start_time_hour"`
		StartTimeMinute int    `json:"start_time_minutes"`
		EndTimeHour     int    `json:"end_time_hour"`
		EndTimeMinute   int    `json:"end_time_minutes"`
		RoomIndex       int    `json:"room_index"`
	}
)

func (h *ScheduleItemEditPresenter) Present(result *port.ScheduleItemEditOutput) *ScheduleItemEditResponse {

	scheduleItem := result.ScheduleItem
	return &ScheduleItemEditResponse{
		HistoryIndex: scheduleItem.HistoryIndex,
		LessonItemList: lo.Map(scheduleItem.LessonItemList, func(item port.ScheduleLessonItem, _ int) ScheduleItemEditLessonItem {
			return ScheduleItemEditLessonItem{
				LessonID:   item.LessonID,
				Identifier: item.Identifier,
				LessonName: item.LessonName,
				Duration:   item.Duration,
			}
		}),
		RoomLessonList: lo.Map(scheduleItem.RoomLessonList, func(item port.ScheduleRoomLesson, _ int) ScheduleItemEditRoomLesson {
			return ScheduleItemEditRoomLesson{
				ItemTag:         item.ItemTag,
				LessonID:        item.LessonID,
				Identifier:      item.Identifier,
				LessonName:      item.LessonName,
				Duration:        item.Duration,
				StartTimeHour:   item.StartTime.ScheduleItemTimeHour,
				StartTimeMinute: item.StartTime.ScheduleItemTimeMinutes,
				EndTimeHour:     item.EndTime.ScheduleItemTimeHour,
				EndTimeMinute:   item.EndTime.ScheduleItemTimeMinutes,
				RoomIndex:       item.RoomIndex,
			}
		}),
	}
}
