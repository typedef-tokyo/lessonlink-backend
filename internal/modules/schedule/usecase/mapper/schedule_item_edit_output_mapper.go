package mapper

import (
	"cmp"
	"errors"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/usecase/public"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/port"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type LessonDTO struct {
	id       vo.LessonID
	name     vo.LessonName
	duration vo.LessonDuration
}

func NewlessonDTO(
	id vo.LessonID,
	name vo.LessonName,
	duration vo.LessonDuration,
) *LessonDTO {

	return &LessonDTO{
		id:       id,
		name:     name,
		duration: duration,
	}
}

type LessonsDTOSlice []*LessonDTO

func NewLessonDTOSlice(lessons []public.LessonGetOutDTO) (LessonsDTOSlice, error) {

	lessonsDTO := make([]*LessonDTO, 0, len(lessons))
	for _, lesson := range lessons {

		var errs error

		var lessonID vo.LessonID
		var lessonName vo.LessonName
		var duration vo.LessonDuration

		errs = errors.Join(errs, utility.SetVOConstructor(&lessonID, vo.NewLessonID, lesson.ID))
		errs = errors.Join(errs, utility.SetVOConstructor(&lessonName, vo.NewLessonName, lesson.Name))
		errs = errors.Join(errs, utility.SetVOConstructor(&duration, vo.NewLessonDuration, lesson.Duration))

		if errs != nil {
			return nil, log.WrapErrorWithStackTraceBadRequest(log.Errorf("%v", errs.Error()))
		}

		lessonsDTO = append(lessonsDTO, NewlessonDTO(
			lessonID,
			lessonName,
			duration,
		))
	}

	return lessonsDTO, nil
}

func (r LessonsDTOSlice) findByID(id vo.LessonID) *LessonDTO {

	dto, _ := lo.Find(r, func(item *LessonDTO) bool {
		return item.id == id
	})

	return dto
}

type ScheduleItemEditOutputMapper struct{}

func NewScheduleItemEditOutputMapper() ScheduleItemEditOutputMapper {
	return ScheduleItemEditOutputMapper{}
}

func (m *ScheduleItemEditOutputMapper) ToScheduleItemEditOutput(
	scheduleData *schedule.RootScheduleModel,
	lessons LessonsDTOSlice,
) port.ScheduleItemEditOutputDTO {

	return port.ScheduleItemEditOutputDTO{
		HistoryIndex:   scheduleData.HistoryIndex().Value(),
		LessonItemList: m.BuildScheduleLessonItems(scheduleData, lessons),
		RoomLessonList: m.BuildScheduleRoomLessonItems(scheduleData, lessons),
	}
}

func (r ScheduleItemEditOutputMapper) BuildScheduleLessonItems(
	scheduleData *schedule.RootScheduleModel,
	lessons LessonsDTOSlice,
) []port.ScheduleLessonItem {

	savedLessonIDs := lo.Uniq(append(scheduleData.RoomItems().LessonIDs(), scheduleData.Items().LessonIDs()...))

	filterdLessons := lo.Filter(lessons, func(lesson *LessonDTO, _ int) bool {
		return !lo.Contains(savedLessonIDs, lesson.id)
	})

	items := make([]port.ScheduleLessonItem, 0, len(filterdLessons)+len(scheduleData.Items()))
	for _, lessonData := range filterdLessons {
		items = append(items, r.newScheduleLessonItem(scheduleData, lessonData)...)
	}

	for _, item := range scheduleData.Items() {

		lessonName := ""
		lesson := lessons.findByID(item.LessonID())
		if lesson != nil {
			lessonName = lesson.name.Value()
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
	lessonData *LessonDTO,
) []port.ScheduleLessonItem {

	existingSlice := scheduleData.FilterScheduleItemByLessonID(lessonData.id)
	if len(existingSlice) > 0 {

		return lo.Map(existingSlice, func(item *schedule.ScheduleItemModel, _ int) port.ScheduleLessonItem {
			return port.ScheduleLessonItem{
				LessonID:   item.LessonID().Value(),
				Identifier: item.Identifier().Value(),
				LessonName: lessonData.name.Value(),
				Duration:   item.Duration().Value(),
			}
		})
	}

	return []port.ScheduleLessonItem{
		{
			LessonID:   lessonData.id.Value(),
			Identifier: vo.NewIdentifierGenerate().Value(),
			LessonName: lessonData.name.Value(),
			Duration:   lessonData.duration.Value(),
		},
	}
}

func (r ScheduleItemEditOutputMapper) BuildScheduleRoomLessonItems(
	scheduleData *schedule.RootScheduleModel,
	lessons LessonsDTOSlice,
) []port.ScheduleRoomLesson {

	return lo.Map(scheduleData.RoomItems(), func(item *schedule.ScheduleRoomItemModel, _ int) port.ScheduleRoomLesson {
		startTimeHour, startTimeMinutes := item.StartTime().Value()
		endTimeHour, endTimeMinutes := item.EndTime().Value()

		lessonName := ""
		lesson := lessons.findByID(item.LessonID())
		if lesson != nil {
			lessonName = lesson.name.Value()
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
