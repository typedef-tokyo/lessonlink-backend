package schedule

import (
	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type ScheduleItemModelSlice []*ScheduleItemModel

func (r ScheduleItemModelSlice) LessonIDs() []vo.LessonID {

	return lo.Map(r, func(item *ScheduleItemModel, _ int) vo.LessonID {
		return item.lessonID
	})
}

func (r ScheduleItemModelSlice) findByIdentifier(identifier vo.Identifier) (*ScheduleItemModel, bool) {

	return lo.Find(r, func(item *ScheduleItemModel) bool {
		return item.identifier == identifier
	})
}

func (r ScheduleItemModelSlice) filterByLessonID(lessonID vo.LessonID) ScheduleItemModelSlice {

	return lo.Filter(r, func(item *ScheduleItemModel, _ int) bool {
		return item.lessonID == lessonID
	})
}

func (r ScheduleItemModelSlice) removeByIdentifier(removeID vo.Identifier) ScheduleItemModelSlice {

	return lo.Filter(r, func(item *ScheduleItemModel, _ int) bool {
		return item.identifier != removeID
	})
}

func (r ScheduleItemModelSlice) addItem(addItem *ScheduleItemModel) ScheduleItemModelSlice {

	return append(r, addItem)
}

type ScheduleItemModel struct {
	lessonID   vo.LessonID
	identifier vo.Identifier
	duration   vo.LessonDuration
}

func NewScheduleItemModel(
	lessonID vo.LessonID,
	identifier vo.Identifier,
	duration vo.LessonDuration,
) *ScheduleItemModel {

	return &ScheduleItemModel{
		lessonID:   lessonID,
		identifier: identifier,
		duration:   duration,
	}
}

func (r *ScheduleItemModel) LessonID() vo.LessonID {
	return r.lessonID
}

func (r *ScheduleItemModel) Identifier() vo.Identifier {
	return r.identifier
}

func (r *ScheduleItemModel) Duration() vo.LessonDuration {
	return r.duration
}

func (r ScheduleItemModel) duplicate() *ScheduleItemModel {

	return &ScheduleItemModel{
		lessonID:   r.lessonID,
		identifier: r.identifier,
		duration:   r.duration,
	}
}
