package mapper

import (
	"cmp"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/lesson"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/port"
)

type ScheduleItemEditOutputMapper struct{}

func NewScheduleItemEditOutputMapper() ScheduleItemEditOutputMapper {
	return ScheduleItemEditOutputMapper{}
}

func (m *ScheduleItemEditOutputMapper) ToScheduleItemEditOutput(
	scheduleData *schedule.RootScheduleModel,
	lessons lesson.RootLessonModelSlice,
) port.ScheduleItemEditOutputDTO {

	return port.ScheduleItemEditOutputDTO{
		HistoryIndex:   scheduleData.HistoryIndex().Value(),
		LessonItemList: m.BuildScheduleLessonItems(scheduleData, lessons),
		RoomLessonList: m.BuildScheduleRoomLessonItems(scheduleData, lessons),
	}
}

func (r ScheduleItemEditOutputMapper) BuildScheduleLessonItems(
	scheduleData *schedule.RootScheduleModel,
	lessons lesson.RootLessonModelSlice,
) []port.ScheduleLessonItem {

	savedLessonIDs := lo.Uniq(append(scheduleData.RoomItems().LessonIDs(), scheduleData.Items().LessonIDs()...))

	filterdLessons := lo.Filter(lessons, func(lesson *lesson.RootLessonModel, _ int) bool {
		return !lo.Contains(savedLessonIDs, lesson.ID())
	})

	items := make([]port.ScheduleLessonItem, 0, len(filterdLessons)+len(scheduleData.Items()))
	for _, lessonData := range filterdLessons {
		items = append(items, r.newScheduleLessonItem(scheduleData, lessonData)...)
	}

	for _, item := range scheduleData.Items() {

		lessonName := ""
		lesson := lessons.FindByID(item.LessonID())
		if lesson != nil {
			lessonName = lesson.Name().Value()
		}

		items = append(items, port.ScheduleLessonItem{
			LessonID:   item.LessonID().Value(),
			Identifier: item.Identifier().Value(),
			LessonName: cmp.Or(lessonName, "不明な講座"),
			Duration:   item.Duration().Value(),
		})
	}

	return items
}

func (r ScheduleItemEditOutputMapper) newScheduleLessonItem(
	scheduleData *schedule.RootScheduleModel,
	lessonData *lesson.RootLessonModel,
) []port.ScheduleLessonItem {

	existingSlice := scheduleData.FilterScheduleItemByLessonID(lessonData.ID())
	if len(existingSlice) > 0 {

		return lo.Map(existingSlice, func(item *schedule.ScheduleItemModel, _ int) port.ScheduleLessonItem {
			return port.ScheduleLessonItem{
				LessonID:   item.LessonID().Value(),
				Identifier: item.Identifier().Value(),
				LessonName: lessonData.Name().Value(),
				Duration:   item.Duration().Value(),
			}
		})
	}

	return []port.ScheduleLessonItem{
		{
			LessonID:   lessonData.ID().Value(),
			Identifier: vo.NewIdentifierGenerate().Value(),
			LessonName: lessonData.Name().Value(),
			Duration:   lessonData.Duration().Value(),
		},
	}
}

func (r ScheduleItemEditOutputMapper) BuildScheduleRoomLessonItems(
	scheduleData *schedule.RootScheduleModel,
	lessons lesson.RootLessonModelSlice,
) []port.ScheduleRoomLesson {

	return lo.Map(scheduleData.RoomItems(), func(item *schedule.ScheduleRoomItemModel, _ int) port.ScheduleRoomLesson {
		startTimeHour, startTimeMinutes := item.StartTime().Value()
		endTimeHour, endTimeMinutes := item.EndTime().Value()

		lessonName := ""
		lesson := lessons.FindByID(item.LessonID())
		if lesson != nil {
			lessonName = lesson.Name().Value()
		} else if item.ItemTag().IsCleaning() {
			lessonName = "清掃"
		}

		return port.ScheduleRoomLesson{
			ItemTag:    item.ItemTag().Value(),
			LessonID:   item.LessonID().Value(),
			Identifier: item.Identifier().Value(),
			LessonName: cmp.Or(lessonName, "不明な講座"),
			Duration:   item.Duration().Value(),
			StartTime: port.ScheduleItemEditRoomLessonTime{
				ScheduleItemTimeHour:    startTimeHour,
				ScheduleItemTimeMinutes: startTimeMinutes,
			},
			EndTime: port.ScheduleItemEditRoomLessonTime{
				ScheduleItemTimeHour:    endTimeHour,
				ScheduleItemTimeMinutes: endTimeMinutes,
			},
			RoomIndex: item.RoomIndex().Value(),
		}
	})
}
