package presenter

import (
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/lessonlist"
)

type (
	ILessonListPresenter interface {
		Present(result *lessonlist.LessonListQueryOutput) *LessonListResponse
	}

	LessonListPresenter struct {
	}
)

func NewLessonListPresenter() ILessonListPresenter {
	return &LessonListPresenter{}
}

type (
	LessonListResponse struct {
		Lessons []*LessonListDTO `json:"lessons"`
	}

	LessonListDTO struct {
		ID             int    `json:"id"`
		LessonName     string `json:"lesson_name"`
		LessonDuration int    `json:"lesson_duration"`
	}
)

func (h *LessonListPresenter) Present(result *lessonlist.LessonListQueryOutput) *LessonListResponse {

	lessons := lo.Map(result.LessonList, func(item *lessonlist.QueryLessonDTO, _ int) *LessonListDTO {
		return &LessonListDTO{
			ID:             item.ID,
			LessonName:     item.LessonName,
			LessonDuration: item.LessonDuration,
		}
	})

	slices.SortFunc(lessons, func(a, b *LessonListDTO) int {
		return strings.Compare(a.LessonName, b.LessonName)
	})

	return &LessonListResponse{
		Lessons: lessons,
	}
}
