package vo

import (
	"errors"
	"fmt"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrItemLessonDurationUnderMin = errors.New("講座時間は0分以上を設定する必要があります。")
var ErrItemLessonDurationOverMax = errors.New("講座時間に設定できる時間を超えています")
var ErrItemLessonDurationRangeOver = errors.New("分割時間が講座時間の範囲外です")

type LessonDuration int

const (
	LESSON_DURATION_INVALID = LessonDuration(-1)
)

func NewLessonDuration(duration int) (LessonDuration, error) {

	if duration < 0 {
		return LESSON_DURATION_INVALID, log.WrapErrorWithStackTrace(ErrItemLessonDurationUnderMin)
	}

	const max_lesson_time_minutes = 60 * 5
	if duration > max_lesson_time_minutes {
		return LESSON_DURATION_INVALID, log.WrapErrorWithStackTrace(fmt.Errorf("%w 最大:%d分", ErrItemLessonDurationOverMax, max_lesson_time_minutes))
	}

	return LessonDuration(duration), nil
}

func (r LessonDuration) Value() int {
	return int(r)
}

func (r LessonDuration) Add(input LessonDuration) LessonDuration {
	return LessonDuration(r.Value() + input.Value())
}

func (r LessonDuration) StartOffsetMinutes() int {
	return int(r) - 1
}

func (r LessonDuration) Divide(inputDivideMinutes ItemDivideMinutes) (LessonDuration, LessonDuration, error) {

	duration := r.Value()
	divideMinutes := inputDivideMinutes.Value()

	if divideMinutes > (duration - 1) {

		return LESSON_DURATION_INVALID, LESSON_DURATION_INVALID, log.WrapErrorWithStackTrace(ErrItemLessonDurationRangeOver)
	}

	durationDivideFrom := divideMinutes
	durationDivideTo := duration - divideMinutes

	return LessonDuration(durationDivideFrom), LessonDuration(durationDivideTo), nil
}
