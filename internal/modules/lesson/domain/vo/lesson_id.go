package vo

import (
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrLessonIDUnderMin = errors.New("講座IDは0以上を設定する必要があります。")

type LessonID int

const (
	LESSON_ID_INVALID = LessonID(-1)
	LESSON_ID_INITIAL = LessonID(0)
)

func NewLessonID(id int) (LessonID, error) {

	if id < 0 {
		return LESSON_ID_INVALID, log.WrapErrorWithStackTrace(ErrLessonIDUnderMin)
	}

	return LessonID(id), nil
}

func (r LessonID) Value() int {
	return int(r)
}
