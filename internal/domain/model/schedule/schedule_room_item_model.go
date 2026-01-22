package schedule

import (
	"cmp"
	"slices"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type ScheduleRoomItemModelSlice []*ScheduleRoomItemModel

func (r ScheduleRoomItemModelSlice) LessonIDs() []vo.LessonID {

	return lo.Map(r, func(item *ScheduleRoomItemModel, _ int) vo.LessonID {
		return item.lessonID
	})
}

func (r ScheduleRoomItemModelSlice) findByIdentifier(identifier vo.Identifier) (*ScheduleRoomItemModel, bool) {

	return lo.Find(r, func(item *ScheduleRoomItemModel) bool {
		return item.identifier == identifier
	})
}

func (r ScheduleRoomItemModelSlice) replaceItem(newItem *ScheduleRoomItemModel) ScheduleRoomItemModelSlice {

	removedItems := lo.Filter(r, func(item *ScheduleRoomItemModel, _ int) bool {
		return item.identifier != newItem.identifier
	})

	return append(removedItems, newItem)
}

func (r ScheduleRoomItemModelSlice) removeByIdentifier(removeItemID vo.Identifier) ScheduleRoomItemModelSlice {

	return lo.Filter(r, func(item *ScheduleRoomItemModel, _ int) bool {
		return item.identifier != removeItemID
	})
}

func (r ScheduleRoomItemModelSlice) removeByRoomIndex(roomIndex vo.RoomIndex) ScheduleRoomItemModelSlice {

	return lo.Filter(r, func(item *ScheduleRoomItemModel, _ int) bool {
		return item.roomIndex != roomIndex
	})
}

func (r ScheduleRoomItemModelSlice) shiftedItems(roomIndex vo.RoomIndex, scheduleTime vo.ScheduleTime) (ScheduleRoomItemModelSlice, error) {

	roomItems := lo.Map(
		lo.Filter(r, func(item *ScheduleRoomItemModel, _ int) bool {
			return item.roomIndex == roomIndex
		}),
		func(item *ScheduleRoomItemModel, _ int) *ScheduleRoomItemModel {
			copy := *item
			return &copy
		},
	)

	slices.SortFunc(roomItems, func(a, b *ScheduleRoomItemModel) int {
		return cmp.Compare(a.startTime.ValueMinutes(), b.startTime.ValueMinutes())
	})

	var createLessonTime = func(minutes int) (vo.ScheduleLessonTime, error) {

		newStartTimeHour := minutes / 60
		newStartTimeMinute := minutes % 60
		return vo.NewScheduleLessonTime(newStartTimeHour, newStartTimeMinute)

	}

	scheduleEndTime := scheduleTime.EndTimeValueMinutes()
	for index, roomItem := range roomItems {

		if index <= 0 {
			continue
		}

		start := roomItems[index-1].endTime.ValueMinutes()
		end := start + roomItem.duration.Value()

		if end > scheduleEndTime {

			diff := end - scheduleEndTime
			end -= diff
			start -= diff
		}

		newStartTime, err := createLessonTime(start)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		newEndTime, err := createLessonTime(end)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		roomItem.startTime = newStartTime
		roomItem.endTime = newEndTime
	}

	return roomItems, nil
}

type ScheduleRoomItemModel struct {
	itemTag    vo.RoomItemTag
	lessonID   vo.LessonID
	identifier vo.Identifier
	duration   vo.LessonDuration
	startTime  vo.ScheduleLessonTime
	endTime    vo.ScheduleLessonTime
	roomIndex  vo.RoomIndex
}

func NewScheduleRoomItemModel(
	itemTag vo.RoomItemTag,
	lessonID vo.LessonID,
	identifier vo.Identifier,
	duration vo.LessonDuration,
	startTime vo.ScheduleLessonTime,
	endTime vo.ScheduleLessonTime,
	roomIndex vo.RoomIndex,
) *ScheduleRoomItemModel {

	return &ScheduleRoomItemModel{
		itemTag:    itemTag,
		lessonID:   lessonID,
		identifier: identifier,
		duration:   duration,
		startTime:  startTime,
		endTime:    endTime,
		roomIndex:  roomIndex,
	}
}

func (r ScheduleRoomItemModel) ItemTag() vo.RoomItemTag {
	return r.itemTag
}

func (r ScheduleRoomItemModel) LessonID() vo.LessonID {
	return r.lessonID
}

func (r ScheduleRoomItemModel) Identifier() vo.Identifier {
	return r.identifier
}

func (r ScheduleRoomItemModel) Duration() vo.LessonDuration {
	return r.duration
}

func (r ScheduleRoomItemModel) StartTime() vo.ScheduleLessonTime {
	return r.startTime
}

func (r ScheduleRoomItemModel) EndTime() vo.ScheduleLessonTime {
	return r.endTime
}

func (r ScheduleRoomItemModel) RoomIndex() vo.RoomIndex {
	return r.roomIndex
}

func (r ScheduleRoomItemModel) duplicate() *ScheduleRoomItemModel {
	return &ScheduleRoomItemModel{
		itemTag:    r.itemTag,
		lessonID:   r.lessonID,
		identifier: r.identifier,
		duration:   r.duration,
		startTime:  r.startTime,
		endTime:    r.endTime,
		roomIndex:  r.roomIndex,
	}
}
