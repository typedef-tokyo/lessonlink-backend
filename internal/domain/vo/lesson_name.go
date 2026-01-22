package vo

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrLessonNameEmpty = errors.New("講座名が未設定です")
var ErrLessonNameLengthOver = errors.New("講座名の最大設定位数を超えています")

type LessonName string

const (
	LESSON_NAME_INVALID = LessonName("invalid")
)

func NewLessonName(name string) (LessonName, error) {

	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return LESSON_NAME_INVALID, log.WrapErrorWithStackTrace(ErrLessonNameEmpty)
	}

	const NAME_MAX_LENGTH = 32
	if utf8.RuneCountInString(name) > NAME_MAX_LENGTH {
		return LESSON_NAME_INVALID, log.WrapErrorWithStackTrace(fmt.Errorf("%w 最大:%d文字", ErrLessonNameLengthOver, NAME_MAX_LENGTH))
	}

	return LessonName(name), nil
}

func (r LessonName) Value() string {
	return string(r)
}
