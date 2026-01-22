package lesson

import (
	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type RootLessonModelSlice []*RootLessonModel

func (r RootLessonModelSlice) FindByID(id vo.LessonID) *RootLessonModel {

	model, _ := lo.Find(r, func(item *RootLessonModel) bool {
		return item.id == id
	})

	return model
}

func (r RootLessonModelSlice) CheckDuplicateEntry(model *RootLessonModel) bool {

	_, found := lo.Find(r, func(item *RootLessonModel) bool {
		return item.id != model.id && item.campus == model.campus && item.name == model.name
	})

	return found
}

type RootLessonModel struct {
	id       vo.LessonID
	campus   vo.Campus
	name     vo.LessonName
	duration vo.LessonDuration
}

func NewRootLessonModel(
	id vo.LessonID,
	campus vo.Campus,
	name vo.LessonName,
	duration vo.LessonDuration,
) *RootLessonModel {

	return &RootLessonModel{
		id:       id,
		campus:   campus,
		name:     name,
		duration: duration,
	}
}

func (r RootLessonModel) ID() vo.LessonID {
	return r.id
}

func (r RootLessonModel) Campus() vo.Campus {
	return r.campus
}

func (r RootLessonModel) Name() vo.LessonName {
	return r.name
}

func (r RootLessonModel) Duration() vo.LessonDuration {
	return r.duration
}

func (r *RootLessonModel) Revise(newLessonName vo.LessonName, newDuration vo.LessonDuration) {
	r.name = newLessonName
	r.duration = newDuration
}
